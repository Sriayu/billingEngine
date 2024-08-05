package main

import (
	"BillingEngine/src/handler"
	"BillingEngine/src/handler/response"
	"BillingEngine/src/infra"
	"BillingEngine/src/infra/persistence/postgres"
	"BillingEngine/src/usecases"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := infra.Make()
	postgresDb := postgres.New(config.SqlDb)
	custLoanRepo := postgres.NewCustomerLoanPersistence(postgresDb.DB)
	loanPaymentRepo := postgres.NewLoanPaymentPersistence(postgresDb.DB)

	allUsecases := usecases.AllUseCases{
		LoanPayment:  usecases.NewLoanPaymentUsecase(loanPaymentRepo, custLoanRepo),
		CustomerLoan: usecases.NewCustomerLoanUsecase(custLoanRepo, loanPaymentRepo),
	}
	r := makeRoute(allUsecases)
	log.Println("Server Running : http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func makeRoute(usecases usecases.AllUseCases) *chi.Mux {
	r := chi.NewRouter()
	respClient := response.NewResponseClient()
	loanPaymentHandler := handler.NewLoanPaymentHandler(usecases.LoanPayment, respClient)
	customerLoanHandler := handler.NewCustomerLoanHandler(usecases.CustomerLoan, respClient)

	r.Route("/customer-loan", func(r chi.Router) {
		r.Get("/outstanding", customerLoanHandler.GetCustomerOutstanding)
		r.Get("/delinquent", customerLoanHandler.GetCustomerDelinquent)
	})

	r.Route("/loan-payment", func(r chi.Router) {
		r.Post("/make-payment", loanPaymentHandler.PayLoanPayment)
	})

	return r
}
