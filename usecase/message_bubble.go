package usecase

import (
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"context"
)

type MessageBubbleUsecase interface {
	AddMessageBubble(ctx context.Context, body entity.MessageBubble) (*entity.MessageBubble, error)
}

type messageBubbleUsecaseImpl struct {
	messageBubbleRepository repository.MessageBubbleRepository
}

func NewMessageBubbleUsecase(
	messageBubbleRepository repository.MessageBubbleRepository,
) *messageBubbleUsecaseImpl {
	return &messageBubbleUsecaseImpl{
		messageBubbleRepository: messageBubbleRepository,
	}
}

func (u *messageBubbleUsecaseImpl) AddMessageBubble(ctx context.Context, body entity.MessageBubble) (*entity.MessageBubble, error) {
	messageBubble, err := u.messageBubbleRepository.InsertOne(ctx, body)
	if err != nil {
		return nil, err
	}

	return messageBubble, nil
}
