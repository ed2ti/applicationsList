// main.go
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

// Configuration represents system settings
type Configuration struct {
	ID          uint `gorm:"primaryKey"`
	CardsPerRow int  `gorm:"not null;default:3"`
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

	// Auto migrate the schemas
	db.AutoMigrate(&Application{}, &Configuration{})

	// Initialize default configuration if not exists
	var config Configuration
	if db.First(&config).Error != nil {
		db.Create(&Configuration{CardsPerRow: 3})
	}
}

func main() {
	initDatabase()

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Add template function for division
	r.SetFuncMap(template.FuncMap{
		"divide": func(a, b int) int {
			return a / b
		},
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Static("/images", "./images")

	r.GET("/", func(c *gin.Context) {
		var apps []Application
		var config Configuration
		db.Find(&apps)
		db.First(&config)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"apps":   apps,
			"config": config,
		})
	})

	r.GET("/newApp", func(c *gin.Context) {
		c.HTML(http.StatusOK, "newApp.html", nil)
	})

	r.GET("/config", func(c *gin.Context) {
		var config Configuration
		db.First(&config)
		c.HTML(http.StatusOK, "config.html", gin.H{
			"Config": config,
		})
	})

	r.POST("/saveConfig", func(c *gin.Context) {
		cardsPerRow := c.PostForm("cardsPerRow")
		cards, err := strconv.Atoi(cardsPerRow)
		if err != nil || cards < 2 || cards > 5 {
			c.String(http.StatusBadRequest, "Invalid cards per row value")
			return
		}

		var config Configuration
		db.First(&config)
		config.CardsPerRow = cards
		db.Save(&config)

		c.Redirect(http.StatusSeeOther, "/")
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
