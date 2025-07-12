package viewmodel

import (
	"time"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type CompanyRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=255"`
	Industry  string `json:"industry,omitempty"`
	Size      string `json:"size,omitempty" validate:"omitempty,oneof=small medium large enterprise"`
	Role      string `json:"role,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

func (c CompanyRequest) ToEntity() entity.Company {
	return entity.Company{
		Name:      c.Name,
		Industry:  c.Industry,
		Size:      c.Size,
		Role:      c.Role,
		IsDefault: c.IsDefault,
	}
}

type CompanyResponse struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Industry  string    `json:"industry,omitempty"`
	Size      string    `json:"size,omitempty"`
	Role      string    `json:"role,omitempty"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *CompanyResponse) FillFromEntity(company entity.Company) {
	c.UUID = company.UUID
	c.Name = company.Name
	c.Industry = company.Industry
	c.Size = company.Size
	c.Role = company.Role
	c.IsDefault = company.IsDefault
	c.CreatedAt = company.CreatedAt
}