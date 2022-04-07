package class

import (
	"context"
	"time"
)

type Class struct {
	Id           int64     `json:"id"`
	CreationTime time.Time `json:"creation_time"`
	Name         string    `json:"name" validate:"required"`
	StartDate    time.Time `json:"start_date" validate:"required"`
	EndDate      time.Time `json:"end_date" validate:"required"`
	Capacity     int16     `json:"capacity" validate:"required"`
}

type ClassRepository interface {
	GetAll(ctx context.Context, limit, offset, name string) ([]Class, error)
	GetByID(ctx context.Context, id string) (Class, error)
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]Class, error)
	GetTotalCount(ctx context.Context) (int64, error)
	Save(ctx context.Context,class Class) (int,error)
	Update(ctx context.Context, id string, class Class) error
	Delete(ctx context.Context, id string) error
}

type ClassService interface {
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]Class, error)
	Save(ctx context.Context, class Class) (int, error)
	Update(ctx context.Context, id string, class Class) error
}