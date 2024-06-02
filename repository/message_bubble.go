package repository

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

type MessageBubbleRepository interface {
	InsertOne(ctx context.Context, newMessageBubble entity.MessageBubble) (*entity.MessageBubble, error)
}

type messageBubbleRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewMessageBubblesRepository(db transaction.DBTransaction) *messageBubbleRepositoryImpl {
	return &messageBubbleRepositoryImpl{
		db: db,
	}
}

func (r *messageBubbleRepositoryImpl) InsertOne(ctx context.Context, newMessageBubble entity.MessageBubble) (*entity.MessageBubble, error) {
	q := `
		INSERT INTO message_bubbles (
			telemedicine_id,
			sender_type,
			payload, 
			payload_type
		) VALUES
			($1, $2, $3, $4)
		RETURNING
			message_bubble_id,
			telemedicine_id,
			sender_type,
			payload, 
			payload_type,
			created_at
	`

	var scan entity.MessageBubble
	err := r.db.QueryRowContext(ctx, q,
		newMessageBubble.TelemedicineID,
		newMessageBubble.SenderType,
		newMessageBubble.Payload,
		newMessageBubble.PayloadType,
	).Scan(
		&scan.ID,
		&scan.TelemedicineID,
		&scan.SenderType,
		&scan.Payload,
		&scan.PayloadType,
		&scan.CreatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}
