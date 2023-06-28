package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var file *os.File
var logger *log.Logger

func InitLogger() {
	var err error
	file, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logger = log.New(io.MultiWriter(file, os.Stdout), "", log.LstdFlags)
}

func LogError(err error) {
	if logger != nil {
		logger.Println("ERROR:", err)
	}
}

// LoggerMiddleware is a middleware that logs the incoming requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		logTime := end.Format("2006/01/02 - 15:04:05")
		logMessage := fmt.Sprintf("[%s] %s %s %s %v", logTime, c.ClientIP(), c.Request.Method, c.Request.URL.Path, latency)
		if logger != nil {
			logger.Println(logMessage + "\n")
		}
	}
}
