package blockchain

import (

    //"go.dedis.ch/onet/v3/network"
)

var genesisBlock0Hash = BlockID([]byte{
    0x8d, 0x15, 0x7a, 0xca, 0x50, 0x1e, 0x91, 0x61,
    0x68, 0x9a, 0x05, 0x61, 0xe0, 0x78, 0x2e, 0x13,
    0xe9, 0x41, 0xc0, 0xec, 0x0e, 0x25, 0x28, 0x2c,
    0x88, 0xd7, 0xbe, 0xcd, 0x5c, 0x94, 0xb6, 0x56,
})

var genesisBlock1Hash = BlockID([]byte{
    0x8d, 0x15, 0x7a, 0xca, 0x50, 0x1e, 0x91, 0x61,
    0x68, 0x9a, 0x05, 0x61, 0xe0, 0x78, 0x2e, 0x13,
    0xe9, 0x41, 0xc0, 0xec, 0x0e, 0x25, 0x28, 0x2c,
    0x88, 0xd7, 0xbe, 0xcd, 0x5c, 0x94, 0xb6, 0x56,
})

var genesisKeyBlockHash = BlockID([]byte{
    0x8d, 0x15, 0x7a, 0xca, 0x50, 0x1e, 0x91, 0x61,
    0x68, 0x9a, 0x05, 0x61, 0xe0, 0x78, 0x2e, 0x13,
    0xe9, 0x41, 0xc0, 0xec, 0x0e, 0x25, 0x28, 0x2c,
    0x88, 0xd7, 0xbe, 0xcd, 0x5c, 0x94, 0xb6, 0x56,
})

func GetGenesisBlock0() *MicroBlock {
    genesisBlock0 := &MicroBlock{
        BlockHeader: &BlockHeader{
            Version: 1,
	    Timestamp: time.Unix(0x60d5a920, 0),
	    PrevBlock: chain.Hash{},
	    MerkleRoot: chain.Hash{},
	    Bits: 0x1d00ffff,
        }
	Payload: []Transaction{},
	PayRoll: []string{},
	PubKey: "b2aa3a0faf75e5b09f048a30361a41d380c999531dd64123ba699fa4b9bcdcb7",
	Hash: genesisBlock0Hash,
    }
    return genesisBlock0
}

func GetGenesisBlock1() *MicroBlock {
    genesisBlock1 := &MicroBlock{
        BlockHeader: &BlockHeader{
            Version: 1,
	    Timestamp: time.Unix(0x60d5a920, 0),
	    PrevBlock: genesisBlock0Hash,
	    MerkleRoot: chain.Hash{},
	    Bits: 0x1d00ffff,
        }
	Payload: []Transaction{},
	PayRoll: []string{},
	PubKey: "b2aa3a0faf75e5b09f048a30361a41d380c999531dd64123ba699fa4b9bcdcb7",
	Hash: genesisBlock1Hash,
    }
    return genesisBlock1
}

func GetGenesisKeyBlock() *KeyBlock {
    keyBlock := &KeyBlock{
        Nonce: 12345,
	Timestamp: time.Unix(0x60d5a920, 0),
	Root: genesisBlock0Hash,
	Parent: genesisBlock1Hash,
	PubKey: "b2aa3a0faf75e5b09f048a30361a41d380c999531dd64123ba699fa4b9bcdcb7",
    }
    keyBlock.hash = genesisKeyBlockHash
    return keyBlock
}
