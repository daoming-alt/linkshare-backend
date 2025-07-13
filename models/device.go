package models

type Device struct {
    ID       int    `json:"id"`
    UserID   int    `json:"user_id"`
    Name     string `json:"name"`
    LastSeen string `json:"last_seen"`
}