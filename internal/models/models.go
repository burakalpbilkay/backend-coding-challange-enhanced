package models

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type Action struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	UserID     int    `json:"userId"`
	TargetUser int    `json:"targetUser"`
	CreatedAt  string `json:"createdAt"`
}
