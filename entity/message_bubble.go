package entity

import "time"

type MessageBubble struct {
	ID             uint
	TelemedicineID uint
	SenderType     string
	Payload        string
	PayloadType    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
