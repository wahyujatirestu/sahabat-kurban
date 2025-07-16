package model

import (
	"time"

	"github.com/google/uuid"
)

type PembayaranKurban struct {
	ID                	uuid.UUID	`db:"id"`
	OrderID           	string		`db:"order_id"`
	TransactionID     	string		`db:"transaction_id"`
	PekurbanID        	uuid.UUID	`db:"pekurban_id"`
	Metode            	string		`db:"metode"`
	PaymentType       	*string		`db:"payment_type"`
	VANumber          	*string		`db:"va_number"`
	Jumlah            	float64		`db:"jumlah"`
	Status            	string		`db:"status"`
	FraudStatus       	*string		`db:"fraud_status"`
	ApprovalCode      	*string		`db:"approval_code"`
	TransactionTime   	*time.Time	`db:"transaction_time"`
	TanggalPembayaran 	time.Time	`db:"tanggal_pembayaran"`
	Created_At         	time.Time	`db:"created_at"`
	Updated_At         	time.Time	`db:"updated_at"`
}

