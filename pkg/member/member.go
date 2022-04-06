package member

import "time"

type Member struct {
	Id           int64     `json:"id"`
	CreationTime time.Time `json:"creation_time"`
	Name         string    `json:"name" validate:"required"`
}

type MemberRepository interface {
	GetAll(limit, offset, name string) ([]Member, error)
	GetByID(id string) (Member, error)
	GetTotalCount() (int64, error)
	Save(class Member) error
	Update(id string, class Member) error
	Delete(id string) error
}
