package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/repository"
	"github.com/shopspring/decimal"
)

type WithdrawRepository interface {
	Withdraw(ctx context.Context, withdraw *domain.Withdraw) error
	Balance(ctx context.Context, withdraw *domain.Withdraw) (float32, error)
	WithdrawBalance(ctx context.Context, withdraw *domain.Withdraw) (float32, error)
}

type Withdraw struct {
	repo    WithdrawRepository
	storage *repository.Storage
}

func NewWithdraw(repo WithdrawRepository, storage *repository.Storage) *Withdraw {
	return &Withdraw{
		repo:    repo,
		storage: storage,
	}
}

// Withdraw реализует списание валюты с кошелька.
func (w *Withdraw) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	if len(strings.TrimSpace(withdraw.WalletNumber)) == 0 {
		return domain.ErrIncorrectWalletNumber
	}

	if withdraw.Amount <= 0 {
		return domain.ErrIncorrectAmount
	}

	if _, ok := domain.Currency[withdraw.Currency]; !ok {
		return domain.ErrIncorrectCurrency

	}

	withdraw.UploadedAt = time.Now().Format(time.RFC3339)

	// Начало транзакции
	tx, err := w.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// сразу списываем бонусы
	err = w.repo.Withdraw(ctx, &withdraw)
	if err != nil {
		return err
	}

	// узнаем баланс валютного кошелька
	balance, err := w.repo.Balance(ctx, &withdraw)
	if err != nil {
		return err
	}

	// ffffff
	fmt.Println(balance)

	// узнаем баланс списанных сумм
	balanceWithdraws, err := w.repo.WithdrawBalance(ctx, &withdraw)
	if err != nil {
		return err
	}

	// ffffff
	fmt.Println(balanceWithdraws)

	// проверка для проведения списания бонусов
	sum := decimal.NewFromFloat32(balance).Sub(decimal.NewFromFloat32(balanceWithdraws))
	if sum.LessThan(decimal.Zero) {
		// если баланс в минусе
		tx.Rollback()

		// ffffff
		fmt.Println("ОТКАТТТТ")
		return domain.ErrNoMoney
	}

	return tx.Commit()
}
