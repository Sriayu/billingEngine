package handler

import (
	"BillingEngine/src/handler/requests"
	"BillingEngine/src/handler/response"
	"BillingEngine/src/usecases"
	"net/http"
)

type ICustomerLoanHandler interface {
	GetCustomerDelinquent(w http.ResponseWriter, r *http.Request)
	GetCustomerOutstanding(w http.ResponseWriter, r *http.Request)
	CreateCustomerLoan(w http.ResponseWriter, r *http.Request)
}

type customerLoanHandler struct {
	usecase  usecases.ICustomerLoanUsecase
	response response.IResponseClient
}

// NewCustomerLoanHandler ...
func NewCustomerLoanHandler(u usecases.ICustomerLoanUsecase, r response.IResponseClient) ICustomerLoanHandler {
	return &customerLoanHandler{
		usecase:  u,
		response: r,
	}
}

func (d *customerLoanHandler) GetCustomerOutstanding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := requests.GetCustomerLoanListRequest{}
	dtoReq, err := req.Validate(r)
	if err != nil {
		d.response.HttpError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := d.usecase.GetCustomerLoanByStatus(ctx, dtoReq)
	if err != nil {
		d.response.HttpError(w, err, http.StatusBadRequest)
		return
	}

	d.response.ResponseJSON(w, "Success get customer outstanding list", resp, nil)
}

func (d *customerLoanHandler) GetCustomerDelinquent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := requests.GetCustomerLoanListRequest{}
	dtoReq, err := req.Validate(r)
	if err != nil {
		d.response.HttpError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := d.usecase.GetCustomerLoanByStatus(ctx, dtoReq)
	if err != nil {
		d.response.HttpError(w, err, http.StatusBadRequest)
		return
	}

	d.response.ResponseJSON(w, "Success get customer delinquent list", resp, nil)

}

func (d *customerLoanHandler) CreateCustomerLoan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := requests.UpdateCustomerLoanRequest{}
	dtoReq, err := req.Validate(r)
	if err != nil {
		d.response.HttpError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := d.usecase.UpdateCustomerLoan(ctx, dtoReq)
	if err != nil {
		d.response.HttpError(w, err, http.StatusBadRequest)
		return
	}

	d.response.ResponseJSON(w, "Success update status dating", resp, nil)
}
