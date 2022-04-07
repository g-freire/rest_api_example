package class

import (
	"context"
	"gym/internal/common"
	"gym/internal/errors"
	"time"
)

type service struct {
	ClassRepository ClassRepository
}

func NewService(ClassRepo ClassRepository) *service {
	return &service{
		ClassRepository: ClassRepo,
	}
}

// this function could be adjusted to this entity to handle hours/minutes depending on the business requirements
func (s service) GetByDateRange(ctx context.Context, startDate, endDate string) ([]Class, error) {
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
	return  s.ClassRepository.GetByDateRange(ctx, startDate, endDate)
}

func (s service) Save(ctx context.Context, class Class) (int,error) {
	// validates time chronology
	if common.CheckTimestampIsUpToDate(class.StartDate) == false {
		return 0, errors.ErrOldTimestamp
	}
	if common.CheckTimestampIsUpToDate(class.EndDate) == false {
		return 0,errors.ErrOldTimestamp
	}
	if common.CheckTimestampIsValid(class.StartDate, class.EndDate) == false {
		return 0,errors.ErrInvalidTimestamp
	}
	return s.ClassRepository.Save(context.TODO(), class)
}

func (s service) Update(ctx context.Context, id string, class Class) error {
	// validates time chronology
	if common.CheckTimestampIsUpToDate(class.StartDate) == false {
		return errors.ErrOldTimestamp
	}
	if common.CheckTimestampIsUpToDate(class.EndDate) == false {
		return errors.ErrOldTimestamp
	}
	if common.CheckTimestampIsValid(class.StartDate, class.EndDate) == false {
		return errors.ErrInvalidTimestamp
	}
	return s.ClassRepository.Update(ctx, id, class)
}