package models

import (
	"time"
)

// Currency codes (ISO 4217) - ECB supported currencies
const (
	CurrencyEUR = "EUR"
	CurrencyUSD = "USD"
	CurrencyJPY = "JPY"
	CurrencyCZK = "CZK"
	CurrencyDKK = "DKK"
	CurrencyGBP = "GBP"
	CurrencyHUF = "HUF"
	CurrencyPLN = "PLN"
	CurrencyRON = "RON"
	CurrencySEK = "SEK"
	CurrencyCHF = "CHF"
	CurrencyISK = "ISK"
	CurrencyNOK = "NOK"
	CurrencyTRY = "TRY"
	CurrencyAUD = "AUD"
	CurrencyBRL = "BRL"
	CurrencyCAD = "CAD"
	CurrencyCNY = "CNY"
	CurrencyHKD = "HKD"
	CurrencyIDR = "IDR"
	CurrencyILS = "ILS"
	CurrencyINR = "INR"
	CurrencyKRW = "KRW"
	CurrencyMXN = "MXN"
	CurrencyMYR = "MYR"
	CurrencyNZD = "NZD"
	CurrencyPHP = "PHP"
	CurrencySGD = "SGD"
	CurrencyTHB = "THB"
	CurrencyZAR = "ZAR"
)

type RateCurrency struct {
	ID           string    `json:"id"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	Rate         float64   `json:"rate"`
	RateDate     time.Time `json:"rate_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpsertRateCurrencyRequest struct {
	FromCurrency string    `json:"from_currency" binding:"required"`
	ToCurrency   string    `json:"to_currency" binding:"required"`
	Rate         float64   `json:"rate" binding:"required"`
	RateDate     time.Time `json:"rate_date" binding:"required"`
}

type RateCurrencySearchParams struct {
	FromCurrency *string
	ToCurrency   *string
	RateDate     *time.Time
	PageSize     int
	PageNumber   int
}

func (p *RateCurrencySearchParams) Offset() int {
	if p.PageSize <= 0 {
		return 0
	}
	return (p.PageNumber - 1) * p.PageSize
}
