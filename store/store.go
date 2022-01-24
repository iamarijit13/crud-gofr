package store

import (
	"database/sql"
	"fmt"
	"employee-crud/model"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type employee struct {}

func New() Store {
	return employee{}
}

func (e employee) Find(c *gofr.Context) ([]model.Employee, error) {
	rows, err := c.DB().QueryContext(c, "SELECT * FROM employees")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	emps := make([]model.Employee, 0)

	for rows.Next() {
		var e model.Employee

		err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.CTC)
		if err != nil {
			return nil, err
		}

		emps = append(emps, e)
	}

	return emps, nil
}

func (e employee) FindByID(c *gofr.Context, id int) (model.Employee, error) {
	var res model.Employee

	err := c.DB().QueryRowContext(c, "SELECT * FROM employees where id=$1", id).Scan(&res.ID, &res.Name, &res.Email, &res.CTC)

	if err == sql.ErrNoRows {
		return model.Employee{}, errors.EntityNotFound{Entity: "Employee", ID: fmt.Sprint(id)}
	}

	return res, nil
}

func (e employee) Update(c *gofr.Context, emp model.Employee) (model.Employee, error) {
	_, err := c.DB().ExecContext(c, "UPDATE employees SET name=$1, email=$2, CTC=$3 WHERE id=$4", emp.Name, emp.Email, emp.CTC, emp.ID)

	if err != nil {
		return model.Employee{}, errors.DB{Err: err}
	}

	return emp, nil
}

func (e employee) Create(c *gofr.Context, emp model.Employee) (model.Employee, error) {
	var res model.Employee

	err := c.DB().QueryRowContext(c, "INSERT INTO employees(name, email, CTC) VALUES($1, $2, $3) RETURNING id, name", emp.Name, emp.Email, emp.CTC).Scan(&res.ID, &res.Name, &res.Email, &res.CTC)

	if err != nil {
		return model.Employee{}, errors.DB{Err: err}
	}

	return res, nil
}

func (e employee) Delete(c *gofr.Context, id int) error {
	_, err := c.DB().ExecContext(c, "DELETE FROM employees where id=$1", id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}