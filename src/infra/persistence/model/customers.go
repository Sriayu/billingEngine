package model

type Customers struct {
	Base
	CustomerName        string         `gorm:"column:customer_name" json:"customerName"`
	CustomerPhoneNumber string         `gorm:"column:customer_phone_number" json:"customerPhoneNumber"`
	RelationName        string         `gorm:"column:relation_name" json:"relationName"`
	RelationPhoneNumber string         `gorm:"column:relation_phone_number" json:"relationPhoneNumber"`
	CustomerLoan        []CustomerLoan `gorm:"foreignKey:CustomerId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
