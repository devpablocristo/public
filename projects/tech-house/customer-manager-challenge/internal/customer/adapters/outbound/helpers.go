package outbound

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	types "github.com/devpablocristo/customer-manager/pkg/types"
	transport "github.com/devpablocristo/customer-manager/projects/tech-house/customer-manager-challenge/internal/customer/adapters/outbound/transport"
)

const (
	// Schema definition
	schema = `
        CREATE TABLE IF NOT EXISTS customers (
            id          INTEGER PRIMARY KEY AUTOINCREMENT,
            name        TEXT NOT NULL,
            last_name   TEXT NOT NULL,
            email       TEXT NOT NULL UNIQUE,
            phone       TEXT NOT NULL,
            age         INTEGER NOT NULL,
            birth_date  DATETIME NOT NULL
        );
    `

	// Base select query
	selectAllCustomersQuery = `
        SELECT  id, 
                name, 
                last_name, 
                email, 
                phone, 
                age, 
                birth_date
        FROM    customers
    `

	// Select queries
	selectCustomerByIDQuery    = selectAllCustomersQuery + ` WHERE id = ?`
	selectCustomerByEmailQuery = selectAllCustomersQuery + ` WHERE email = ?`

	// Insert query
	insertCustomerQuery = `
        INSERT INTO customers (
            name,
            last_name,
            email,
            phone,
            age,
            birth_date
        ) VALUES (?, ?, ?, ?, ?, ?)
    `

	// Update query
	updateCustomerQuery = `
        UPDATE  customers 
        SET     name = ?, 
                last_name = ?, 
                email = ?, 
                phone = ?, 
                age = ?, 
                birth_date = ?
        WHERE   id = ?
    `

	// Delete query
	deleteCustomerQuery = `DELETE FROM customers WHERE id = ?`
)

func validateRows(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return types.NewError(
			types.ErrOperationFailed,
			"failed to get affected rows",
			err,
		)
	}

	if rows == 0 {
		return types.NewError(
			types.ErrNotFound,
			"No rows were affected",
			errors.New("no details available"),
		)
	}
	return nil
}

func scanCustomer(row *sql.Row) (*transport.CustomerDataModel, error) {
	var model transport.CustomerDataModel
	err := row.Scan(
		&model.ID, &model.Name, &model.LastName, &model.Email,
		&model.Phone, &model.Age, &model.BirthDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewError(
				types.ErrNotFound,
				"customer not found",
				err,
			)
		}
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to fetch customer",
			err,
		)
	}
	return &model, nil
}

func (r *repository) validateEmailConflict(ctx context.Context, customerID int64, email string) error {
	existing, err := r.GetByEmail(ctx, email)
	if err != nil {
		if !types.IsNotFound(err) {
			return err
		}
		return nil
	}
	if existing != nil && existing.ID != customerID {
		return types.NewError(
			types.ErrConflict,
			fmt.Sprintf("email %s is already in use by another customer", email),
			errors.New("no details available"),
		)
	}
	return nil
}
