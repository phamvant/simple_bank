package db

import (
	"TestProj/utils"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomName(),
		Balance:  1000,
		Currency: utils.RandomCurrency(),
	}

	newAccount, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, arg.Owner, newAccount.Owner)
	require.Equal(t, arg.Balance, newAccount.Balance)
	require.Equal(t, arg.Currency, newAccount.Currency)

	return newAccount
}
