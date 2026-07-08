package panorama

import "errors"

// --------------------------------- GET BALANCE REQUEST ---------------------------------

// GetBalanceRequest is the request structure for GetBalance
type GetBalanceRequest struct {
	Signature string `json:"signature"`
}

func (r *GetBalanceRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

// --------------------------------- GET BALANCE RESPONSE ---------------------------------

// GetBalanceResponse is the response structure for GetBalance
type GetBalanceResponse struct {
	Result struct {
		BasicPanoramaResult
		Response BalanceResponse `json:"response"`
	} `json:"result"`
}

// BalanceResponse is the response payload for GetBalance
type BalanceResponse struct {
	Balance                     float64 `json:"balance"`
	CustomerParentBalance       int     `json:"customer_parent_balance"`
	CreditLimit                 int     `json:"credit_limit"`
	CustomerParentCurrencyCode  string  `json:"customer_parent_currency_code"`
	CurrencyCode                string  `json:"currency_code"`
	IsShowBalance               bool    `json:"is_show_balance"`
	IsShowCustomerParentBalance bool    `json:"is_show_customer_parent_balance"`
	IsShowCreditLimit           bool    `json:"is_show_credit_limit"`
	PointReward                 float64 `json:"point_reward"`
	IsShowPointReward           bool    `json:"is_show_point_reward"`
}
