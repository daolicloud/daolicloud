package blockchain

import (
    "go.dedis.ch/onet/v3/network"
)

const HashSize = 32

const DEFAULT_BITS = 0x1e016800

// BlockID represents the Hash of the Block
type BlockID []byte

// Hash is used in several of the bitcoin messages and common structures.  It
// typically represents the double sha256 of data.
type Hash [HashSize]byte

type SignaturePair struct {
    PublicKey string
    Signature string
}

func newSignaturePair(publicKey, signature string) *Collection {
    return &SignaturePair{
        PublicKey: publicKey,
	Signature: signature,
    }
}

func init() {
    network.RegisterMessage(&BlockHeader{})
    network.RegisterMessage(&Block{})
}

// BlockChain provides functions for working with the block chain.
type BlockChain struct {
    db	*bbolt.DB
}
