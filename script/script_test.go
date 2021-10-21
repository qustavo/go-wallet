package script_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/qustavo/go-wallet/script"
)

func TestInvalidKeys(t *testing.T) {
	expr := Sh(Pkh("the_key"))
	_, err := expr.Eval()
	require.Error(t, err)
}
