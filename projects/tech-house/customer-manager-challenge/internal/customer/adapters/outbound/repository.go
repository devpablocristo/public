package outbound

import (
	"context"
	"database/sql"
	"errors"

	sqrepo "github.com/devpablocristo/golang-monorepo/pkg/databases/sql/sqlite"
	sqdefs "github.com/devpablocristo/golang-monorepo/pkg/databases/sql/sqlite/defs"

	types "github.com/devpablocristo/golang-monorepo/pkg/types"
	transport "github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/customer/adapters/outbound/transport"
	domain "github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/customer/core/domain"
	ports "github.com/devpablocristo/golang-monorepo/projects/tech-house/customer-manager-challenge/internal/customer/core/ports"
)

type repository struct {
	sqliteRepo sqdefs.Repository
}

func NewRepository() (ports.Repository, error) {
	r, err := sqrepo.Bootstrap()
	if err != nil {
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to bootstrap sqlite",
			err,
		)
	}

	if err := initSchema(r.DB()); err != nil {
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to initialize schema",
			err,
		)
	}

	return &repository{
		sqliteRepo: r,
	}, nil
}

func initSchema(sqliteRepo *sql.DB) error {
	_, err := sqliteRepo.Exec(schema)
	if err != nil {
		return types.NewError(
			types.ErrOperationFailed,
			"failed to create schema",
			err,
		)
	}
	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Customer, error) {
	var models []transport.CustomerDataModel
	err := r.sqliteRepo.SelectContext(ctx, &models, selectAllCustomersQuery)
	if err != nil {
		return nil, types.NewError(
			types.ErrOperationFailed,
			"failed to fetch customers",
			err,
		)
	}

	customers := make([]domain.Customer, len(models))
	for i, model := range models {
		customers[i] = *transport.CustomerDataModelToDomain(&model)
	}
	return customers, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	model, err := scanCustomer(r.sqliteRepo.QueryRowContext(ctx, selectCustomerByIDQuery, id))
	if err != nil {
		return nil, err
	}
	return transport.CustomerDataModelToDomain(model), nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	model, err := scanCustomer(r.sqliteRepo.QueryRowContext(ctx, selectCustomerByEmailQuery, email))
	if err != nil {
		return nil, err
	}
	return transport.CustomerDataModelToDomain(model), nil
}

func (r *repository) Create(ctx context.Context, customer *domain.Customer) error {
	if err := r.validateEmailConflict(ctx, 0, customer.Email); err != nil {
		return err
	}

	model := transport.DomainToCustomerDataModel(customer)
	_, err := r.sqliteRepo.DB().ExecContext(ctx, insertCustomerQuery,
		model.Name, model.LastName, model.Email,
		model.Phone, model.Age, model.BirthDate,
	)
	if err != nil {
		return types.NewError(
			types.ErrOperationFailed,
			"failed to create customer",
			err,
		)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, customer *domain.Customer) error {
	// Verificar que existe el customer primero
	var customers []transport.CustomerDataModel
	err := r.sqliteRepo.SelectContext(ctx, &customers, selectCustomerByIDQuery, customer.ID)
	if err != nil {
		return types.NewError(
			types.ErrOperationFailed,
			"failed to fetch customer",
			err,
		)
	}
	if len(customers) == 0 {
		return types.NewError(
			types.ErrNotFound,
			"customer not found",
			errors.New("no details available"),
		)
	}

	// Validar conflicto de email
	if err := r.validateEmailConflict(ctx, customer.ID, customer.Email); err != nil {
		return err
	}

	// Hacer el update - usando el mismo wrapper
	model := transport.DomainToCustomerDataModel(customer)
	result, err := r.sqliteRepo.ExecContext(ctx, updateCustomerQuery, // <-- Aquí está el cambio
		model.Name, model.LastName, model.Email,
		model.Phone, model.Age, model.BirthDate, model.ID,
	)
	if err != nil {
		return types.NewError(
			types.ErrOperationFailed,
			"failed to update customer",
			err,
		)
	}

	return validateRows(result)
}

// func (r *repository) Update(ctx context.Context, customer *domain.Customer) error {

// 	fmt.Println(customer)
// 	var existingCustomer transport.CustomerDataModel
// 	query := selectCustomerByIDQuery // usar la constante que ya tienes definida
// 	fmt.Printf("Query: %s\nID: %d\n", query, customer.ID)
// 	err := r.sqliteRepo.SelectContext(ctx, &existingCustomer, query, customer.ID)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// 	fmt.Printf("Existing customer: %+v\n", existingCustomer)

// 	if err := r.validateEmailConflict(ctx, customer.ID, customer.Email); err != nil {
// 		return err
// 	}

// 	model := transport.DomainToCustomerDataModel(customer)
// 	result, err := r.sqliteRepo.DB().ExecContext(ctx, updateCustomerQuery,
// 		model.Name, model.LastName, model.Email,
// 		model.Phone, model.Age, model.BirthDate, model.ID,
// 	)
// 	if err != nil {
// 		return types.NewError(
// 			types.ErrOperationFailed,
// 			"failed to update customer",
// 			err,
// 		)
// 	}

// 	return validateRows(result)
// }

func (r *repository) Delete(ctx context.Context, id int64) error {
	result, err := r.sqliteRepo.DB().ExecContext(ctx, deleteCustomerQuery, id)
	if err != nil {
		return types.NewError(
			types.ErrOperationFailed,
			"failed to delete customer",
			err,
		)
	}

	return validateRows(result)
}
