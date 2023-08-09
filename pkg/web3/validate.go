package web3

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func ValidateMessageSignature(walletAddress, signature string, message []byte) bool {
	if signature[0:2] != "0x" || len(signature) != 132 {
		return false
	}

	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false
	}

	message = accounts.TextHash(message)
	sig[crypto.RecoveryIDOffset] -= 27

	recovered, err := crypto.SigToPub(message, sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	if walletAddress == recoveredAddr.Hex() {
		return true
	} else {
		return false
	}
}
