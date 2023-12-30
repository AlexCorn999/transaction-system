package service

import (
	"context"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/shopspring/decimal"
)

// Balance выводит сумму баллов лояльности и использованных за весь период регистрации баллов пользователя.
func (w *Withdraw) Balance(ctx context.Context) (*domain.BalanceResult, error) {
	// userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	// if !ok {
	// 	return nil, errors.New("incorrect user id")
	// }

	// узнаем баланс бонусов пользователя
	balanceUser, err := w.repo.Balance(ctx, w)
	if err != nil {
		return nil, err
	}

	// узнаем баланс списанных бонусов пользователя
	balanceWithdraws, err := b.repo.WithdrawBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// чтобы узнать баланс пользователя вычитаем кол-во использованных бонусов
	newBalanceUser := decimal.NewFromFloat32(balanceUser).Sub(decimal.NewFromFloat32(balanceWithdraws))

	var balance domain.BalanceOutput
	balance.Bonuses = float32(newBalanceUser.InexactFloat64())
	balance.Withdraw = balanceWithdraws
	return &balance, nil
}
