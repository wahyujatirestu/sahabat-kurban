package service

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/wahyujatirestu/sahabat-kurban/payments/model"
)

type MidtransService interface {
	Charge(req *model.MidtransChargeRequest) (*model.MidtransChargeResponse, error)
}

type midtransService struct {
	client		*resty.Client
	serverKey	string
	baseURL		string
}

func NewMidtransService() MidtransService {
	key := os.Getenv("MIDTRANS_SERVER_KEY")
	if key == "" {
		panic("MIDTRANS_SERVER_KEY is required")
	}
	return &midtransService{
		client: resty.New(),
		serverKey: key,
		baseURL: "https://api.sandbox.midtrans.com",
	}
}

func (m *midtransService) Charge(req *model.MidtransChargeRequest) (*model.MidtransChargeResponse, error) {
	endpoint := m.baseURL + "/v2/charge"

	auth := base64.StdEncoding.EncodeToString([]byte(m.serverKey + ":"))

	var response model.MidtransChargeResponse
	res, err := m.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Basic "+auth).
		SetBody(req).SetResult(&response).Post(endpoint)

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("Midtrans error: %s", response.StatusMessage)
	}

	for _, action := range response.Actions {
		if action.Name == "generate-qr-code" {
			response.QRUrl = &action.URL
			break
		}
	}

	return &response, nil
}