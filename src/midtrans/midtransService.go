package midtrans

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fp_pinjaman_online/model/dto"
	"log"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransService interface {
	Pay(payload dto.MidtransSnapRequest) (dto.MidtransSnapResponse, error)
	VerifyPayment(orderId int) (bool, error)
}

type midtransService struct {
	client *resty.Client
	url    string
}

func (m *midtransService) Pay(payload dto.MidtransSnapRequest) (dto.MidtransSnapResponse, error) {

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	encodedKey := base64.StdEncoding.EncodeToString([]byte(serverKey))

	resp, err := m.client.R().
		SetHeader("Authorization", "Basic "+encodedKey).
		SetBody(payload).Post(m.url)

	if err != nil {
		log.Println("Error bayar: ", err.Error())
		return dto.MidtransSnapResponse{}, err
	}

	var snapResponse dto.MidtransSnapResponse

	if err := json.Unmarshal(resp.Body(), &snapResponse); err != nil {
		log.Println("Error: ", err.Error())
		return dto.MidtransSnapResponse{}, err
	}

	redirectUrl := fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/%s", snapResponse.Token)
	snapResponse.RedirectUrl = redirectUrl
	return snapResponse, nil
}

func NewMidtransService(client *resty.Client) MidtransService {
	return &midtransService{
		client: client,
		url:    "https://app.sandbox.midtrans.com/snap/v1/transactions",
	}
}

func (m midtransService) VerifyPayment(orderId int) (bool, error) {
	var client coreapi.Client
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	client.New(serverKey, midtrans.Sandbox)

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(strconv.Itoa(orderId))
	if e != nil {
		return false, e
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				return true, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return false, nil
}
