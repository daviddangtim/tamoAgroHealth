package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
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

	// CORS configuration (allow requests from w3spaces)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://tamo-front.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "hx-current-url"},
		AllowCredentials: true,
	}))

	// Routes

	// Patient Routes
	r.GET("/patients/all", func(c *gin.Context) {
		var patients []Patient
		db.Find(&patients)

		// Return HTML with data for htmx to parse
		var html string
		for _, patient := range patients {
			html += fmt.Sprintf("<div>%s - %d</div>", patient.Name, patient.Age)
		}
		c.Data(http.StatusOK, "text/html", []byte(html))
	})

	r.POST("/patients", func(c *gin.Context) {
		var patient Patient
		if err := c.ShouldBindJSON(&patient); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&patient)

		// Return HTML with the created patient data
		html := fmt.Sprintf("<div>%s - %d</div>", patient.Name, patient.Age)
		c.Data(http.StatusCreated, "text/html", []byte(html))
	})

	// Appointment Routes
	r.GET("/appointments/all", func(c *gin.Context) {
		var appointments []Appointment
		db.Find(&appointments)

		// Return HTML with appointment data for htmx to parse
		var html string
		for _, appointment := range appointments {
			html += fmt.Sprintf("<div>Appointment: %s at %s</div>", appointment.Date, appointment.Time)
		}
		c.Data(http.StatusOK, "text/html", []byte(html))
	})

	r.POST("/appointments", func(c *gin.Context) {
		var appointment Appointment
		if err := c.ShouldBindJSON(&appointment); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&appointment)

		// Return HTML with the created appointment data
		html := fmt.Sprintf("<div>Appointment: %s at %s</div>", appointment.Date, appointment.Time)
		c.Data(http.StatusCreated, "text/html", []byte(html))
	})

	// Inventory Routes
	r.GET("/inventory/all", func(c *gin.Context) {
		var inventory []Inventory
		db.Find(&inventory)

		// Return HTML with inventory data for htmx to parse
		var html string
		for _, item := range inventory {
			html += fmt.Sprintf("<div>%s: %d</div>", item.ItemName, item.Quantity)
		}
		c.Data(http.StatusOK, "text/html", []byte(html))
	})

	r.POST("/inventory", func(c *gin.Context) {
		var inventory Inventory
		if err := c.ShouldBindJSON(&inventory); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&inventory)

		// Return HTML with the created inventory data
		html := fmt.Sprintf("<div>%s: %d</div>", inventory.ItemName, inventory.Quantity)
		c.Data(http.StatusCreated, "text/html", []byte(html))
	})

	r.Run() // Run the server
}