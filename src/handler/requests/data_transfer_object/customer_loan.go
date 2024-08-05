package datatransferobject

import (
	"BillingEngine/src/infra/constants"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ICustomerLoanRequest interface {
	Validate() error
}

type (
	CustomerLoanListRequestDto struct {
		CustomerId int
		LoanId     int
		Status     string
	}
	UpdateCustomerLoanRequestDto struct {
		CustomerLoanId  int
		TotalLoanPaidId float64
	}
)

func NewCustomerLoanDelinquentRequestDto(
	customerId int,
	loanId int,
	status string,
) *CustomerLoanListRequestDto {
	return &CustomerLoanListRequestDto{
		CustomerId: customerId,
		LoanId:     loanId,
		Status:     status,
	}
}

func (dto *CustomerLoanListRequestDto) Validate() (err error) {
	if err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.CustomerId, validation.Min(1)),
		validation.Field(&dto.LoanId, validation.Min(1)),
		validation.Field(&dto.Status, validation.Required, validation.In(
			constants.REQ_STATUS_CUSTOMER_LOAN_DELINQUENT,
			constants.REQ_STATUS_CUSTOMER_LOAN_OUTSTANDING,
		)),
	); err != nil {
		retErr := fmt.Errorf("Invalid request customer loan list")
		return retErr
	}
	return nil
}

func NewCustomerLoanRequestDto(
	customerLoanId int,
	totalLoanPaidId float64,
) *UpdateCustomerLoanRequestDto {
	return &UpdateCustomerLoanRequestDto{
		CustomerLoanId:  customerLoanId,
		TotalLoanPaidId: totalLoanPaidId,
	}
}

func (dto *UpdateCustomerLoanRequestDto) Validate() (err error) {
	if err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.CustomerLoanId, validation.Required, validation.Min(1)),
		validation.Field(&dto.TotalLoanPaidId, validation.Required),
	); err != nil {
		retErr := fmt.Errorf("Invalid body request update status update customer loan")
		return retErr
	}
	return nil
}
