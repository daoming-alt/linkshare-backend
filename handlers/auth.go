// handlers/auth.go
package handlers

import (
	"crypto/sha256"
	"fmt"
	"linkshare-backend/db"
	"linkshare-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // Replace with a secure key

// Register handles user registration
// @Summary Register a new user
// @Description Creates a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Hash password (use bcrypt in production)
    hash := sha256.Sum256([]byte(user.Password))
    user.Password = fmt.Sprintf("%x", hash)

    // Insert user into database
    _, err := db.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

// Login handles user login and returns a JWT
// @Summary User login
// @Description Authenticates a user and returns a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Hash password for comparison
    hash := sha256.Sum256([]byte(user.Password))
    hashedPassword := fmt.Sprintf("%x", hash)

    // Check credentials
    var dbUser models.User
    err := db.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)
    if err != nil || dbUser.Password != hashedPassword {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Generate JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": dbUser.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}