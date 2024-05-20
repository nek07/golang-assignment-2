package models

import "time"

type ChatMessage struct {
	ChatID    string    `bson:"chat_id"`
	Sender    string    `bson:"sender"`
	Message   string    `bson:"message"`
	Timestamp time.Time `bson:"timestamp"`
}
