package blockchain

import (
    "crypto/sha256"
    "encoding/binary"
    "encoding/hex"
    "errors"
    "fmt"
    "strings"
    "time"

    //"go.dedis.ch/onet/v3/network"
)

type MicroBlock struct {
    *BlockHeader
    Payload []Transaction
    PayRoll []string
    PubKey string
    Hash chain.Hash
    Signature string
}

type KeyBlock struct {
    Nonce uint64
    Timestamp time.Time
    Root chain.Hash
    Parent chain.Hash
    PubKey string
    Hash chain.Hash
}

type KeyBlockTree struct {
    Root chain.Hash
    Parent chain.Hash
    KeyBlocks []KeyBlock
    SignaturePairs []SignaturePair
}

// Deserialize decodes a block header from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field.
func (b *MicroBlock) Deserialize(r io.Reader) error {
    return readBlock(r, 0, b)
}

// Serialize encodes a block header from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field.
func (b *MicroBlock) Serialize(w io.Writer) error {
    return writeBlockHeader(w, 0, b)
}

// readBlock reads a micro  block from r.  See Deserialize for
// decoding block stored to disk, such as in a database, as opposed to
// decoding from the wire.
func readBlock(r io.Reader, pver uint32, b *MicroBlock) error {
    return readElements(r, &b.Version, &b.PrevBlock, &b.MerkleRoot,
	(*uint32Time)(&bh.Timestamp), &bh.Bits)
}

// writeBlock writes a bitcoin block to w.  See Serialize for
// encoding block to be stored to disk, such as in a database, as
// opposed to encoding for the wire.
func writeBlock(w io.Writer, pver uint32, b *MicroBlock) error {
	sec := uint32(bh.Timestamp.Unix())
	return writeElements(w, b.Version, &b.PrevBlock, &b.MerkleRoot,
		sec, b.Bits)
}

func (hb *NewLeaderHelloBlock) CalculateHash() (BlockID, error) {
    var err error
    hash := sha256.New()
    for _, val := range []uint32{hb.Version} {
        err = binary.Write(hash, binary.LittleEndian, val)
        if err != nil {
            return nil, errors.New("error writing to hash:" + err.Error())
        }
    }
    for _, val := range []uint64{hb.Index, hb.Timestamp} {
        err = binary.Write(hash, binary.LittleEndian, val)
        if err != nil {
            return nil, errors.New("error writing to hash:" + err.Error())
        }
    }
    hash.Write(hb.PrevBlock)
    hash.Write(hb.MerkleRoot)
    hash.Write([]byte(hb.PublicKey))
    buf := hash.Sum(nil)
    return buf, nil
}

func (hb *NewLeaderHelloBlock) Copy() *NewLeaderHelloBlock {
    prevBlock := make(BlockID, len(hb.PrevBlock))
    copy(prevBlock, hb.PrevBlock)
    merkleRoot := make(BlockID, len(hb.MerkleRoot))
    copy(merkleRoot, hb.MerkleRoot)
    hash := make([]byte, len(hb.Hash))
    copy(hash, hb.Hash)
    var transactions []*Transaction
    for _, transaction := range hb.Transactions {
        transactions = append(transactions, transaction) // Copy
    }
    return &NewLeaderHelloBlock{
        Index:		hb.Index,
        Version:	hb.Version,
        Timestamp:	hb.Timestamp,
        PrevBlock:	prevBlock,
        MerkleRoot:	merkleRoot,
        PublicKey:  hb.PublicKey,
        Hash:       hash,
        Block:      hb.Block.Copy(),
        Addresses:  hb.Addresses[:],
        Transactions: transactions,
    }
}

func (hb *NewLeaderHelloBlock) String() string {
    var builder strings.Builder
    builder.WriteString(fmt.Sprintf("Block %d", hb.Index))
    builder.WriteString(fmt.Sprintf("\n\tHeight: %d", hb.Index))
    builder.WriteString(fmt.Sprintf("\n\tVersion: %d", hb.Version))
    builder.WriteString(fmt.Sprintf("\n\tTimestamp: %s", time.Unix(0, int64(hb.Timestamp)).Format("2006-01-02 15:04:05")))
    builder.WriteString(fmt.Sprintf("\n\tPrevBlock: %s", hex.EncodeToString(hb.PrevBlock)))
    builder.WriteString(fmt.Sprintf("\n\tMerkleRoot: %s", hex.EncodeToString(hb.MerkleRoot)))
    builder.WriteString(fmt.Sprintf("\n\tPublicKey: %s", hb.PublicKey))
    builder.WriteString(fmt.Sprintf("\n\tAddresses: %v", hb.Addresses))
    builder.WriteString(fmt.Sprintf("\n\tHash: %s", hex.EncodeToString(hb.Hash)))
    return builder.String()
}

// DeFork including pow winner's block
type DeForkBlock struct {
    // Time the block was created.
    Timestamp uint64
    OrderBlocks []*Block
    Hash []byte
}

