package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type MessageBubble struct {
	TelemedicineID uint   `json:"telemedicine_id" binding:"required"`
	SenderType     string `json:"sender_type" binding:"required,oneof=doctor user"`
	Payload        string `json:"payload" binding:"required"`
	PayloadType    string `json:"payload_type" binding:"required,oneof=text image document"`
}

func (req *MessageBubble) MessageBubble() entity.MessageBubble {
	return entity.MessageBubble{
		TelemedicineID: req.TelemedicineID,
		SenderType:     req.SenderType,
		PayloadType:    req.PayloadType,
		Payload:        req.Payload,
	}
}
