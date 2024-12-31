package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	errC := make(chan error)
	resultC := make(chan TransferTxResult)
	for i := 0; i < 2; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: 19,
				ToAccountID:   19,
				Amount:        10,
			})
			errC <- err
			resultC <- result
		}()
	}
	for i := 0; i < 2; i++ {
		err := <-errC
		require.NoError(t, err)
		result := <-resultC
		require.NotEmpty(t, result)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, int64(19), transfer.FromAccountID)
		require.Equal(t, int64(19), transfer.ToAccountID)
		require.Equal(t, int64(10), transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// _, err = store.GetTransfer(context.Background())
		// require.NoError(t, err)

		// // check entries
		// fromEntry := result.FromEntry
		// require.NotEmpty(t, fromEntry)
		// require.Equal(t, 19, fromEntry.AccountID)
		// require.Equal(t, -10, fromEntry.Amount)
		// require.NotZero(t, fromEntry.ID)
		// require.NotZero(t, fromEntry.CreatedAt)

		// _, err = store.GetEntry(context.Background())
		// require.NoError(t, err)

		// toEntry := result.ToEntry
		// require.NotEmpty(t, toEntry)
		// require.Equal(t, int64(19), toEntry.AccountID)
		// require.Equal(t, int64(10), toEntry.Amount)
		// require.NotZero(t, toEntry.ID)
		// require.NotZero(t, toEntry.CreatedAt)

		// _, err = store.GetEntry(context.Background())
		// require.NoError(t, err)

		// // check accounts
		// fromAccount := result.FromAccount
		// require.NotEmpty(t, fromAccount)
		// require.Equal(t, 19, fromAccount.ID)

		// toAccount := result.ToAccount
		// require.NotEmpty(t, toAccount)
		// require.Equal(t, 19, toAccount.ID)
	}
}
