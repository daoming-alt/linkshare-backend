// handlers/link.go
package handlers

import (
	"fmt"
	"linkshare-backend/db"
	"linkshare-backend/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins for simplicity
}

// WebSocket connections map (user_id -> device_id -> connection)
var (
    connections = make(map[int]map[int]*websocket.Conn)
    connMutex   sync.Mutex
)

// SendLink sends a link to a target device
// @Summary Send a link
// @Description Sends a link from one device to another
// @Tags links
// @Accept json
// @Produce json
// @Param link body models.Link true "Link data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /links [post]
func SendLink(c *gin.Context) {
    var link models.Link
    if err := c.ShouldBindJSON(&link); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    userID, _ := c.Get("user_id")
    link.UserID = userID.(int)

    // Store link in database
    _, err := db.DB.Exec("INSERT INTO links (user_id, from_device_id, to_device_id, url) VALUES (?, ?, ?, ?)",
        link.UserID, link.FromDeviceID, link.ToDeviceID, link.URL)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to store link"})
        return
    }

    // Send link to target device via WebSocket
    connMutex.Lock()
    if deviceConns, exists := connections[link.UserID]; exists {
        if conn, exists := deviceConns[link.ToDeviceID]; exists {
            err := conn.WriteJSON(link)
            if err != nil {
                log.Println("WebSocket write error:", err)
            }
        }
    }
    connMutex.Unlock()

    c.JSON(http.StatusOK, gin.H{"message": "Link sent"})
}

// WebSocketHandler handles WebSocket connections for real-time link delivery
// @Summary WebSocket for link delivery
// @Description Establishes a WebSocket connection for a device
// @Tags websocket
// @Produce json
// @Router /ws [get]
func WebSocketHandler(c *gin.Context) {
    userID, _ := c.Get("user_id")
    deviceID := c.Query("device_id") // Device ID passed as query param

    // Validate device_id
    parsedDeviceID := parseDeviceID(deviceID)
    if parsedDeviceID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device_id"})
        return
    }

    // Upgrade HTTP to WebSocket
    // Ensure the Authorization header is passed correctly
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }

    // Register connection
    connMutex.Lock()
    if _, exists := connections[userID.(int)]; !exists {
        connections[userID.(int)] = make(map[int]*websocket.Conn)
    }
    connections[userID.(int)][parsedDeviceID] = conn
    connMutex.Unlock()

    // Send connection confirmation message
    err = conn.WriteJSON(map[string]string{"message": "Connected"})
    if err != nil {
        log.Println("WebSocket write error for connection message:", err)
        conn.Close()
        return
    }

    // Update device last seen
    _, err = db.DB.Exec("UPDATE devices SET last_seen = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?", parsedDeviceID, userID)
    if err != nil {
        log.Println("Failed to update last_seen:", err)
    }

    // Keep connection open
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Println("WebSocket read error:", err)
            connMutex.Lock()
            delete(connections[userID.(int)], parsedDeviceID)
            if len(connections[userID.(int)]) == 0 {
                delete(connections, userID.(int))
            }
            connMutex.Unlock()
            conn.Close()
            return
        }
    }
}

func parseDeviceID(deviceID string) int {
    // Convert deviceID string to int with error handling
    var id int
    _, err := fmt.Sscanf(deviceID, "%d", &id)
    if err != nil {
        log.Println("Invalid device_id:", deviceID)
        return 0
    }
    return id
}