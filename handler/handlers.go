package handler

import (
	"employee-crud/model"
	"employee-crud/store"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type handler struct {
	store store.Store
}

// Struct use for response.
type response struct {
	Employees []model.Employee
}

// Factory function, entry point for store.
func New(s store.Store) handler {
	return handler{store: s}
}

func (h handler) Find(c *gofr.Context) (interface{}, error) {
	res, err := h.store.Find(c)

	if err != nil {
		return nil, err
	}

	re := response{Employees: res}

	return re, nil
}

func (h handler) FindByID(c *gofr.Context) (interface{}, error) {
	in := c.PathParam("id")

	if in == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(in)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	res, err := h.store.FindByID(c, id)

	if err != nil {
		return nil, errors.EntityNotFound{Entity: "Employee", ID: in}
	}

	return res, nil
}

func (h handler) Create(c *gofr.Context) (interface{}, error) {
	var emp model.Employee

	if err := c.Bind(&emp); err != nil {
		c.Logger.Errorf("Error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	if emp.ID != 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	res, err := h.store.Create(c, emp)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h handler) Update(c *gofr.Context) (interface{}, error) {
	in := c.PathParam("id")

	if in == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(in)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var emp model.Employee
	if err = c.Bind(&emp); err != nil {
		c.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	emp.ID = id

	res, err := h.store.Update(c, emp)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h handler) Delete(c *gofr.Context) (interface{}, error) {
	in := c.PathParam("id")
	if in == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(in)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	if err := h.store.Delete(c, id); err != nil {
		return nil, err
	}

	return "Deleted Successfully", nil
}
