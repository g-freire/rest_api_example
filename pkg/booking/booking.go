package booking

import (
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
	GetAll(limit, offset string) ([]Booking, error)
	GetByID(id string) (Booking, error)
	GetByDateRange(startDate, endDate string) ([]Booking, error)
	GetTotalCount() (int64, error)
	GetAllClassesByMemberId(memberId string) ([]class.Class, error)
	GetAllMembersByClassId(classId string) ([]member.Member, error)
	Save(class Booking) error
	Update(id string, class Booking) error
	Delete(id string) error
}

type BookingService interface {
	GetByDateRange(startDate, endDate string) ([]Booking, error)
	Save(class Booking) error
	Update(id string, class Booking) error
}
