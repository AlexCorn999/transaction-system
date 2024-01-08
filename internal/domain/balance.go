package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrIncorrectWalletNumber = errors.New("incorrect wallet number")
	ErrIncorrectWallet       = errors.New("there is no such wallet")
	ErrIncorrectAmount       = errors.New("enter an amount greater than zero")
	ErrNoMoney               = errors.New("not enough money")
)

type InvoiceStatus string

const (
	// the transaction has been loaded into the system, but has not been processed.
	Created InvoiceStatus = "CREATED"
	// crediting has been completed.
	Success InvoiceStatus = "SUCCESS"
	// an error occurred while processing an operation.
	Error InvoiceStatus = "ERROR"
)

type Invoice struct {
	Currency     string        `json:"currency"`
	Amount       float64       `json:"amount"`
	UploadedAt   string        `json:"-"`
	WalletNumber string        `json:"wallet_number"`
	Status       InvoiceStatus `json:"-"`
	UserID       int64         `json:"-"`
}

type InvoiceDB struct {
	Currency     string
	Amount       decimal.Decimal
	UploadedAt   string
	WalletNumber string
	Status       InvoiceStatus
	UserID       int64
}

type Withdraw struct {
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount"`
	UploadedAt   string  `json:"-"`
	WalletNumber string  `json:"wallet_number"`
	UserID       int64   `json:"-"`
}

type WithdrawDB struct {
	Currency     string
	Amount       decimal.Decimal
	UploadedAt   string
	WalletNumber string
	UserID       int64
}

type BalanceOutput struct {
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}
