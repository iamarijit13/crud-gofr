package stores

import (
	"employee-crud/models"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Employee interface {
	GetAll(ctx *gofr.Context) ([]models.Employee, error)
	GetByID(ctx *gofr.Context, id int) (models.Employee, error)
	Update(ctx *gofr.Context, employee models.Employee) (models.Employee, error)
	Create(ctx *gofr.Context, employee models.Employee) (models.Employee, error)
	Delete(ctx *gofr.Context, id int) error
}