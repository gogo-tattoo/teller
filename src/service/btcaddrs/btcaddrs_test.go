package btcaddrs

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/skycoin/teller/src/logger"
	"github.com/stretchr/testify/require"
)

func TestNewBtcAddrs(t *testing.T) {
	db, shutdown := setupDB(t)
	defer shutdown()

	invalidAddr := "invalid address"
	addrJSON := addressJSON{
		BtcAddresses: []string{
			"14JwrdSxYXPxSi6crLKVwR4k2dbjfVZ3xj",
			"1JNonvXRyZvZ4ZJ9PE8voyo67UQN1TpoGy",
			"1JrzSx8a9FVHHCkUFLB2CHULpbz4dTz5Ap",
			invalidAddr,
		},
	}

	v, err := json.Marshal(addrJSON)
	require.Nil(t, err)

	btca, err := New(db, bytes.NewReader(v), logger.NewLogger("", true))
	require.Nil(t, err)

	for _, addr := range addrJSON.BtcAddresses {
		_, ok := btca.addresses[addr]
		if addr == invalidAddr {
			require.False(t, ok)
			continue
		}
		require.True(t, ok)
	}
}

func TestNewAddress(t *testing.T) {
	db, shutdown := setupDB(t)
	defer shutdown()

	addrJSON := addressJSON{
		BtcAddresses: []string{
			"14JwrdSxYXPxSi6crLKVwR4k2dbjfVZ3xj",
			"1JNonvXRyZvZ4ZJ9PE8voyo67UQN1TpoGy",
			"1JrzSx8a9FVHHCkUFLB2CHULpbz4dTz5Ap",
		},
	}

	v, err := json.Marshal(addrJSON)
	require.Nil(t, err)

	btca, err := New(db, bytes.NewReader(v), logger.NewLogger("", true))
	require.Nil(t, err)

	addr, err := btca.NewAddress()
	require.Nil(t, err)

	// check if the addr still in the address pool
	_, ok := btca.addresses[addr]
	require.False(t, ok)

	// check if the addr is in used storage
	require.True(t, btca.used.IsExsit(addr))

	btca1, err := New(db, bytes.NewReader(v), logger.NewLogger("", true))
	require.Nil(t, err)

	_, ok = btca1.addresses[addr]
	require.False(t, ok)

	require.True(t, btca1.used.IsExsit(addr))

	// run out all addresses
	for i := 0; i < 2; i++ {
		btca1.NewAddress()
	}

	_, err = btca1.NewAddress()
	require.Equal(t, ErrDepositAddressEmpty, err)
}
