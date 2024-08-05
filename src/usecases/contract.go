package usecases

type AllUseCases struct {
	LoanPayment  ILoanPaymentUsecase
	CustomerLoan ICustomerLoanUsecase
}
