package postgres

import (
	"BillingEngine/src/infra/persistence/model"
	"context"

	"gorm.io/gorm"
)

type CreateLoan struct {
	LoanName  string
	LoanPrice float64
}

type ILoansRepository interface {
	LoansDetailByName(ctx context.Context, name string) (resp *model.Loans, err error)
	CreateLoan(ctx context.Context, req CreateLoan) (resp *model.Loans, err error)
}

type LoansPersistence struct {
	dBConn *gorm.DB
}

// NewLoansPersistence ...
func NewLoansPersistence(db *gorm.DB) ILoansRepository {
	return &LoansPersistence{
		dBConn: db,
	}
}

func (u *LoansPersistence) LoansDetailByName(ctx context.Context, name string) (resp *model.Loans, err error) {
	db := u.dBConn.WithContext(ctx)
	err = db.Find(&resp, "loan_name = ?", name).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (u *LoansPersistence) CreateLoan(ctx context.Context, req CreateLoan) (resp *model.Loans, err error) {
	create := model.Loans{
		LoanName:  req.LoanName,
		LoanPrice: req.LoanPrice,
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
