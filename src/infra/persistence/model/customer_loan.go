package model

type CustomerLoan struct {
	Base
	CustomerId         uint          `gorm:"column:customer_id;default:0;index:idx_customer_id" json:"customerId"`
	LoanId             uint          `gorm:"column:loan_id;default:0;index:idx_loan_id" json:"loanId"`
	OutstandingBalance float64       `gorm:"column:outstanding_balance" json:"outstandingBalance"`
	Status             string        `gorm:"column:status" json:"status"`
	LoanPayment        []LoanPayment `gorm:"foreignKey:CustomerLoanId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CustomerLoanResp struct {
	CustomerId          uint    `gorm:"column:customer_id;default:0;index:idx_customer_id" json:"customerId"`
	OutstandingBalance  float64 `gorm:"column:outstanding_balance" json:"outstandingBalance"`
	Status              string  `gorm:"column:status" json:"status"`
	CustomerName        string  `gorm:"column:customer_name" json:"customerName"`
	CustomerPhoneNumber string  `gorm:"column:customer_phone_number" json:"customerPhoneNumber"`
}
