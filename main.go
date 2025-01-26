package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// Models

type Patient struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	History  string `json:"history"` // Placeholder for medical history
}

type Appointment struct {
	ID         uint   `gorm:"primaryKey"`
	PatientID  uint   `json:"patient_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Description string `json:"description"`
}

type Inventory struct {
	ID       uint   `gorm:"primaryKey"`
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
}

// Initialize database
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("healthcare.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Patient{}, &Appointment{}, &Inventory{})
	return db, nil
}

func main() {
	db, err := initDB()
	if err != nil {
		panic("Failed to connect to database")
	}

	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Routes

	// Patient Routes
	r.GET("/patients", func(c *gin.Context) {
		var patients []Patient
		db.Find(&patients)
		c.HTML(http.StatusOK, "patients.html", gin.H{
			"patients": patients,
		})
	})

	r.POST("/patients", func(c *gin.Context) {
		var patient Patient
		if err := c.ShouldBindJSON(&patient); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&patient)
		c.HTML(http.StatusCreated, "patient_created.html", gin.H{
			"patient": patient,
		})
	})

	// Appointment Routes
	r.GET("/appointments", func(c *gin.Context) {
		var appointments []Appointment
		db.Find(&appointments)
		c.HTML(http.StatusOK, "appointments.html", gin.H{
			"appointments": appointments,
		})
	})

	r.POST("/appointments", func(c *gin.Context) {
		var appointment Appointment
		if err := c.ShouldBindJSON(&appointment); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&appointment)
		c.HTML(http.StatusCreated, "appointment_created.html", gin.H{
			"appointment": appointment,
		})
	})

	// Inventory Routes
	r.GET("/inventory", func(c *gin.Context) {
		var inventory []Inventory
		db.Find(&inventory)
		c.HTML(http.StatusOK, "inventory.html", gin.H{
			"inventory": inventory,
		})
	})

	r.POST("/inventory", func(c *gin.Context) {
		var inventory Inventory
		if err := c.ShouldBindJSON(&inventory); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&inventory)
		c.HTML(http.StatusCreated, "inventory_created.html", gin.H{
			"inventory": inventory,
		})
	})

	r.Run() // Run the server
}
