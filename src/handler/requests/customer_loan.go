package requests

import (
	datatransferobject "BillingEngine/src/handler/requests/data_transfer_object"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type (
	GetCustomerLoanListRequest struct {
		CustomerId *int64
		LoanId     *int64
		Status     string
	}

	UpdateCustomerLoanRequest struct {
		CustomerLoanId int
		TotalPaidLoan  float64
	}
)

func (req *GetCustomerLoanListRequest) Validate(r *http.Request) (dto datatransferobject.ICustomerLoanRequest, err error) {
	customerId := r.URL.Query().Get("customerId")
	custId, err := strconv.ParseInt(customerId, 10, 32)
	if err != nil {
		custId = 0
	}
	loanId := r.URL.Query().Get("loanId")
	lId, err := strconv.ParseInt(loanId, 10, 32)
	if err != nil {
		lId = 0
	}

	status := r.URL.Query().Get("status")

	req.CustomerId = &custId
	req.LoanId = &lId
	req.Status = status
	dto = datatransferobject.NewCustomerLoanDelinquentRequestDto(int(custId), int(lId), req.Status)
	err = dto.Validate()
	if err != nil {
		return dto, err
	}

	return dto, nil
}

func (req *UpdateCustomerLoanRequest) Validate(r *http.Request) (dto datatransferobject.ICustomerLoanRequest, err error) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		return dto, fmt.Errorf("Invalid id param url")
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return dto, fmt.Errorf("Invalid payload format json")
	}

	req.CustomerLoanId = int(id)

	dto = datatransferobject.NewCustomerLoanRequestDto(
		req.CustomerLoanId,
		req.TotalPaidLoan,
	)
	err = dto.Validate()
	if err != nil {
		return dto, err
	}

	return dto, nil
}
