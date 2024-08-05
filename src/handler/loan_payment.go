package handler

import (
	"BillingEngine/src/handler/requests"
	"BillingEngine/src/handler/response"
	"BillingEngine/src/usecases"
	"net/http"
)

type ILoanPaymentHandler interface {
	PayLoanPayment(w http.ResponseWriter, r *http.Request)
}

type loanPaymentHandler struct {
	usecase             usecases.ILoanPaymentUsecase
	customerLoanUsecase usecases.ICustomerLoanUsecase
	response            response.IResponseClient
}

// NewLoanPaymentHandler ...
func NewLoanPaymentHandler(u usecases.ILoanPaymentUsecase, r response.IResponseClient) ILoanPaymentHandler {
	return &loanPaymentHandler{
		usecase:  u,
		response: r,
	}
}

func (u *loanPaymentHandler) PayLoanPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := requests.PayCustomerLoanRequest{}
	dtoReq, err := req.Validate(r)
	if err != nil {
		u.response.HttpError(w, err, http.StatusBadRequest)
		return
	}
	err = u.usecase.PayLoanPayment(ctx, dtoReq)
	if err != nil {
		u.response.HttpError(w, err, http.StatusBadRequest)
		return
	}

	u.response.ResponseJSON(w, "Customer making payment", nil, nil)
}
