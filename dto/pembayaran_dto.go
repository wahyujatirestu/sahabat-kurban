package dto

import (

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	payment "github.com/wahyujatirestu/sahabat-kurban/payments/model"
)

type CreatePaymentRequest struct {
	PekurbanID    uuid.UUID `json:"pekurban_id" binding:"required"`
	Metode        string    `json:"metode" binding:"required"`         
	Bank          string    `json:"bank,omitempty"`
}


type PaymentResponse struct {
	ID              string   `json:"id"`
	OrderID         string   `json:"order_id"`
	TransactionID   string   `json:"transaction_id"`
	PekurbanID      string   `json:"pekurban_id"`
	Metode          string   `json:"metode"`
	PaymentType     *string  `json:"payment_type,omitempty"`
	VANumber        *string  `json:"va_number,omitempty"`
	Status          string   `json:"status"`
	FraudStatus     *string  `json:"fraud_status,omitempty"`
	ApprovalCode    *string  `json:"approval_code,omitempty"`
	TransactionTime *string  `json:"transaction_time,omitempty"`
	RedirectURL     *string  `json:"redirect_url,omitempty"`
	Jumlah          float64  `json:"jumlah"`
}

func ToPaymentResponse(p *model.PembayaranKurban, jumlah float64, mid *payment.MidtransChargeResponse) PaymentResponse {
	var trxTime *string
	if p.TransactionTime != nil {
		str := p.TransactionTime.Format("2006-01-02 15:04:05")
		trxTime = &str
	}

	var redirectURL *string
	if mid != nil && mid.QRUrl != nil {
		redirectURL = mid.QRUrl
	}

	return PaymentResponse{
		ID:              p.ID.String(),
		OrderID:         p.OrderID,
		TransactionID:   p.TransactionID,
		PekurbanID:      p.PekurbanID.String(),
		Metode:          p.Metode,
		PaymentType:     p.PaymentType,
		VANumber:        p.VANumber,
		Status:          p.Status,
		FraudStatus:     p.FraudStatus,
		ApprovalCode:    p.ApprovalCode,
		TransactionTime: trxTime,
		RedirectURL:     redirectURL,
		Jumlah:          jumlah,
	}
}

func ToMidtransChargeRequest(orderID string, grossAmount float64, name, email, phone string, req CreatePaymentRequest) *payment.MidtransChargeRequest {
	payload := &payment.MidtransChargeRequest{
		PaymentType: req.Metode,
		TransactionDetails: payment.TransactionDetails{
			OrderID:     orderID,
			GrossAmount: grossAmount,
		},
		CustomerDetails: payment.CustomerDetails{
			FirstName: name,
			Email:     email,
			Phone:     phone,
		},
	}

	if req.Metode == "bank_transfer" {
		payload.BankTransfer = &payment.BankTransfer{Bank: req.Bank}
	}

	if req.Metode == "qris" {
		payload.QR = &payment.QRIS{}
	}

	return payload
}
