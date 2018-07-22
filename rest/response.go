package rest

type ItemResponse struct {
	Result *ItemInfo `json:"result,omitempty"`
	Errors []ErrorRs `json:"errors,omitempty"`
}

type ItemInfo struct {
	Name  string `json:"name,omitempty"`
	Price *Money `json:"money,omitempty"`
}

type Money struct {
	Currency string  `json:"currency,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
}

type ErrorRs struct {
	Message string `json:"message,omitempty"`
}
