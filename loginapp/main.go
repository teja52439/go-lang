package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Struct for JSON input
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Struct for DB model
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"unique"`
    Password string
}

var db *gorm.DB

func main() {
    var err error
    db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
    fmt.Print(err)
    if err != nil {
        panic("failed to connect database")
    }

    // Auto-create table if not exists
    db.AutoMigrate(&User{})

    r := gin.Default()
    r.LoadHTMLGlob("template/*")

    // Show login page
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
    })

    // Register user (you can create users here)
    r.POST("/register", func(c *gin.Context) {
        var req LoginRequest

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
            return
        }

        // Validation
        if req.Username != strings.ToUpper(req.Username) {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Username must be in CAPITAL letters"})
            return
        }
        if _, err := strconv.Atoi(req.Password); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Password must be integer"})
            return
        }

        user := User{Username: req.Username, Password: req.Password}
        if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Username already exists"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User registered successfully!"})
    })

    // Login route
    r.POST("/login", func(c *gin.Context) {
        var login LoginRequest
        if err := c.ShouldBindJSON(&login); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON input"})
            return
        }

        // Validation checks
        if login.Username != strings.ToUpper(login.Username) {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Username must be in CAPITAL letters"})
            return
        }
        if _, err := strconv.Atoi(login.Password); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Password must be an integer"})
            return
        }

        // Database lookup
        var user User
        if err := db.Where("username = ? AND password = ?", login.Username, login.Password).First(&user).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
    })

    r.Run(":8080")
}
