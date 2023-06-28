package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetUniversities(c *gin.Context) {

	InitLogger()
	var universities []University
	DB.Find(&universities)
	c.JSON(http.StatusOK, gin.H{"data": universities})

}

func GetUnivByRank(c *gin.Context) {

	InitLogger()
	var university University
	var source string
	data, err, source := GetCache(c.Param("rank"))
	if !err {
		fmt.Println("Cache miss")
		rankint, _ := strconv.Atoi(c.Param("rank"))
		// Queries the database for the university record with the given rank
		dberr := DB.Where(&University{Rank: uint16(rankint)}).First(&university).Error
		if errors.Is(dberr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"data": "Record not found"})
			if file != nil {
				file.WriteString(fmt.Sprintf("[%s] given rank university record not found\n", time.Now().Format("2006/01/02 - 15:04:05")))
			}

		} else {

			// Sets the retrieved university record to the cache
			SetCache(strconv.Itoa(int(university.Rank)), university)
			source = "database"
			c.JSON(http.StatusOK, gin.H{"data": university, "source": source})
			if file != nil {
				file.WriteString(fmt.Sprintf("[%s] Data retrieved from %s\n", time.Now().Format("2006/01/02 - 15:04:05"), source))
			}
		}

	} else {

		// Retrieves the university record from the cache
		university = data.(University)
		c.JSON(http.StatusOK, gin.H{"data": university, "source": source})
		if file != nil {
			file.WriteString(fmt.Sprintf("[%s] Data retrieved from Cache\n", time.Now().Format("2006/01/02 - 15:04:05")))
		}
	}
}

// Delete University By Rank

func deleteUniversity(c *gin.Context) {

	var university University
	rankUpd, err := strconv.Atoi(c.Param("rank"))
	if err != nil {
		panic(err)
	}
	if err := DB.First(&university, rankUpd).Error; err != nil {
		logger.Printf("University Not Found for rank %d", rankUpd)
		c.JSON(http.StatusNotFound, gin.H{"Error": "University Not Found!!"})
		return
	}

	if err := DB.Delete(&University{}, c.Param("rank")).Error; err != nil {
		logger.Printf("Failed to delete university with rank %d", rankUpd)
		c.JSON(500, gin.H{"Error": "Failed to delete university"})
		return
	}

	logger.Printf("Successfully deleted university with rank %d", rankUpd)

	c.JSON(http.StatusOK, gin.H{"Status": "Deleted University succcessfully!!"})

}

// Update University data by Rank

func updateUniversity(c *gin.Context) {

	var institution University
	var input University
	rankUpd, _ := strconv.Atoi(c.Param("rank"))
	if err := DB.Where(University{Rank: uint16(rankUpd)}).First(&institution).Error; err != nil {
		logger.Printf("University Not Found for rank %d", rankUpd)
		c.JSON(http.StatusNotFound, gin.H{"Error": "University Not Found!!!"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Printf("Unable to update university for rank %d", rankUpd)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"Error": "Not able to Update data"})
		return
	}

	log.Printf("Successfully updated university for rank %d", rankUpd)
	DB.Model(&institution).Updates(&input)
	c.JSON(http.StatusOK, gin.H{"data": institution})

}
