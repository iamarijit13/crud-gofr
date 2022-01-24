package store

import (
	"employee-crud/model"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Store interface {
	Find(ctx *gofr.Context) ([]model.Employee, error)
	FindByID(ctx *gofr.Context, id int) (model.Employee, error)
	Update(ctx *gofr.Context, employee model.Employee) (model.Employee, error)
	Create(ctx *gofr.Context, employee model.Employee) (model.Employee, error)
	Delete(ctx *gofr.Context, id int) error
}