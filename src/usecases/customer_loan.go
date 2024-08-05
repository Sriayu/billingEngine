package usecases

import (
	dto "BillingEngine/src/handler/requests/data_transfer_object"
	"BillingEngine/src/infra/constants"
	"BillingEngine/src/infra/persistence/model"
	"BillingEngine/src/infra/persistence/postgres"
	"context"
	"errors"
)

type ICustomerLoanUsecase interface {
	GetCustomerLoanByStatus(ctx context.Context, do dto.ICustomerLoanRequest) (resp []model.CustomerLoanResp, err error)
	UpdateCustomerLoan(ctx context.Context, do dto.ICustomerLoanRequest) (resp *model.CustomerLoan, err error)
}

type customerLoanUsecase struct {
	custLoanRepo    postgres.ICustomerLoanRepository
	loanPaymentRepo postgres.ILoanPaymentRepository
}

// NewCustomerLoanUsecase ...
func NewCustomerLoanUsecase(
	cl postgres.ICustomerLoanRepository,
	lp postgres.ILoanPaymentRepository,
) ICustomerLoanUsecase {
	return &customerLoanUsecase{
		custLoanRepo:    cl,
		loanPaymentRepo: lp,
	}
}

func (d *customerLoanUsecase) GetCustomerLoanByStatus(ctx context.Context, do dto.ICustomerLoanRequest) (resp []model.CustomerLoanResp, err error) {
	custLoan, ok := do.(*dto.CustomerLoanListRequestDto)
	if !ok {
		err := errors.New("type assertion failed to dto.CustomerLoanDelinquentListRequestDto")
		return resp, err
	}

	custLoanData := []model.CustomerLoanResp{}
	if custLoan.Status == constants.REQ_STATUS_CUSTOMER_LOAN_OUTSTANDING {
		custLoanData, err = d.custLoanRepo.GetCustomerLoanOutstanding(ctx, postgres.GetCustomerLoanList{
			CustomerId: custLoan.CustomerId,
			LoanId:     custLoan.LoanId,
			Status:     custLoan.Status,
		})
		if err != nil {
			err := errors.New("Error Get Customer Loan data")
			return resp, err
		}
	} else if custLoan.Status == constants.REQ_STATUS_CUSTOMER_LOAN_DELINQUENT {
		custLoanData, err = d.custLoanRepo.GetCustomerLoanDelinquent(ctx, postgres.GetCustomerLoanList{
			CustomerId: custLoan.CustomerId,
			LoanId:     custLoan.LoanId,
			Status:     custLoan.Status,
		})
		if err != nil {
			err := errors.New("Error Get Customer Loan data")
			return resp, err
		}
	}

	return custLoanData, err
}

func (d *customerLoanUsecase) UpdateCustomerLoan(ctx context.Context, do dto.ICustomerLoanRequest) (resp *model.CustomerLoan, err error) {
	custLoan, ok := do.(*dto.UpdateCustomerLoanRequestDto)
	if !ok {
		err := errors.New("type assertion failed to dto.UpdateCustomerLoanRequestDto")
		return resp, err
	}

	loanPay, err := d.loanPaymentRepo.GetLoanPayment(ctx, postgres.GetLoanPaymentRequest{
		CustomerLoanId: uint(custLoan.CustomerLoanId),
	})

	updateCUstLoan, err := d.custLoanRepo.UpdateCustomerLoan(ctx, postgres.UpdateCustomerLoanRequest{
		CustomerLoanId: custLoan.CustomerLoanId,
		TotalLoan:      custLoan.TotalLoanPaidId,
		Status:         loanPay.Status,
	})

	return updateCUstLoan, err
}
