package dto

type MidtransSnapRequest struct {
	TransactionDetails struct {
		OrderID  int     `json:"order_id"`
		GrossAmt float64 `json:"gross_amount"`
	} `json:"transaction_details"`
	PaymentType string `json:"payment_type"`
	Customer    string `json:"customer"`
	// Items       []Item `json:"item_details"`
}

type MidtransSnapResponse struct {
	Token        string   `json:"token" binding:"omitempty"`
	RedirectUrl  string   `json:"redirect_url"`
	ErrorMessage []string `json:"error_message" binding:"omitempty"`
}
