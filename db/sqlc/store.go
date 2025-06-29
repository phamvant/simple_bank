package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
	Currency      string
}

type TransferTxResult struct {
	Transfer    Transfer
	FromEntry   Entry
	ToEntry     Entry
	FromAccount Account
	ToAccount   Account
}

var txKey = struct{}{}

func (store *Store) TransferTX(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})

		if err != nil {
			return err
		}

		if params.FromAccountID < params.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, params.FromAccountID, -params.Amount, params.ToAccountID, params.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, params.ToAccountID, params.Amount, params.FromAccountID, -params.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, account1ID int64, amount1 int64, account2ID int64, amount2 int64) (account1 Account, account2 Account, err error) {

	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})

	return

}
