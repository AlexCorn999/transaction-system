package domain

import "errors"

var (
	ErrIncorrectCurrency = errors.New("incorrect currency")
)

var Currency = map[string]string{
	"USD": "USD",
	"EUR": "EUR",
	"GBP": "GBP",
	"JPY": "JPY",
	"AUD": "AUD",
	"CAD": "CAD",
	"CHF": "CHF",
	"CNY": "CNY",
	"SEK": "SEK",
	"NZD": "NZD",
	"DKK": "DKK",
	"NOK": "NOK",
	"SGD": "SGD",
	"CZK": "CZK",
	"HKD": "HKD",
	"MXN": "MXN",
	"PLN": "PLN",
	"RUB": "RUB",
	"TRY": "TRY",
	"ZAR": "ZAR",
	"CNH": "CNH",
}
