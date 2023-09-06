package models

import "time"

type SessionInfo struct {
	SessionID string    `bson:"session"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}
