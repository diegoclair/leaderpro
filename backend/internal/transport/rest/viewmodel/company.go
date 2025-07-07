package viewmodel

import (
	"time"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type CompanyRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description,omitempty"`
	Industry    string `json:"industry,omitempty"`
	Size        string `json:"size,omitempty" validate:"omitempty,oneof=small medium large enterprise"`
}

func (c CompanyRequest) ToEntity() entity.Company {
	return entity.Company{
		Name:        c.Name,
		Description: c.Description,
		Industry:    c.Industry,
		Size:        c.Size,
	}
}

type CompanyResponse struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Industry    string    `json:"industry,omitempty"`
	Size        string    `json:"size,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *CompanyResponse) FillFromEntity(company entity.Company) {
	c.UUID = company.UUID
	c.Name = company.Name
	c.Description = company.Description
	c.Industry = company.Industry
	c.Size = company.Size
	c.CreatedAt = company.CreatedAt
}