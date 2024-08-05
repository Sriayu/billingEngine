package model

import "time"

type LoanPayment struct {
	Base
	CustomerLoanId uint      `gorm:"column:customer_loan_id;default:0;index:idx_customer_loan_id" json:"customerLoanId"`
	MethodId       uint      `gorm:"column:method_id;default:0;index:idx_method_id" json:"methodId"`
	LoanWeek       uint      `gorm:"column:loan_week" json:"loanWeek"`
	LoanPaidPrice  float64   `gorm:"column:loan_paid_price;default:0" json:"loanPaidPrice"`
	Status         string    `gorm:"column:status" json:"status"`
	DeadlineDate   time.Time `gorm:"column:deadline_date" json:"deadlineDate"`
}
