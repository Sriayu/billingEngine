package requests

import (
	datatransferobject "BillingEngine/src/handler/requests/data_transfer_object"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	PayCustomerLoanRequest struct {
		CustomerLoanId int
		TotalPricePaid float64
	}
)

func (req *PayCustomerLoanRequest) Validate(r *http.Request) (dto datatransferobject.ILoanPaymentRequest, err error) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return dto, fmt.Errorf("Invalid payload format json")
	}

	dto = datatransferobject.NewPayCustomerLoanRequestDto(req.CustomerLoanId, req.TotalPricePaid)
	err = dto.Validate()
	if err != nil {
		return dto, err
	}

	return dto, nil
}
