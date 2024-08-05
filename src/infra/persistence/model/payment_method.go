package model

type PaymentMethods struct {
	Base
	MethodName  string        `gorm:"column:method_name" json:"methodName"`
	LoanPayment []LoanPayment `gorm:"foreignKey:MethodId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
