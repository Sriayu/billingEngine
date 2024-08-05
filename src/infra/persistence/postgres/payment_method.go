package postgres

import (
	"BillingEngine/src/infra/persistence/model"
	"context"

	"gorm.io/gorm"
)

type CreatePaymentMethod struct {
	MethodName string
}

type IPaymentMethodRepository interface {
	CreatePaymentMethod(ctx context.Context, req CreatePaymentMethod) (resp *model.PaymentMethods, err error)
	GetPaymentMethod(ctx context.Context, methodName string) (resp *model.PaymentMethods, err error)
}

type paymentMethodPersistence struct {
	dBConn *gorm.DB
}

// NewPaymentMethodPersistence ...
func NewPaymentMethodPersistence(db *gorm.DB) IPaymentMethodRepository {
	return &paymentMethodPersistence{
		dBConn: db,
	}
}

func (u *paymentMethodPersistence) CreatePaymentMethod(ctx context.Context, req CreatePaymentMethod) (resp *model.PaymentMethods, err error) {
	// create := model.PaymentMethods{
	// 	MethodName: "Transfer",
	// }
	create := model.PaymentMethods{
		MethodName: req.MethodName,
	}
	trx := u.dBConn.WithContext(ctx).Begin()

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

func (u *paymentMethodPersistence) GetPaymentMethod(ctx context.Context, methodName string) (resp *model.PaymentMethods, err error) {
	db := u.dBConn.WithContext(ctx)
	err = db.Find(&resp, "method_name = ?", methodName).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}
