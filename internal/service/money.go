package service

import (
	"context"
	"strings"
	"time"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/repository"
	"github.com/shopspring/decimal"
)

type MoneyManagementRepository interface {
	Invoice(ctx context.Context, invoice *domain.InvoiceDB) error
	Withdraw(ctx context.Context, withdraw *domain.WithdrawDB) error
	CheckWallet(ctx context.Context, withdraw *domain.WithdrawDB) (int, error)
	InvoiceToUser(ctx context.Context, invoice *domain.InvoiceDB) error
	Balance(ctx context.Context, withdraw *domain.WithdrawDB) (float64, error)
	BalanceActual(userID int64) ([]domain.BalanceOutput, error)
	BalanceFrozen(userID int64) ([]domain.BalanceOutput, error)
	WithdrawBalance(ctx context.Context, withdraw *domain.WithdrawDB) (float64, error)
}

type Money struct {
	repo    MoneyManagementRepository
	storage *repository.Storage
}

func NewInvoices(repo MoneyManagementRepository, storage *repository.Storage) *Money {
	return &Money{
		repo:    repo,
		storage: storage,
	}
}

// Invoice credits money to the user's account.
func (m *Money) Invoice(ctx context.Context, invoice *domain.Invoice) error {
	// if the wallet number is empty
	if len(strings.TrimSpace(invoice.WalletNumber)) == 0 {
		return domain.ErrIncorrectWalletNumber
	}

	// if the replenishment amount is negative or equal to zero
	if invoice.Amount <= 0 {
		return domain.ErrIncorrectAmount
	}

	// if no such currency exists
	if _, ok := domain.Currency[invoice.Currency]; !ok {
		return domain.ErrIncorrectCurrency

	}

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return domain.ErrIncorrectUserID
	}

	// invoice for database with decimal amount
	invoiceForDB := domain.InvoiceDB{
		Currency:     invoice.Currency,
		Amount:       decimal.NewFromFloat(invoice.Amount),
		UploadedAt:   time.Now().Format(time.RFC3339),
		WalletNumber: invoice.WalletNumber,
		Status:       domain.Created,
		UserID:       userID,
	}

	return m.repo.Invoice(ctx, &invoiceForDB)
}

// Withdraw realizes currency debit from the wallet.
func (m *Money) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {

	// if the wallet number is empty
	if len(strings.TrimSpace(withdraw.WalletNumber)) == 0 {
		return domain.ErrIncorrectWalletNumber
	}

	// if the replenishment amount is negative or equal to zero
	if withdraw.Amount <= 0 {
		return domain.ErrIncorrectAmount
	}

	// if no such currency exists
	if _, ok := domain.Currency[withdraw.Currency]; !ok {
		return domain.ErrIncorrectCurrency

	}

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return domain.ErrIncorrectUserID
	}

	// withdraw for database with decimal amount
	withdrawForDB := domain.WithdrawDB{
		Currency:     withdraw.Currency,
		Amount:       decimal.NewFromFloat(withdraw.Amount),
		UploadedAt:   time.Now().Format(time.RFC3339),
		WalletNumber: withdraw.WalletNumber,
		UserID:       userID,
	}

	// transaction initiation
	tx, err := m.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// bonus debit
	err = m.repo.Withdraw(repository.InjectTx(ctx, tx), &withdrawForDB)
	if err != nil {
		return err
	}

	invoiceUserID, err := m.repo.CheckWallet(repository.InjectTx(ctx, tx), &withdrawForDB)
	if err != nil {
		return err
	}

	// invoice for database with decimal amount
	invoice := domain.InvoiceDB{
		Currency:     withdraw.Currency,
		Amount:       decimal.NewFromFloat(withdraw.Amount),
		UploadedAt:   time.Now().Format(time.RFC3339),
		WalletNumber: withdraw.WalletNumber,
		Status:       domain.Created,
		UserID:       int64(invoiceUserID),
	}

	// transferring money to another user
	err = m.repo.InvoiceToUser(repository.InjectTx(ctx, tx), &invoice)
	if err != nil {
		return err
	}

	// find out the wallet balance where the status is success
	balance, err := m.repo.Balance(repository.InjectTx(ctx, tx), &withdrawForDB)
	if err != nil {
		return err
	}

	// find out the balance of amounts written off
	balanceWithdraws, err := m.repo.WithdrawBalance(repository.InjectTx(ctx, tx), &withdrawForDB)
	if err != nil {
		return err
	}

	// verification for bonus debit execution
	sum := decimal.NewFromFloat(balance).Sub(decimal.NewFromFloat(balanceWithdraws))
	if sum.LessThan(decimal.Zero) {
		// if the balance is negative
		tx.Rollback()
		return domain.ErrNoMoney
	}

	return tx.Commit()
}

// Balance returns the user's wallet balance with success status.
func (m *Money) Balance(ctx context.Context) ([]domain.BalanceOutput, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, domain.ErrIncorrectUserID
	}
	return m.repo.BalanceActual(userID)
}

// BalanceFrozen displays the user's balance in the created status.
func (m *Money) BalanceFrozen(ctx context.Context) ([]domain.BalanceOutput, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, domain.ErrIncorrectUserID
	}
	return m.repo.BalanceFrozen(userID)
}
