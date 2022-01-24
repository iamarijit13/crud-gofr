package handler

import (
	"bytes"
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"

	"employee-crud/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"

	"github.com/stretchr/testify/assert"
)

type mockStore struct {}

func (m mockStore) Find(c *gofr.Context) ([]model.Employee, error){
	p := c.Param("mock")

	if p == "success" {
		return nil, nil
	}

	return nil, errors.Error("Error fetching employees.")
}

func (m mockStore) FindByID(c *gofr.Context, id int) (model.Employee, error) {
	if (id == 1) {
		return model.Employee{ID: 1, Name: "First_Name Second_Name", Email: "example@exaple.com", CTC: 0.0}, nil
	}

	return model.Employee{}, errors.EntityNotFound{Entity: "employee", ID: fmt.Sprint(id)}
}

func (m mockStore) Update(c *gofr.Context, customer model.Employee) (model.Employee, error) {
	if customer.Name == "First_Name Second_Name" {
		return model.Employee{}, nil
	}

	return model.Employee{}, errors.Error("error updating employee.")
}

func (m mockStore) Create(c *gofr.Context, emp model.Employee) (model.Employee, error) {
	switch emp.Name {
	case "First_Name Second_Name":
		return model.Employee{ID: 1, Name: emp.Name, Email: emp.Email, CTC: emp.CTC}, nil
	case "mock body error":
		return model.Employee{}, errors.InvalidParam{Param: []string{"body"}}
	case `{"id":1}`:
		return model.Employee{}, errors.InvalidParam{Param: []string{"id"}}
	}

	return model.Employee{}, errors.Error("error adding new employee")
}

func (m mockStore) Delete(c *gofr.Context, id int) error {
	if c.PathParam("id") == "123" {
		return nil
	}

	return errors.Error("error deleting employee.")
}

func TestModel_AddEmployee(t *testing.T) {
	h := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		desc string
		body []byte
		err error
	}{
		{"create with invalid id", []byte(`{"id":1}`), errors.InvalidParam{Param: []string{"id"}}},
		{"create succuss", []byte(`{"name":"First_Name Second_Name"}`), nil},
		{"create invalid body", []byte(`mock body error`), errors.InvalidParam{Param: []string{"body"}}},
		{"create error", []byte(`{"name":"creation error"}`), errors.Error("error adding new employee")},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, k)

		_, err := h.Create(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_UpdateCustomer(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		body []byte
		err  error
		id   string
	}{
		{"missing id", nil, errors.MissingParam{Param: []string{"id"}}, ""},
		{"invalid id", nil, errors.InvalidParam{Param: []string{"id"}}, "abc123"},
		{"invalid body", []byte(`{`), errors.InvalidParam{Param: []string{"body"}}, "123"},
		{"update succuss", []byte(`{"name":"First_Name Second_Name"}`), nil, "123"},
		{"update error", []byte(`{"name":"creation error"}`), errors.Error("error updating employee."), "123"},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		_, err := h.Update(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_GetCustomerById(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		id   string
		resp interface{}
		err  error
	}{
		{"get by id succuss", "1", model.Employee{ID: 1, Name: "First_Name Second_Name", Email: "example@exaple.com", CTC: 0.0}, nil},
		{"invalid id", "absd123", nil, errors.InvalidParam{Param: []string{"id"}}},
		{"missing id", "", nil, errors.MissingParam{Param: []string{"id"}}},
		{"id not found", "2", nil, errors.EntityNotFound{Entity: "Employee", ID: "2"}},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		resp, err := h.FindByID(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_DeleteCustomer(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"delete succuss", "123", nil},
		{"delete fail", "1234", errors.Error("error deleting employee.")},
		{"invalid id", "absd123", errors.InvalidParam{Param: []string{"id"}}},
		{"missing id", "", errors.MissingParam{Param: []string{"id"}}},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		_, err := h.Delete(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_GetCustomers(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc         string
		mockParamStr string
		err          error
	}{
		{"get success", "mock=success", nil},
		{"get fail", "", errors.Error("Error fetching employees.")},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://dummy?"+tc.mockParamStr, nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := h.Find(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}