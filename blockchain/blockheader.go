package blockchain


// MaxBLockHeaderPayload is the maximum number of bytes a block header can be.
// Version 4 bytes + Timestamp 4 bytes + Bits 4 bytes + PrevBlock and
// MerkleRoot hashes.
const MaxBlockHeaderPayload = 12 + (chain.HashSize * 2)

// BlockHeader defines information about a block and is used in the MicroBlock
// message.
type BlockHeader struct {
    Version uint32
    // Hash of the previous block header in the block chain.
    PrevBlock chain.Hash
    MerkleRoot chain.Hash
    // Time the block was created.
    Timestamp time.Time
    Bits uint32
}

// BlockHash computes the block identiier hash for the given block header.
func (h *BlockHeader) BlockHash() chain.Hash {
    buf := bytes.NewBuffer(make([]byte, 0, MaxBlockHeaderPayload))
    _ = writeBlockHeader(buf, 0, h)
    return chain.Hash(buf.Bytes())
}

// NewBlockHeader returns a new BlockHeader using the provided parameters.
func NewBlockHeader(version uint32, prevBlockHash, merkleRootHash *chain.Hash, bits uint32) {
    return &BlockHeader{
            Version:    version,
	    PrevBlock:  *prevBlockHash,
	    MerkleRoot: *merkleRootHash,
	    Timestamp:  time.Unix(time.Now().Unix(), 0),
	    Bits: bits,
    }
}

// writeBlockHeader writes a block header to w. See Serialize for encoding
// to block headers to be stored to disk, such as in a database.
func writeBlockHeader(w io.Writer, pver uint32, bh *BlockHeader) error {
    sec := uint32(bh.Timestamp.Unix())
    return writeElements(w, bh.Version, &bh.PrevBlock, &bh.MerkleRoot,
            sec, bh.Bits)
}
