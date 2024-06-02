package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type MessageBubbleDTO struct {
	ID             uint      `json:"id"`
	TelemedicineID uint      `json:"telemedicine_id"`
	SenderType     string    `json:"sender_type"`
	Payload        string    `json:"payload"`
	PayloadType    string    `json:"payload_type"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewMessageBubbleDTO(p entity.MessageBubble) MessageBubbleDTO {
	return MessageBubbleDTO{
		ID:             p.ID,
		TelemedicineID: p.TelemedicineID,
		SenderType:     p.SenderType,
		Payload:        p.Payload,
		PayloadType:    p.PayloadType,
		CreatedAt:      p.CreatedAt,
	}
}
