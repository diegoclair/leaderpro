package viewmodel

import (
	"time"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type PersonRequest struct {
	Name       string     `json:"name" validate:"required,min=2,max=450"`
	Email      string     `json:"email,omitempty" validate:"omitempty,email"`
	Position   string     `json:"position,omitempty"`
	Department string     `json:"department,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	Gender     *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
}

func (p PersonRequest) ToEntity() entity.Person {
	return entity.Person{
		Name:       p.Name,
		Email:      p.Email,
		Position:   p.Position,
		Department: p.Department,
		Phone:      p.Phone,
		StartDate:  p.StartDate,
		Notes:      p.Notes,
		Gender:     p.Gender,
		// Set defaults for fields not in the simplified form
		IsManager:   false,
		HasKids:     false,
		Interests:   "",
		Personality: "",
	}
}

type PersonResponse struct {
	UUID               string     `json:"uuid"`
	Name               string     `json:"name"`
	Email              string     `json:"email,omitempty"`
	Position           string     `json:"position,omitempty"`
	Department         string     `json:"department,omitempty"`
	Phone              string     `json:"phone,omitempty"`
	Birthday           *time.Time `json:"birthday,omitempty"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	IsManager          bool       `json:"is_manager"`
	ManagerUUID        string     `json:"manager_uuid,omitempty"`
	Notes              string     `json:"notes,omitempty"`
	HasKids            bool       `json:"has_kids"`
	Gender             *string    `json:"gender,omitempty"`
	Interests          string     `json:"interests,omitempty"`
	Personality        string     `json:"personality,omitempty"`
	LastOneOnOneDate   *time.Time `json:"last_one_on_one_date,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	Age                *int       `json:"age,omitempty"`
	Tenure             *int       `json:"tenure,omitempty"`
}

func (p *PersonResponse) FillFromEntity(person entity.Person) {
	p.UUID = person.UUID
	p.Name = person.Name
	p.Email = person.Email
	p.Position = person.Position
	p.Department = person.Department
	p.Phone = person.Phone
	p.Birthday = person.Birthday
	p.StartDate = person.StartDate
	p.IsManager = person.IsManager
	p.Notes = person.Notes
	p.HasKids = person.HasKids
	p.Gender = person.Gender
	p.Interests = person.Interests
	p.Personality = person.Personality
	p.LastOneOnOneDate = person.LastOneOnOneDate
	p.CreatedAt = person.CreatedAt
	p.Age = person.GetAge()
	p.Tenure = person.GetTenure()
	
	// Manager UUID will be resolved in a future enhancement if needed
	// For now, we keep it empty since we don't have a direct relationship
	p.ManagerUUID = ""
}