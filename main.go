// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Application represents an app in the panel
type Application struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Logo string
	Link string `gorm:"not null"`
}

var db *gorm.DB

func initDatabase() {
	var err error
	// Ensure the directory exists
	os.MkdirAll("db", os.ModePerm)
	db, err = gorm.Open(sqlite.Open("db/applications.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	db.AutoMigrate(&Application{})
}

func main() {
	initDatabase()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Static("/images", "./images")

	r.GET("/", func(c *gin.Context) {
		var apps []Application
		db.Find(&apps)
		c.HTML(http.StatusOK, "index.html", gin.H{"apps": apps})
	})

	r.POST("/add", func(c *gin.Context) {
		name := c.PostForm("name")
		link := c.PostForm("link")

		// Upload the image
		file, err := c.FormFile("logo")
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to upload logo: %s", err.Error())
			return
		}
		imagePath := filepath.Join("images", file.Filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.String(http.StatusInternalServerError, "Failed to save logo: %s", err.Error())
			return
		}

		db.Create(&Application{Name: name, Logo: imagePath, Link: link})
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.POST("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		var app Application
		if err := db.First(&app, id).Error; err != nil {
			c.String(http.StatusNotFound, "Application not found")
			return
		}

		// Delete the image file
		if err := os.Remove(app.Logo); err != nil {
			log.Printf("Failed to delete image file: %s", err.Error())
		}

		db.Delete(&app)
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.Run(":80")
}
