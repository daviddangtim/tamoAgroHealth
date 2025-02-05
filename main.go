package main

import (
	"fmt"
	"net/http"
	// "github.com/gin-gonic/gin/binding"
	"strings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Models

type Patient struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	History  string `json:"history"` // Placeholder for medical history
}

type Appointment struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey"`
	PatientID  uint   `json:"patient_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Description string `json:"description"`
}

type Inventory struct {
	gorm.Model
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

	r.Static("/static", "./static") 
	r.LoadHTMLGlob("templates/*")   
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://tamo-front.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "hx-current-url"},
		AllowCredentials: true,
	}))


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

	r.GET("/patients/all", func(c *gin.Context) {
		var patients []Patient
		db.Find(&patients)

		var html string
		for _, patient := range patients {
			html += fmt.Sprintf("<div>%s - %d</div>", patient.Name, patient.Age)
		}
		c.Data(http.StatusOK, "text/html", []byte(html))
	})

	r.POST("/patients", func(c *gin.Context) {
		var patient Patient
		if err := c.Bind(&patient); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&patient)

		html := fmt.Sprintf("<div>%s - %d</div>", patient.Name, patient.Age)
		c.Data(http.StatusCreated, "text/html", []byte(html))
	})

	r.GET("/appointments/all", func(c *gin.Context) {
		var appointments []Appointment
		db.Find(&appointments)

		var html string
		for _, appointment := range appointments {
			html += fmt.Sprintf("<div>Appointment: %s at %s</div>", appointment.Date, appointment.Time)
		}
		c.Data(http.StatusOK, "text/html", []byte(html))
	})

	r.POST("/appointments", func(c *gin.Context) {
		var appointment Appointment
		if err := c.Bind(&appointment); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&appointment)

		html := fmt.Sprintf("<div>Appointment: %s at %s</div>", appointment.Date, appointment.Time)
		c.Data(http.StatusCreated, "text/html", []byte(html))
	})

	r.GET("/inventory/all", func(c *gin.Context) {
		var inventory []Inventory
		db.Find(&inventory)

		var html string
		for _, item := range inventory {
			html += fmt.Sprintf("<div>%s: %d</div>", item.ItemName, item.Quantity)
		}
		c.Data(http.StatusOK, "text/html", []byte(html))
	})

	r.GET("/inventory/list", func(c *gin.Context) {
		var inventories []Inventory
		db.Find(&inventories)
		fmt.Printf("Fetched inventories: %+v\n", inventories) // Log the fetched data
	
		var html strings.Builder
		if len(inventories) == 0 {
			html.WriteString(`<p class="text-center text-gray-400">No inventory items available.</p>`)
		} else {
			for _, inventory := range inventories {
				html.WriteString(fmt.Sprintf(`<div class="inventory-item">%s: %d</div>`, inventory.ItemName, inventory.Quantity))
			}
		}
	
		c.Data(http.StatusOK, "text/html", []byte(html.String()))
	})

	r.POST("/inventory", func(c *gin.Context) {
		// Log the raw form data
		fmt.Printf("Raw form data: %v\n", c.Request.PostForm)
		
		var inventory Inventory
		if err := c.Bind(&inventory); err != nil {
			fmt.Printf("Binding error: %v\n", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	
		fmt.Printf("After binding: %+v\n", inventory)
	
		result := db.Create(&inventory)
		if result.Error != nil {
			fmt.Printf("Database error: %v\n", result.Error)
			c.JSON(500, gin.H{"error": "Failed to save to database"})
			return
		}
	
		fmt.Printf("After save: %+v\n", inventory)
	
		html := fmt.Sprintf("<div>%s: %d</div>", inventory.ItemName, inventory.Quantity)
		c.Data(http.StatusCreated, "text/html", []byte(html))
	})

	r.Run() // Run the server
}