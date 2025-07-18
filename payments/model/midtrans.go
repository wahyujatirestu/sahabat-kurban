package model

type MidtransChargeRequest struct {
	PaymentType 		string               `json:"payment_type"`
	TransactionDetails 	TransactionDetails 	 `json:"transaction_details"`
	CustomerDetails    	CustomerDetails    	 `json:"customer_details"`
	BankTransfer       	*BankTransfer      	 `json:"bank_transfer,omitempty"`
	QR                 	*QRIS              	 `json:"qris,omitempty"`
}

type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type CustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type BankTransfer struct {
	Bank string `json:"bank"`
}

type QRIS struct{} 

type Action struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MidtransChargeResponse struct {
	StatusCode        string  	 `json:"status_code"`
	StatusMessage     string  	 `json:"status_message"`
	TransactionID     string  	 `json:"transaction_id"`
	OrderID           string   	 `json:"order_id"`
	GrossAmount       string  	 `json:"gross_amount"`
	PaymentType       string   	 `json:"payment_type"`
	TransactionTime   string  	 `json:"transaction_time"`
	TransactionStatus string  	 `json:"transaction_status"`
	FraudStatus       *string    `json:"fraud_status,omitempty"`
	ApprovalCode      *string 	 `json:"approval_code,omitempty"`
	VANumbers 		  []VANumber `json:"va_numbers,omitempty"`
	Actions           []Action    `json:"actions"`
	QRUrl 			  *string 	 `json:"qr_code_url,omitempty"`
}

type VANumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

