package db

import (
	"context"
	"testing"

	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T){
	arg := CreatedAccountParams{
		Owner: "tom",
		Balance: 100,
		Currency: "USD",
	}

	account, err := testQueries.CreatedAccount(context.Background(), arg)

	require.NoError(t,err)
	require.NotEmpty(t,account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	
}