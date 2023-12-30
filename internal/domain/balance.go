package domain

import "errors"

var (
	ErrIncorrectWalletNumber = errors.New("incorrect wallet number")
	ErrIncorrectWallet       = errors.New("there is no such wallet")
	ErrIncorrectAmount       = errors.New("enter an amount greater than zero")
	ErrNoMoney               = errors.New("not enough money")
)

type InvoiceStatus string

const (
	// заказ загружен в систему, но не попал в обработку.
	Created InvoiceStatus = "CREATED"
	// зачисление выполнено.
	Success InvoiceStatus = "SUCCESS"
	// возникла ошибка при обработки операции.
	Error InvoiceStatus = "ERROR"
)

type Invoice struct {
	Currency     string        `json:"currency"`
	Amount       float32       `json:"amount"`
	UploadedAt   string        `json:"-"`
	WalletNumber string        `json:"wallet_number"`
	Status       InvoiceStatus `json:"-"`
}

type Withdraw struct {
	Currency     string  `json:"currency"`
	Amount       float32 `json:"amount"`
	UploadedAt   string  `json:"-"`
	WalletNumber string  `json:"wallet_number"`
}
