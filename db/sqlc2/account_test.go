package db

import (
	"context"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	args := CreateAccountParams{
		Owner:    "asd",
		Balance:  123,
		Currency: "USD",
	}
	testQueries.CreateAccount(context.Background(), args)
	act, _ := testQueries.GetAccount(context.Background())
	require.Equal(t, args.Balance, act.Balance)
	require.Equal(t, args.Owner, act.Owner)
	require.Equal(t, args.Currency, act.Currency)
	require.NotZero(t, act.ID)
}

func TestGetListAccount(t *testing.T) {
	// for i := 0; i < 5; i++ {
	// 	args := CreateAccountParams{
	// 		Owner:    "asd",
	// 		Balance:  123,
	// 		Currency: "USD",
	// 	}
	// 	testQueries.CreateAccount(context.Background(), args)
	// }
	arg := GetListAccountParams{
		Limit:  1,
		Offset: 1,
	}
	acts, _ := testQueries.GetListAccount(context.Background(), arg)
	require.NotEmpty(t, acts)
	for _, v := range acts {
		require.NotEmpty(t, v)
	}
	// for i := 0; i < 5; i++ {
	// 	args := CreateAccountParams{}
	// 	testQueries.CreateAccount(context.Background(), args)
	// }
	arg.Limit = -1
	arg.Offset = -1
	_, err := testQueries.GetListAccount(context.Background(), arg)
	require.Error(t, err)
	// err = testQueries.DeleteAccountById(context.Background())
	// require.NoError(t, err)
}
