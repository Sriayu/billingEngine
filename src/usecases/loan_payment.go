package usecases

import (
	dto "BillingEngine/src/handler/requests/data_transfer_object"
	"BillingEngine/src/infra/constants"
	"BillingEngine/src/infra/persistence/postgres"
	"context"
	"errors"
)

type ILoanPaymentUsecase interface {
	PayLoanPayment(ctx context.Context, do dto.ILoanPaymentRequest) (err error)
}

type loanPaymentUsecase struct {
	loanPaymentRepo postgres.ILoanPaymentRepository
	custLoanRepo    postgres.ICustomerLoanRepository
}

// NewLoanPaymentUsecase ...
func NewLoanPaymentUsecase(lp postgres.ILoanPaymentRepository, cl postgres.ICustomerLoanRepository) ILoanPaymentUsecase {
	return &loanPaymentUsecase{
		loanPaymentRepo: lp,
		custLoanRepo:    cl,
	}
}

func (u *loanPaymentUsecase) PayLoanPayment(ctx context.Context, do dto.ILoanPaymentRequest) (err error) {
	customerLoan, ok := do.(*dto.PayCustomerLoanRequestDto)
	if !ok {
		err = errors.New("type assertion failed to dto.UsersRegisterRequestDto")
		return err
	}

	if int(customerLoan.TotalPricePaid)%110000 != 0 {
		err = errors.New("Borrower can only pay the exact 1100000")
		return err
	}
	weekPaid := int(customerLoan.TotalPricePaid) / 110000
	var (
		latestWeek     int = 0
		statusCustLoan     = constants.STATUS_CUSTOMER_ACTIVE_LOAN
	)
	for i := 1; i <= weekPaid; i++ {
		getLatestPayment, err := u.loanPaymentRepo.GetLatestLoanPayment(ctx, uint(customerLoan.CustomerLoanId))
		if err != nil {
			return err
		}

		err = u.loanPaymentRepo.UpdateLoanPayment(ctx, postgres.UpdateLoanPaymentRequest{
			CustomerLoanId: getLatestPayment.CustomerLoanId,
			LoanWeek:       int(getLatestPayment.LoanWeek),
			Status:         "Paid",
		})
		if err != nil {
			return err
		}

		latestWeek = int(getLatestPayment.LoanWeek)
	}

	custLoanData, err := u.custLoanRepo.CustomerLoanDetail(ctx, customerLoan.CustomerLoanId)
	if err != nil {
		return err
	}
	if latestWeek == 50 {
		statusCustLoan = constants.STATUS_CUSTOMER_LOAN_DONE_PAID
	}
	outstandingBalance := custLoanData.OutstandingBalance - customerLoan.TotalPricePaid
	_, err = u.custLoanRepo.UpdateCustomerLoan(ctx, postgres.UpdateCustomerLoanRequest{
		CustomerLoanId: customerLoan.CustomerLoanId,
		TotalLoan:      outstandingBalance,
		Status:         statusCustLoan,
	})

	if err != nil {
		return err
	}

	return nil
}
