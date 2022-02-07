package employee

import (
	"context"
	"testing"

	"employee-crud/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/stretchr/testify/assert"
)

func TestCoreLayer(t *testing.T) {
	app := gofr.New()

	createTable(app)
	testAddEmployee(t, app)
}

func createTable(app *gofr.Gofr) {
	_, err := app.DB().Exec("DROP TABLE employees;")
	if err != nil {
		return
	}

	_, err = app.DB().Exec("CREATE TABLE customers (id serial primary key,name varchar (50), email varchar(50), CTC NUMERIC (3, 2));")
	if err != nil {
		return
	}
}

func testAddEmployee(t *testing.T, app *gofr.Gofr) {
	tests := []struct {
		desc     string
		employee models.Employee
		err      error
	}{
		{"create succuss test #1", models.Employee{Name: "Test123", Email: "example@example.com", CTC: 0.0}, nil},
		{"create succuss test #2", models.Employee{Name: "Test234", Email: "example2@example.com", CTC: 1.0}, nil},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		store := New()
		resp, err := store.Create(ctx, tc.employee)

		app.Logger.Log(resp)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

