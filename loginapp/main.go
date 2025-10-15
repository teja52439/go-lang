package main

import (
    "net/http"
    "strconv"
    "strings"
    "github.com/gin-gonic/gin"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("template/*")

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
    })

    r.POST("/login", func(c *gin.Context) {
        var login LoginRequest

        // Parse JSON body into struct
        if err := c.ShouldBindJSON(&login); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON input"})
            return
        }

        // Check 1: Username must be all uppercase
        if login.Username != strings.ToUpper(login.Username) {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Username must be in CAPITAL letters"})
            return
        }

        // Check 2: Password must be an integer (numeric only)
        if _, err := strconv.Atoi(login.Password); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Password must be an integer"})
            return
        }

        // Login validation (example)
        if login.Username == "ADMIN" && login.Password == "1234" {
            c.JSON(http.StatusOK, gin.H{"message": " Login successful!"})
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"message": " Invalid credentials"})
        }
    })

    r.Run(":8080")}