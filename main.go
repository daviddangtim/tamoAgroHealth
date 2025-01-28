package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	// Serve static files (frontend)
	r.Static("/static", "./static") // Serve static files like CSS, JS
	r.LoadHTMLGlob("templates/*")   // Load HTML templates

	// CORS configuration (if still needed for other external APIs)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://tamo-front.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "hx-current-url"},
		AllowCredentials: true,
	}))

	// Routes

	// Serve the main frontend pages
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/patients", func(c *gin.Context) {
		c.HTML(http.StatusOK, "patients_page.html", nil)
	})

	r.GET("/appointments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "appointments_page.html", nil)
	})

	r.GET("/inventory", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inventory_page.html", nil)
	})

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
		// Read the raw request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to read request body"})
			return
		}
		fmt.Printf("Raw request body: %s\n", string(body))
	
		// Bind JSON to struct
		var inventory Inventory
		if err := c.ShouldBindJSON(&inventory); err != nil {
			fmt.Printf("Binding error: %v\n", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	
		// Save to database
		db.Create(&inventory)
	
		// Return HTML with the created inventory data
		html := fmt.Sprintf("<div>%s: %d</div>", inventory.ItemName, inventory.Quantity)
		c.Data(http.StatusCreated, "text/html", []byte(html))
	})

	r.Run() // Run the server
}