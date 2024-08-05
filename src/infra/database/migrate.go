package database

import (
	"BillingEngine/src/infra/constants"
	"BillingEngine/src/infra/persistence/model"
	"BillingEngine/src/infra/persistence/postgres"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Migrate represent migration schema models
func Migrate(db *gorm.DB) error {
	Loans := model.Loans{}
	Customers := model.Customers{}
	PaymentMethods := model.PaymentMethods{}
	CustomerLoan := model.CustomerLoan{}
	LoanPayment := model.LoanPayment{}

	err := db.AutoMigrate(&Loans, &Customers, &PaymentMethods, &CustomerLoan, &LoanPayment)
	CreateDataMasterLoan(db)
	CreateDataMasterPaymentMethod(db)
	CreateNewCustomerLoan(db)
	CreateOutStandingCustomerLoan(db)
	CreateDelinquentCustomerLoan(db)
	return err
}

func CreateDataMasterLoan(db *gorm.DB) (err error) {
	ctx := context.Background()
	loanRepo := postgres.NewLoansPersistence(db)
	//Create loan master data
	_, err = loanRepo.CreateLoan(ctx, postgres.CreateLoan{
		LoanName:  "Loan Amartha - 5000000",
		LoanPrice: 5000000,
	})
	if err != nil {
		err := errors.New("Error Create Loan")
		return err
	}

	return nil
}

func CreateDataMasterPaymentMethod(db *gorm.DB) (err error) {
	ctx := context.Background()
	paymentMethodRepo := postgres.NewPaymentMethodPersistence(db)
	//create Payment Method
	_, err = paymentMethodRepo.CreatePaymentMethod(ctx, postgres.CreatePaymentMethod{
		MethodName: "Virtual Account",
	})
	if err != nil {
		err := errors.New("Error Get Payment Method")
		return err
	}

	return nil
}

func CreateNewCustomerLoan(db *gorm.DB) (err error) {
	ctx := context.Background()
	loanRepo := postgres.NewLoansPersistence(db)
	custRepo := postgres.NewCustomersPersistence(db)
	paymentMethodRepo := postgres.NewPaymentMethodPersistence(db)
	custLoanRepo := postgres.NewCustomerLoanPersistence(db)
	loanPaymentRepo := postgres.NewLoanPaymentPersistence(db)

	loanData, err := loanRepo.LoansDetailByName(ctx, "Loan Amartha - 5000000")
	if err != nil {
		err := errors.New("Error Create Loan")
		return err
	}

	// create customer master data
	custData, err := custRepo.CreateCustomers(ctx, postgres.CreateCustomer{
		CustomerName:        "Customer Test 1",
		CustomerPhoneNumber: "081511331199",
		RelationName:        "Customer Relation",
		RelationPhoneNumber: "082160259686",
	})

	if err != nil {
		err := errors.New("Error Create Customer")
		return err
	}

	//get Payment Method "Virtual Account"
	paymentMethodData, err := paymentMethodRepo.GetPaymentMethod(ctx, "Virtual Account")
	if err != nil {
		err := errors.New("Error Get Payment Method")
		return err
	}

	createCustLoan, err := custLoanRepo.CreateCustomerLoan(ctx, postgres.CreateCustomerLoanRequest{
		CustomerId:         custData.ID,
		LoanId:             loanData.ID,
		OutstandingBalance: 5500000,
		Status:             constants.STATUS_CUSTOMER_ACTIVE_LOAN,
	})

	if err != nil {
		err := errors.New("Error Create Customer Loan")
		return err
	}

	//create loan payment for 50 weeks
	today := time.Now()
	saveLoanPayment := []postgres.CreateLoanPaymentRequest{}
	for i := 1; i <= 50; i++ {
		saveLoanPayment = append(saveLoanPayment, postgres.CreateLoanPaymentRequest{
			CustomerLoanId: createCustLoan.ID,
			MethodId:       paymentMethodData.ID,
			LoanWeek:       uint(i),
			Status:         constants.STATUS_LOAN_PAYMENT_NOT_PAID,
			DeadlineDate:   today.AddDate(0, 0, i*7),
		})
	}
	_, err = loanPaymentRepo.CreateLoanPayment(ctx, saveLoanPayment)
	if err != nil {
		err := errors.New("Error Create Loan Payment")
		return err
	}

	return err
}

func CreateOutStandingCustomerLoan(db *gorm.DB) (err error) {
	ctx := context.Background()
	custRepo := postgres.NewCustomersPersistence(db)
	loanRepo := postgres.NewLoansPersistence(db)
	paymentMethodRepo := postgres.NewPaymentMethodPersistence(db)
	custLoanRepo := postgres.NewCustomerLoanPersistence(db)
	loanPaymentRepo := postgres.NewLoanPaymentPersistence(db)

	loanData, err := loanRepo.LoansDetailByName(ctx, "Loan Amartha - 5000000")
	if err != nil {
		err := errors.New("Error Create Loan")
		return err
	}
	if err != nil {
		err := errors.New("Error Create Loan")
		return err
	}

	// create customer master data
	custData, err := custRepo.CreateCustomers(ctx, postgres.CreateCustomer{
		CustomerName:        "Customer Test OutStanding",
		CustomerPhoneNumber: "081511331121",
		RelationName:        "Customer Relation OutStanding",
		RelationPhoneNumber: "082160259633",
	})

	if err != nil {
		err := errors.New("Error Create Customer")
		return err
	}

	//get Payment Method "Virtual Account"
	paymentMethodData, err := paymentMethodRepo.GetPaymentMethod(ctx, "Virtual Account")
	if err != nil {
		err := errors.New("Error Get Payment Method")
		return err
	}
	createCustLoan, err := custLoanRepo.CreateCustomerLoan(ctx, postgres.CreateCustomerLoanRequest{
		CustomerId:         custData.ID,
		LoanId:             loanData.ID,
		OutstandingBalance: 0,
		Status:             constants.STATUS_CUSTOMER_LOAN_DONE_PAID,
	})

	if err != nil {
		err := errors.New("Error Create Customer Loan")
		return err
	}

	//create loan payment for 50 weeks
	today := time.Now()
	saveLoanPayment := []postgres.CreateLoanPaymentRequest{}
	for i := 50; i >= 1; i-- {
		saveLoanPayment = append(saveLoanPayment, postgres.CreateLoanPaymentRequest{
			CustomerLoanId: createCustLoan.ID,
			MethodId:       paymentMethodData.ID,
			LoanWeek:       uint(50 - i + 1),
			Status:         constants.STATUS_LOAN_PAYMENT_PAID,
			DeadlineDate:   today.AddDate(0, 0, -(i * 7)),
		})
	}
	_, err = loanPaymentRepo.CreateLoanPayment(ctx, saveLoanPayment)
	if err != nil {
		err := errors.New("Error Create Loan Payment")
		return err
	}

	return err
}

func CreateDelinquentCustomerLoan(db *gorm.DB) (err error) {
	ctx := context.Background()
	custRepo := postgres.NewCustomersPersistence(db)
	loanRepo := postgres.NewLoansPersistence(db)
	paymentMethodRepo := postgres.NewPaymentMethodPersistence(db)
	custLoanRepo := postgres.NewCustomerLoanPersistence(db)
	loanPaymentRepo := postgres.NewLoanPaymentPersistence(db)

	loanData, err := loanRepo.LoansDetailByName(ctx, "Loan Amartha - 5000000")
	if err != nil {
		err := errors.New("Error Create Loan")
		return err
	}
	if err != nil {
		err := errors.New("Error Create Loan")
		return err
	}

	// create customer master data
	custData, err := custRepo.CreateCustomers(ctx, postgres.CreateCustomer{
		CustomerName:        "Customer Test Delinquent",
		CustomerPhoneNumber: "081511331114",
		RelationName:        "Customer Relation Delinquent",
		RelationPhoneNumber: "082160259614",
	})

	if err != nil {
		err := errors.New("Error Create Customer")
		return err
	}

	//get Payment Method "Virtual Account"
	paymentMethodData, err := paymentMethodRepo.GetPaymentMethod(ctx, "Virtual Account")
	if err != nil {
		err := errors.New("Error Get Payment Method")
		return err
	}
	createCustLoan, err := custLoanRepo.CreateCustomerLoan(ctx, postgres.CreateCustomerLoanRequest{
		CustomerId:         custData.ID,
		LoanId:             loanData.ID,
		OutstandingBalance: 5500000 - 4950000,
		Status:             constants.STATUS_CUSTOMER_ACTIVE_LOAN,
	})

	if err != nil {
		err := errors.New("Error Create Customer Loan")
		return err
	}

	//create loan payment for 50 weeks
	today := time.Now()
	saveLoanPayment := []postgres.CreateLoanPaymentRequest{}
	for i := 48; i >= -2; i-- {
		week := 48 - i + 1
		if i < 0 {
			week = 48 - i
		}
		status := constants.STATUS_LOAN_PAYMENT_PAID
		if week >= 46 {
			status = constants.STATUS_LOAN_PAYMENT_NOT_PAID
		}
		saveLoanPayment = append(saveLoanPayment, postgres.CreateLoanPaymentRequest{
			CustomerLoanId: createCustLoan.ID,
			MethodId:       paymentMethodData.ID,
			LoanWeek:       uint(week),
			Status:         status,
			DeadlineDate:   today.AddDate(0, 0, -(i * 7)),
		})
	}
	_, err = loanPaymentRepo.CreateLoanPayment(ctx, saveLoanPayment)
	if err != nil {
		err := errors.New("Error Create Loan Payment")
		return err
	}

	return err
}
