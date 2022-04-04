package booking

import "time"

type Booking struct {
	Id           int64     `json:"id"`
	CreationTime time.Time `json:"creation_time"`
	Name         string    `json:"name" validate:"required"`
	Date         time.Time `json:"date" validate:"required"`
}

type BookingRepository interface {
	GetAll(limit, offset, name string) ([]Booking, error)
	GetByID(id string) (Booking, error)
	GetByDateRange(startDate, endDate string) ([]Booking, error)
	GetTotalCount() (int64, error)
	Save(class Booking) error
	Update(id string, class Booking) error
	Delete(id string) error
}
