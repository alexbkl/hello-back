package web3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateMessageSignature(t *testing.T) {
	testCases := []struct {
		name          string
		signature     string
		checkResponse func(t *testing.T, result bool)
	}{
		{
			name:      "ok",
			signature: "0x1ba84a5e7944426650938492b9ebf398587f601db6d92d75d50b36b549bc5f644804ae83e12ad686fe6641f414356e75dae96668aeeee48d433a63cd5c4e83171c",
			checkResponse: func(t *testing.T, result bool) {
				require.Equal(t, result, true)
			},
		},
		{
			name:      "invalid signature length",
			signature: "0x1ba84a5e7944426650938492b9ebf398587f601db6d92d75d50b36b549bc5f644804ae83e12ad686fe6641f414356e75dae96668aeeee48d433a63cd5c",
			checkResponse: func(t *testing.T, result bool) {
				require.Equal(t, result, false)
			},
		},
		{
			name:      "invalid",
			signature: "0x53b3255e84bbd3d9e3adfc1f6520fab4091c53c2125f186f21fb344dcedb736471c596490716b11b3a59dbcef1b478e6a9390488a8f176e542d7b6836099170a1b",
			checkResponse: func(t *testing.T, result bool) {
				require.Equal(t, result, false)
			},
		},
		{
			name:      "invalid signature type",
			signature: "53b3255e84bbd3d9e3adfc1f6520fab4091c53c2125f186f21fb344dcedb736471c596490716b11b3a59dbcef1b478e6a9390488a8f176e542d7b6836099170a1b0x",
			checkResponse: func(t *testing.T, result bool) {
				require.Equal(t, result, false)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateMessageSignature(
				"0x2C0b73164AF92a89d30Af163912B38F45b7f7b65",
				tc.signature,
				[]byte("joinhello"),
			)

			tc.checkResponse(t, result)
		})
	}
}
