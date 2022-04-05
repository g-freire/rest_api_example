package class

import (
	"gym/internal/errors"
	"time"
)

type service struct {
	ClassRepository ClassRepository
}

func NewService(ClassRepo ClassRepository) ClassService {
	return &service{
		ClassRepository: ClassRepo,
	}
}

// check if start date is before end date
func checkTimestampIsValid(startTime, endTime time.Time) bool {
	if startTime.Before(endTime) {
		return true
	} else {
		return false
	}
}

func (s service) GetByDateRange(startDate, endDate string) ([]Class, error) {
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
	if checkTimestampIsValid(startTime, endTime) == false {
		return nil, errors.ErrInvalidTimestamp
	}
	return  s.ClassRepository.GetByDateRange(startDate, endDate)
}

func (s service) Save(class Class) error {
	// validates time chronology
	if checkTimestampIsValid(class.StartDate, class.EndDate) == false {
		return errors.ErrInvalidTimestamp
	}
	return s.ClassRepository.Save(class)
}

func (s service) Update(id string, class Class) error {
	// validates time chronology
	if checkTimestampIsValid(class.StartDate, class.EndDate) == false {
		return errors.ErrInvalidTimestamp
	}
	return s.ClassRepository.Update(id, class)
}
