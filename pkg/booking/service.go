package booking

import (
	"gym/internal/common"
	"gym/internal/errors"
	"time"
)

type service struct {
	BookingRepository BookingRepository
}

func NewService(ClassRepo BookingRepository) BookingService {
	return &service{
		BookingRepository: ClassRepo,
	}
}

// this function could be adjusted to this entity to handle hours/minutes depending on the business requirements
func (s service) GetByDateRange(startDate, endDate string) ([]Booking, error) {
	// cast string to time
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}
	// validate time chronology
	if common.CheckTimestampIsValid(startTime, endTime) == false {
		return nil, errors.ErrInvalidTimestamp
	}
	return s.BookingRepository.GetByDateRange(startDate, endDate)
}

func (s service) Save(booking Booking) error {
	// validates time chronology
	if common.CheckTimestampIsUpToDate(booking.Date) == false {
		return errors.ErrInvalidTimestamp
	}
	return s.BookingRepository.Save(booking)
}

func (s service) Update(id string, booking Booking) error {
	// validates time chronology
	if common.CheckTimestampIsUpToDate(booking.Date) == false {
		return errors.ErrInvalidTimestamp
	}
	return s.BookingRepository.Update(id, booking)
}
