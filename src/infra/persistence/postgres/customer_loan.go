package postgres

import (
	"BillingEngine/src/infra/persistence/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type CreateCustomerLoanRequest struct {
	CustomerId         uint
	LoanId             uint
	OutstandingBalance float64
	Status             string
}

type UpdateCustomerLoanRequest struct {
	CustomerLoanId int
	TotalLoan      float64
	Status         string
}

type GetCustomerLoanList struct {
	CustomerId int
	LoanId     int
	Status     string
}

type ICustomerLoanRepository interface {
	GetCustomerLoanOutstanding(ctx context.Context, req GetCustomerLoanList) (resp []model.CustomerLoanResp, err error)
	GetCustomerLoanDelinquent(ctx context.Context, req GetCustomerLoanList) (resp []model.CustomerLoanResp, err error)
	CustomerLoanDetail(ctx context.Context, id int) (resp *model.CustomerLoan, err error)
	CreateCustomerLoan(ctx context.Context, req CreateCustomerLoanRequest) (resp *model.CustomerLoan, err error)
	UpdateCustomerLoan(ctx context.Context, req UpdateCustomerLoanRequest) (resp *model.CustomerLoan, err error)
}

type customerLoanPersistence struct {
	dBConn *gorm.DB
}

// NewCustomerLoanPersistence ...
func NewCustomerLoanPersistence(db *gorm.DB) ICustomerLoanRepository {
	return &customerLoanPersistence{
		dBConn: db,
	}
}

func (d *customerLoanPersistence) CustomerLoanDetail(ctx context.Context, id int) (resp *model.CustomerLoan, err error) {
	db := d.dBConn.WithContext(ctx)
	err = db.Find(&resp, "id = ?", id).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (d *customerLoanPersistence) GetCustomerLoanOutstanding(ctx context.Context, req GetCustomerLoanList) (resp []model.CustomerLoanResp, err error) {
	db := d.dBConn.WithContext(ctx)
	query := fmt.Sprintf(`SELECT customer_loans.customer_id, 
		customer_loans.outstanding_balance, 
		customer_loans.status, 
		customers.customer_name, 
		customers.customer_phone_number
	FROM customer_loans
	join loans on loans.id = customer_loans.loan_id
	JOIN customers on customers.id = customer_loans.customer_id
	join (
		select lp.customer_loan_id from public.loan_payments lp 
		where lp.deadline_date::date < now()::date and status <> 'Paid'
	) cl on cl.customer_loan_id <> customer_loans.id 
	WHERE loans.loan_name = 'Loan Amartha - 5000000' AND customer_loans.deleted_at IS NULL
	GROUP BY customers.id, customer_loans.id  `)
	err = db.Raw(query).
		Scan(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (d *customerLoanPersistence) GetCustomerLoanDelinquent(ctx context.Context, req GetCustomerLoanList) (resp []model.CustomerLoanResp, err error) {
	db := d.dBConn.WithContext(ctx)
	query := fmt.Sprintf(`SELECT customer_loans.customer_id, 
							customer_loans.outstanding_balance, 
							customer_loans.status, 
							customers.customer_name, 
							customers.customer_phone_number
						FROM customer_loans
						join loans on loans.id = customer_loans.loan_id
						JOIN customers on customers.id = customer_loans.customer_id
						join public.loan_payments on loan_payments.customer_loan_id = customer_loans.id  
						join (
							select lp.customer_loan_id,   COUNT(*) as total from public.loan_payments lp 
							where lp.deadline_date::date < now()::date and status = 'Not Paid'
							group by lp.customer_loan_id 
						) cl on cl.customer_loan_id = customer_loans.id 
						WHERE loans.loan_name = 'Loan Amartha - 5000000' AND customer_loans.deleted_at IS null
						and cl.total >=2
						group by customers.id, customer_loans.id `)
	err = db.Raw(query).
		Scan(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (d *customerLoanPersistence) CreateCustomerLoan(ctx context.Context, req CreateCustomerLoanRequest) (resp *model.CustomerLoan, err error) {
	create := model.CustomerLoan{
		CustomerId:         req.CustomerId,
		LoanId:             req.LoanId,
		OutstandingBalance: req.OutstandingBalance,
		Status:             req.Status,
	}
	trx := d.dBConn.WithContext(ctx).Begin()
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

func (d *customerLoanPersistence) UpdateCustomerLoan(ctx context.Context, req UpdateCustomerLoanRequest) (resp *model.CustomerLoan, err error) {
	trx := d.dBConn.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	custLoan := model.CustomerLoan{}
	err = trx.Model(&custLoan).Where("id = ?", req.CustomerLoanId).Updates(map[string]interface{}{
		"outstanding_balance": req.TotalLoan,
		"status":              req.Status,
	}).Error
	if err != nil {
		return resp, err
	}

	err = trx.Commit().Error
	if err != nil {
		return resp, err
	}
	return resp, nil
}
