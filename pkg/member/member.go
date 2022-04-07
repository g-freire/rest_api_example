package member

import (
	"context"
	"time"
)

type Member struct {
	Id           int64     `json:"id"`
	CreationTime time.Time `json:"creation_time"`
	Name         string    `json:"name" validate:"required"`
}

type MemberRepository interface {
	GetAll(ctx context.Context, limit, offset, name string) ([]Member, error)
	GetByID(ctx context.Context, id string) (Member, error)
	GetTotalCount(ctx context.Context) (int64, error)
	Save(ctx context.Context, class Member) (int, error)
	Update(ctx context.Context, id string, class Member) error
	Delete(ctx context.Context, id string) error
}
