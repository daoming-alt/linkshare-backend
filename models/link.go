package models

type Link struct {
    ID           int    `json:"id"`
    UserID       int    `json:"user_id"`
    FromDeviceID int    `json:"from_device_id"`
    ToDeviceID   int    `json:"to_device_id"`
    URL          string `json:"url"`
    CreatedAt    string `json:"created_at"`
}