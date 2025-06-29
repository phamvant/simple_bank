package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	store := NewStore(conn)

	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)

	n := 3
	amount := int64(10)

	results := make(chan TransferTxResult)
	errs := make(chan error)

	// existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i)
		ctx := context.WithValue(context.Background(), txKey, txName)

		go func() {
			result, err := store.TransferTX(ctx, TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		fromAccountResult := result.FromAccount
		require.NotEmpty(t, fromAccountResult)
		require.Equal(t, fromAccountResult.ID, fromAccount.ID)

		toAccountResult := result.ToAccount
		require.NotEmpty(t, toAccountResult)
		require.Equal(t, toAccountResult.ID, toAccount.ID)

		fmt.Println(">> tx:", fromAccountResult.Balance, toAccountResult.Balance)
		diff1 := fromAccount.Balance - fromAccountResult.Balance
		diff2 := toAccountResult.Balance - toAccount.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
	}

	updatedAccount1, err := store.GetAccountById(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccountById(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-updatedAccount1.Balance, amount*int64(n))
	require.Equal(t, updatedAccount2.Balance-toAccount.Balance, amount*int64(n))
}

func TestCreateTransfer2Way(t *testing.T) {
	store := NewStore(conn)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	n := 10
	amount := int64(10)

	results := make(chan TransferTxResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {

		var fromAccountID int64
		var toAccountID int64

		if i%2 == 0 {
			fromAccountID = account1.ID
			toAccountID = account2.ID
		} else {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			result, err := store.TransferTX(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := store.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
