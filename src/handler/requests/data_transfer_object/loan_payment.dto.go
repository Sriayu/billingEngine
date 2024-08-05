package datatransferobject

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ILoanPaymentRequest interface {
	Validate() error
}

type (
	PayCustomerLoanRequestDto struct {
		CustomerLoanId int
		TotalPricePaid float64
	}
)

func NewPayCustomerLoanRequestDto(
	customerLoanId int,
	totalPricePaid float64,
) *PayCustomerLoanRequestDto {
	return &PayCustomerLoanRequestDto{
		CustomerLoanId: customerLoanId,
		TotalPricePaid: totalPricePaid,
	}
}

func (dto *PayCustomerLoanRequestDto) Validate() (err error) {
	if err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.CustomerLoanId, validation.Required),
		validation.Field(&dto.TotalPricePaid, validation.Required),
	); err != nil {
		retErr := fmt.Errorf("Invalid request pay customer loan")
		return retErr
	}
	return nil
}