func NewDeForkBlock(blocks []*Block) *DeForkBlock {
    return &DeForkBlock{
        Timestamp: uint64(time.Now().UnixNano()),
        OrderBlocks: blocks,
    }
}

func (bh *DeForkBlock) Copy() *DeForkBlock {
    var orderBlocks []*Block
    for _, block := range bh.OrderBlocks {
        orderBlocks = append(orderBlocks, block.Copy())
    }
    deForkBlock := NewDeForkBlock(orderBlocks)
    //deForkBlock.Hash = make([]byte, len(hb.Hash))
    copy(deForkBlock.Hash, bh.Hash)
    return deForkBlock
}

type BlockHeader struct {
    // Version of the block.
    Version uint32
    // Difficulty target for the block.
    Bits uint32
    // Nonce used to generate the block.
    Nonce uint64
    // Time the block was created.
    Timestamp uint64
    // Hash of the previous block header in the block chain.
    PrevBlock BlockID
    // Merkle tree reference to hash of all transactions for the block.
    // MerkleRoot BlockID
    // Public Key
    PublicKey string

    // Data is any data to be stored in that Block.
    Data []byte
}

func (bh *BlockHeader) Copy() *BlockHeader {
    prevBlock := make(BlockID, len(bh.PrevBlock))
    copy(prevBlock, bh.PrevBlock)
    /*merkleRoot := make(BlockID, len(bh.MerkleRoot))
    copy(merkleRoot, bh.MerkleRoot)*/
    /*var addresses []*network.ServerIdentity
    for _, addr := range bh.Addresses {
        addresses = append(addresses, addr) // Copy
    }*/
    data := make([]byte, len(bh.Data))
    copy(data, bh.Data)
    return &BlockHeader{
	    Version:	bh.Version,
	    Bits:		bh.Bits,
	    Nonce:		bh.Nonce,
	    Timestamp:	bh.Timestamp,
	    PrevBlock:	prevBlock,
        PublicKey:  bh.PublicKey,
	    Data:		data,
    }
}

type Block struct {
    *BlockHeader
    Hash []byte
    Collections []*Collection
}

func NewBlock() *Block {
    return &Block{
        BlockHeader: &BlockHeader{
            Data: make([]byte, 0),
	    },
	    Hash: make([]byte, 0),
	    Collections: make([]*Collection, 0),
    }
}

// CalculateHash hashes all block header of the block.
func (b *Block) CalculateHash() (BlockID, error) {
    var err error
    hash := sha256.New()
    /*err := binary.Write(hash, binary.LittleEndian, int32(b.Index))
    if err != nil {
        return nil, errors.New("error writing to hash:" + err.Error())
    }*/
    for _, val := range []uint32{b.Version, b.Bits} {
        err = binary.Write(hash, binary.LittleEndian, val)
        if err != nil {
            return nil, errors.New("error writing to hash:" + err.Error())
        }
    }
    for _, val := range []uint64{b.Nonce, b.Timestamp} {
        err = binary.Write(hash, binary.LittleEndian, val)
        if err != nil {
            return nil, errors.New("error writing to hash:" + err.Error())
        }
    }

    hash.Write(b.PrevBlock)
    hash.Write([]byte(b.PublicKey))
    hash.Write(b.Data)
    buf := hash.Sum(nil)
    return buf, nil
}

// Copy makes a deep copy of the Block
func (b *Block) Copy() *Block {
    if b == nil {
        return nil
    }
    block := &Block{
        BlockHeader:	b.BlockHeader.Copy(),
	    Hash:           make([]byte, len(b.Hash)),
	    Collections:	make([]*Collection, len(b.Collections)),
    }
    copy(block.Hash, b.Hash)
    for _, c := range b.Collections {
        block.Collections = append(block.Collections, newCollection(c.PublicKey))
    }
    return block
}

func (b *Block) Sign(key string) {
    keyIndex := -1
    for index, c := range b.Collections {
        if c.PublicKey == key {
            keyIndex = index
            break
        }
    }
    if keyIndex < 0 {
        b.Collections = append(b.Collections, newCollection(key))
    }
}

func (b *Block) String() string {
    var builder strings.Builder
    builder.WriteString(fmt.Sprintf("\n\tVersion: %d", b.Version))
    builder.WriteString(fmt.Sprintf("\n\tBits: 0x%x", b.Bits))
    builder.WriteString(fmt.Sprintf("\n\tNonce: %d", b.Nonce))
    builder.WriteString(fmt.Sprintf("\n\tTimestamp: %s", time.Unix(0, int64(b.Timestamp)).Format("2006-01-02 15:04:05")))
    builder.WriteString(fmt.Sprintf("\n\tPrevBlock: %s", hex.EncodeToString(b.PrevBlock)))
    builder.WriteString(fmt.Sprintf("\n\tPublicKey: %s", b.PublicKey))
    builder.WriteString(fmt.Sprintf("\n\tData: %s", hex.EncodeToString(b.Data)))
    builder.WriteString(fmt.Sprintf("\n\tHash: %s", hex.EncodeToString(b.Hash)))
    return builder.String()
}
