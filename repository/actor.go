package repository

import "context"

type ActorRepository interface {
	UpdatePassword(ctx context.Context, id uint, password string) error
}
