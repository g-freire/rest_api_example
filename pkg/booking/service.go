package booking

import (
	"context"
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
func (s service) GetByDateRange(ctx context.Context,startDate, endDate string) ([]Booking, error) {
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
	return s.BookingRepository.GetByDateRange(ctx, startDate, endDate)
}

func (s service) Save(ctx context.Context,booking Booking)  (int64 ,error) {
	// validates time chronology
	if common.CheckTimestampIsUpToDate(booking.Date) == false {
		return 0, errors.ErrOldTimestamp
	}
	return s.BookingRepository.Save(ctx, booking)
}

func (s service) Update(ctx context.Context,id string, booking Booking) error {
	// validates time chronology
	if common.CheckTimestampIsUpToDate(booking.Date) == false {
		return errors.ErrOldTimestamp
	}
	return s.BookingRepository.Update(ctx, id, booking)
}
