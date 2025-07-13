// handlers/device.go
package handlers

import (
	"linkshare-backend/db"
	"linkshare-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterDevice registers a new device for a user
// @Summary Register a device
// @Description Registers a device with a name for the authenticated user
// @Tags devices
// @Accept json
// @Produce json
// @Param device body models.Device true "Device data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /devices [post]
func RegisterDevice(c *gin.Context) {
    var device models.Device
    if err := c.ShouldBindJSON(&device); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    userID, _ := c.Get("user_id")
    device.UserID = userID.(int)

    _, err := db.DB.Exec("INSERT INTO devices (user_id, name) VALUES (?, ?)", device.UserID, device.Name)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to register device"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Device registered"})
}

// ListDevices returns a list of a user’s devices
// @Summary List devices
// @Description Returns all devices registered by the authenticated user
// @Tags devices
// @Produce json
// @Success 200 {array} models.Device
// @Failure 500 {object} map[string]string
// @Router /devices [get]
func ListDevices(c *gin.Context) {
    userID, _ := c.Get("user_id")
    rows, err := db.DB.Query("SELECT id, user_id, name, last_seen FROM devices WHERE user_id = ?", userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
        return
    }
    defer rows.Close()

    var devices []models.Device
    for rows.Next() {
        var device models.Device
        if err := rows.Scan(&device.ID, &device.UserID, &device.Name, &device.LastSeen); err != nil {
            continue
        }
        devices = append(devices, device)
    }

    c.JSON(http.StatusOK, devices)
}