package employee

import (
	"database/sql"
	"fmt"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"employee-crud/models"
	"employee-crud/stores"
)

type store struct{}

func New() stores.Employee {
	return store{}
}

func (s store) GetAll(ctx *gofr.Context) ([]models.Employee, error) {
	rows, err := ctx.DB().QueryContext(ctx, "SELECT * FROM employees")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer func() {
		err = rows.Err()
		ctx.Logger.Errorf("Rows Error: ", err)
	}()

	defer func() {
		err = rows.Close()
		ctx.Logger.Errorf("Error closing rows: ", err)
	}()

	emps := make([]models.Employee, 0)

	for rows.Next() {
		var e models.Employee

		err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.CTC)
		if err != nil {
			return nil, errors.DB{Err: err}
		}

		emps = append(emps, e)
	}

	return emps, nil
}

func (s store) GetByID(ctx *gofr.Context, id int) (models.Employee, error) {
	var res models.Employee

	err := ctx.DB().QueryRowContext(ctx, "SELECT * FROM employees where id=$1", id).Scan(&res.ID, &res.Name, &res.Email, &res.CTC)
	if err != nil {
		switch err {
		case sql.ErrNoRows: {
			return models.Employee{}, errors.EntityNotFound{Entity: "Employee", ID: fmt.Sprint(id)}
		}
		default: 
			return models.Employee{}, errors.DB{Err: err}
		}
	}

	return res, nil
}

func (s store) Update(ctx *gofr.Context, emp models.Employee) (models.Employee, error) {
	_, err := ctx.DB().ExecContext(ctx, "UPDATE employees SET name=$1, email=$2, CTC=$3 WHERE id=$4", emp.Name, emp.Email, emp.CTC, emp.ID)

	if err != nil {
		return models.Employee{}, errors.DB{Err: err}
	}

	return emp, nil
}

func (e store) Create(ctx *gofr.Context, emp models.Employee) (models.Employee, error) {
	var res models.Employee
	err := ctx.DB().QueryRowContext(ctx, "INSERT INTO employees(name, email, CTC) VALUES($1, $2, $3) RETURNING id, name", emp.Name, emp.Email, emp.CTC).Scan(&res.ID, &res.Name, &res.Email, &res.CTC)

	if err != nil {
		fmt.Println("here..")
		return models.Employee{}, errors.DB{Err: err}
	}

	return res, nil
}

func (s store) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.DB().ExecContext(ctx, "DELETE FROM employees where id=$1", id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
