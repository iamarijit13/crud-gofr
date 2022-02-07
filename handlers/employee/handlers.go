package employee

import (
	"fmt"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"employee-crud/models"
	"employee-crud/stores"
)

type handler struct {
	store stores.Employee
}

// Factory function, entry point for store.
func New(s stores.Employee) handler {
	return handler{store: s}
}

func (h handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	res, err := h.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h handler) GetByID(ctx *gofr.Context) (interface{}, error) {
	idParam := ctx.PathParam("id")

	if idParam == "" {
		return nil, errors.MissingParam{Param: []string{"in"}}
	}

	id, err := strconv.Atoi(idParam)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	res, err := h.store.GetByID(ctx, id)

	if err != nil {
		return nil, errors.EntityNotFound{Entity: "Employee", ID: idParam}
	}

	return res, nil
}

func (h handler) Create(ctx *gofr.Context) (interface{}, error) {
	var emp models.Employee

	if err := ctx.Bind(&emp); err != nil {
		ctx.Logger.Errorf("Error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	fmt.Println("on handler...")
	res, err := h.store.Create(ctx, emp)
	
	if err != nil {
		fmt.Println("on handler")
		return nil, err
	}

	return res, nil
}

func (h handler) Update(ctx *gofr.Context) (interface{}, error) {
	idParam := ctx.PathParam("id")

	if idParam == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var emp models.Employee
	if err = ctx.Bind(&emp); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	emp.ID = id

	res, err := h.store.Update(ctx, emp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h handler) Delete(ctx *gofr.Context) (interface{}, error) {
	idParam := ctx.PathParam("id")
	if idParam == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(idParam)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	if err := h.store.Delete(ctx, id); err != nil {
		return nil, err
	}

	return "Deleted Successfully", nil
}
