package booking

import (
	"context"
	"gym/pkg/class"
	"gym/pkg/member"
	"time"
)

type Booking struct {
	Id           int64     `json:"id"`
	ClassId      int64     `json:"class_id"`
	MemberId     int64     `json:"member_id"`
	Date         time.Time `json:"date" validate:"required"`
	CreationTime time.Time `json:"creation_time"`
}

type BookingRepository interface {
	GetAll(ctx context.Context, limit, offset string) ([]Booking, error)
	GetByID(ctx context.Context, id string) (Booking, error)
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]Booking, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetAllClassesByMemberId(ctx context.Context, memberId string) ([]class.Class, error)
	GetAllMembersByClassId(ctx context.Context, classId string) ([]member.Member, error)
	Save(ctx context.Context, class Booking) (int64, error)
	Update(ctx context.Context, id string, class Booking) error
	Delete(ctx context.Context, id string) error
}

type BookingService interface {
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]Booking, error)
	Save(ctx context.Context, class Booking) (int64, error)
	Update(ctx context.Context, id string, class Booking) error
}
