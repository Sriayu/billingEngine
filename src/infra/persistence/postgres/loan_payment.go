package postgres

import (
	"BillingEngine/src/infra/persistence/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type CreateLoanPaymentRequest struct {
	CustomerLoanId uint
	MethodId       uint
	LoanWeek       uint
	LoanPaidPrice  float64
	DeadlineDate   time.Time
	Status         string
}

type GetLoanPaymentRequest struct {
	CustomerLoanId uint
	LoanWeek       int
}

type UpdateLoanPaymentRequest struct {
	CustomerLoanId uint
	LoanWeek       int
	Status         string
}

type ILoanPaymentRepository interface {
	CreateLoanPayment(ctx context.Context, req []CreateLoanPaymentRequest) (resp []*model.LoanPayment, err error)
	GetLoanPayment(ctx context.Context, req GetLoanPaymentRequest) (resp *model.LoanPayment, err error)
	GetLatestLoanPayment(ctx context.Context, customerLoanId uint) (resp *model.LoanPayment, err error)
	UpdateLoanPayment(ctx context.Context, req UpdateLoanPaymentRequest) (err error)
}

type loanPaymentPersistence struct {
	dBConn *gorm.DB
}

// NewLoanPaymentPersistence ...
func NewLoanPaymentPersistence(db *gorm.DB) ILoanPaymentRepository {
	return &loanPaymentPersistence{
		dBConn: db,
	}
}
func (c *loanPaymentPersistence) GetLoanPayment(ctx context.Context, req GetLoanPaymentRequest) (resp *model.LoanPayment, err error) {
	db := c.dBConn.WithContext(ctx)
	err = db.Find(&resp, "customer_loan_id = ? AND loan_week = ?", req.CustomerLoanId, req.LoanWeek).Error
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *loanPaymentPersistence) GetLatestLoanPayment(ctx context.Context, customerLoanId uint) (resp *model.LoanPayment, err error) {
	db := c.dBConn.WithContext(ctx)
	err = db.Find(&resp, "customer_loan_id = ? AND status = 'Not Paid'", customerLoanId).Error
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *loanPaymentPersistence) UpdateLoanPayment(ctx context.Context, req UpdateLoanPaymentRequest) (err error) {
	trx := c.dBConn.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	loanPayment := model.LoanPayment{}
	err = trx.Model(&loanPayment).Where("customer_loan_id = ? AND loan_week = ?", req.CustomerLoanId, req.LoanWeek).Updates(map[string]interface{}{
		"status": req.Status,
	}).Error
	if err != nil {
		return err
	}

	err = trx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (c *loanPaymentPersistence) CreateLoanPayment(ctx context.Context, req []CreateLoanPaymentRequest) (resp []*model.LoanPayment, err error) {
	create := []*model.LoanPayment{}
	for _, lp := range req {
		create = append(create, &model.LoanPayment{
			CustomerLoanId: lp.CustomerLoanId,
			MethodId:       lp.MethodId,
			LoanWeek:       lp.LoanWeek,
			LoanPaidPrice:  lp.LoanPaidPrice,
			Status:         lp.Status,
			DeadlineDate:   lp.DeadlineDate,
		})
	}
	trx := c.dBConn.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	err = trx.Create(&create).Error
	if err != nil {
		return resp, err
	}
	err = trx.Commit().Error
	if err != nil {
		return resp, err
	}

	return create, nil
}
