package class

import "time"

type Class struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"name" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
	Capacity  int16     `json:"capacity" validate:"required"`
}

type ClassRepository interface {
	GetAll(string, string) ([]Class, error)
	GetByID(id string) (Class, error)
	GetByDateRange(startDate, endDate string) ([]Class, error)
	GetByName(name string) ([]Class, error)
	GetTotalCount() (int64, error)
	Save(class Class) error
	Update(id string, class Class) error
	Delete(id string) error
}
