Berikut Entity Billing Engine : 
<img width="890" alt="image" src="https://github.com/user-attachments/assets/7c503944-e273-4774-bd47-6c600c3bbc7a">

Steps to running service : 
1. Provide postgreSQl database local, that must be set in file .env
2. run commmand: **go installs** to get all package that used in service
2. create database configuration based on file .env
3. run commmand: **go run src/infra/cmd/migrate.go** to automatic create tables
4. run commmand: **go run main.go** to running services


API ENDPOINT : 
1. GetOutstanding : {{base_url}}/customer-loan/outstanding?status=Outstanding (GET)
2. MakePayment : {{base_url}}/loan-payment/make-payment (POST)
                Body req : 
                {
                    "customerLoanId" : 2, 
                    "totalPricePaid" : 2200000
                }
3. IsDelinquent : {{base_url}}/customer-loan/outstanding?status=Delinquent

