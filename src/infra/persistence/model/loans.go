package model

type Loans struct {
	Base
	LoanName     string         `gorm:"column:loan_name" json:"loanName"`
	LoanPrice    float64        `gorm:"column:loan_price" json:"loanPrice"`
	CustomerLoan []CustomerLoan `gorm:"foreignKey:LoanId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
