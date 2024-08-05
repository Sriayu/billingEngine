package postgres

import (
	"BillingEngine/src/infra/persistence/model"
	"context"

	"gorm.io/gorm"
)

type CreateCustomer struct {
	CustomerName        string
	CustomerPhoneNumber string
	RelationName        string
	RelationPhoneNumber string
}

type ICustomersRepository interface {
	CustomersDetail(ctx context.Context, id int) (resp *model.Customers, err error)
	CreateCustomers(ctx context.Context, req CreateCustomer) (resp *model.Customers, err error)
	GetCustomerByName(ctx context.Context, name string) (resp *model.Customers, err error)
}

type CustomersPersistence struct {
	dBConn *gorm.DB
}

// NewCustomersPersistence ...
func NewCustomersPersistence(db *gorm.DB) ICustomersRepository {
	return &CustomersPersistence{
		dBConn: db,
	}
}

func (c *CustomersPersistence) CustomersDetail(ctx context.Context, id int) (resp *model.Customers, err error) {
	db := c.dBConn.WithContext(ctx)
	err = db.Find(&resp, "id = ?", id).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (u *CustomersPersistence) GetCustomerByName(ctx context.Context, name string) (resp *model.Customers, err error) {
	db := u.dBConn.WithContext(ctx)
	err = db.Find(&resp, "customer_name = ?", name).Error
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *CustomersPersistence) CreateCustomers(ctx context.Context, req CreateCustomer) (resp *model.Customers, err error) {
	create := model.Customers{
		CustomerName:        req.CustomerName,
		CustomerPhoneNumber: req.CustomerPhoneNumber,
		RelationName:        req.RelationName,
		RelationPhoneNumber: req.RelationPhoneNumber,
	}

	trx := c.dBConn.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	result := trx.Create(&create)
	err = result.Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	err = trx.Commit().Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}
