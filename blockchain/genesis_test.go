package blockchain

import (
    "testing"
)

func TestNewGenesisBlock(t *testing.T) {
    hash, err := GetGenesisBlock().CalculateHash()
    if err != nil {
        t.Errorf(err.Error())
    } else {
        t.Logf("Genesis Block hash: %#x", hash)
    }
}
