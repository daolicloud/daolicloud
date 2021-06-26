package service

import (
	//"bytes"
	//"math/rand"
    "sync"

    bc "daolicloud/blockchain"

    bbolt "go.etcd.io/bbolt"
    "go.dedis.ch/onet/v3/log"
    "go.dedis.ch/onet/v3/network"
    "golang.org/x/xerrors"
)

type BlockID bc.BlockID

type txPool struct {
    sync.Mutex
    txsMap map[string][]bc.Transaction
}

func newTxPool() txPool {
    return txPool{
        txsMap: make(map[string][]bc.Transaction),
    }
}

func (p *txPool) add(key string, tx bc.Transaction) {
    p.Lock()
    defer p.Unlock()

    if txs, ok := p.txsMap[key]; !ok {
	p.txsMap[key] = []bc.Transaction{tx}
    } else {
	txs = append(txs, tx)
	p.txsMap[key] = txs
    }
}

// Announce is used to pass a message to all children.
type PingMessage struct {
    Message string
}

type TunnelMessage PingMessage

type BlockMessage struct {
    Type int // 0:CandidateBlock,1:RefererBlock,2:TxBlock
    Block interface{}
}

type HandshakeMessage struct {
    GenesisID BlockID
    LatestBlock *bc.NewLeaderHelloBlock
    Answer  bool
}

type DownloadBlockRequest struct {
    GenesisID BlockID
    Start uint64
    Size uint64
}

type DownloadBlockResponse struct {
    Blocks []*bc.NewLeaderHelloBlock
    GenesisID BlockID
}

type SignatureRequest struct {
	Blocks []*bc.Block
}

// SignatureResponse is what the Cosi service will reply to clients.
type SignatureResponse struct {
    PublicKey string
    BlockMap map[string]*bc.Block
}

type ProxyRequest struct {
	Blocks []*bc.Block
}

type ProxyResponse struct {
    PublicKey string
    BlockMap map[string]*bc.Block
}

type AddressMessage struct {
    Addresses []string
}

type RemoteServerIndex struct {
    Index uint64
    ServerIdentity *network.ServerIdentity
}

type BlockDB struct {
    *bbolt.DB
    bucketName []byte
    latestBlockID BlockID
    latestMutex sync.Mutex
    blockIndexMap map[uint64]BlockID
}

// NewBlockDB returns an initialized BlockDB structure.
func NewBlockDB(db *bbolt.DB, bn[]byte) *BlockDB {
	return &BlockDB{
		DB: db,
		bucketName: bn,
		blockIndexMap: make(map[uint64]BlockID),
	}
}

// storeToTx stores the block into the database.
func (db *BlockDB) storeToTx(tx *bbolt.Tx, block *bc.NewLeaderHelloBlock) (BlockID, error) {
    log.Lvlf2("Storing block %d / %x", block.Index, block.Hash)
	val, err := network.Marshal(block)
	if err != nil {
		return nil, err
	}
	return block.Hash, tx.Bucket([]byte(db.bucketName)).Put(block.Hash, val)
}

// getFromTx returns the skipblock identified by sbID.
// nil is returned if the key does not exist.
func (db *BlockDB) getFromTx(tx *bbolt.Tx, sbID BlockID) (*bc.NewLeaderHelloBlock, error) {
	if sbID == nil {
		return nil, xerrors.New("cannot look up block with ID == nil")
	}

	val := tx.Bucket(db.bucketName).Get(sbID)
	if val == nil {
		return nil, nil
	}

	// For some reason boltdb changes the val before Unmarshal finishes. When
	// copying the value into a buffer, there is no SIGSEGV anymore.
	buf := make([]byte, len(val))
	copy(buf, val)
	_, sbMsg, err := network.Unmarshal(buf, suite)
	if err != nil {
		return nil, err
	}

	return sbMsg.(*bc.NewLeaderHelloBlock).Copy(), nil
}

// Stores the set of blocks in the boltdb
func (db *BlockDB) StoreBlocks(blocks []*bc.NewLeaderHelloBlock) ([]BlockID, error) {
    var result []BlockID
    err := db.Update(func(tx *bbolt.Tx) error {
        for _, block := range blocks {
	   blockID, err := db.storeToTx(tx, block)
	   if err != nil {
	        return err
	   }
	   result = append(result, blockID)
	}
        return nil
    })
    return result, err
}

// Stores the given Block into db
func (db *BlockDB) Store(block *bc.NewLeaderHelloBlock) BlockID {
	ids, err := db.StoreBlocks([]*bc.NewLeaderHelloBlock{block})
	if err != nil {
		log.Error(err)
	}
	if len(ids) > 0 {
		return ids[0]
	}
	return nil
}

// GetGenesisID
func (db *BlockDB) GetGenesisID() (BlockID) {
    blockID, ok := db.blockIndexMap[0]
    if !ok {
	log.Error("no genesis block found")
        return nil
    }
    return blockID
}


// GetLatest searches for the latest available block.
func (db *BlockDB) GetLatest() (*bc.NewLeaderHelloBlock) {
    db.latestMutex.Lock()
    latest := db.GetByID(db.latestBlockID)
    db.latestMutex.Unlock()
    return latest
}

func (db *BlockDB) UpdateLatest(block *bc.NewLeaderHelloBlock) {
    db.latestMutex.Lock()
    db.latestBlockID = BlockID(block.Hash)
	db.blockIndexMap[block.Index] = BlockID(block.Hash)
    db.latestMutex.Unlock()
}

// GetByID returns a new copy of the skip-block or nil if it doesn't exist
func (db *BlockDB) GetByID(bID BlockID) *bc.NewLeaderHelloBlock {
    var result *bc.NewLeaderHelloBlock
    if bID == nil {
        return nil
    }
    err := db.View(func(tx *bbolt.Tx) error {
        block, err := db.getFromTx(tx, bID)
	    if err != nil {
	        return err
	    }
	    result = block
	    return nil
    })

    if err != nil {
        log.Error(err)
    }
    return result
}

// Build block index
func (db *BlockDB) BuildIndex() error {
    return db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket(db.bucketName)
        if b == nil {
            return xerrors.New("Missing bucket")
        }
        return b.ForEach(func(k []byte, v []byte) error {
            _, msg, err := network.Unmarshal(v, suite)
	    if err != nil {
                return err
            }
	    block := msg.(*bc.MicroBlock).Copy()
	    log.Lvlf3("Loading block %d / %x", block.Index, block.Hash)
            //db.blockIndexMap[block.Index] = BlockID(k)
            db.UpdateLatest(block)
	    return nil
	})
    })
}

func (db *BlockDB) GetBlockByIndex(blockIndex uint64) (*bc.NewLeaderHelloBlock, error) {
    blockID, ok := db.blockIndexMap[blockIndex]
    if !ok {
        return nil, xerrors.New("no block found")
    }
    block := db.GetByID(blockID)
    return block, nil
}

// blockBuffer will cache a referer block when the conode has
// verified it and it will later store it in the DB after the protocol
// has succeeded.
// Blocks are stored per skipchain to prevent the cache to grow and it
// is cleared after a new block is added
type blockBuffer struct {
    blockMap map[string]*bc.Block
	sync.Mutex
}

func newBlockBuffer() *blockBuffer {
    return &blockBuffer{
    	blockMap: make(map[string]*bc.Block),
	}
}

func (bb *blockBuffer) List() []*bc.Block {
	var blocks []*bc.Block
	for _, block := range bb.blockMap {
		blocks = append(blocks, block)
	}
	return blocks
}

func (bb *blockBuffer) Put(block *bc.Block) {
    bb.Lock()
    defer bb.Unlock()
    //bb.blocks = append(bb.blocks, block.Copy())
    bb.blockMap[string(block.Hash)] = block
}

func (bb *blockBuffer) Len() int {
	return len(bb.blockMap)
}

func (bb *blockBuffer) Set(block *bc.Block) {
	bb.Lock()
	defer bb.Unlock()
	if _, ok := bb.blockMap[string(block.Hash)]; !ok {
		bb.blockMap[string(block.Hash)] = block
	} else {
		bb.blockMap[string(block.Hash)].Sign(string(block.Hash))
	}
}

func (bb *blockBuffer) HalfMore(memberLen uint64) []*bc.Block {
	var blocks []*bc.Block
	for _, block := range bb.blockMap {
		if len(block.Collections) > int(memberLen) {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func (bb *blockBuffer) in(hash BlockID) bool {
	bb.Lock()
	defer bb.Unlock()
	_, ok := bb.blockMap[string(hash)]
	return ok
	/*for _, b := range bb.blocks {
		if bytes.Compare(b.Hash, hash) == 0 {
			return true
		}
	}
	return false*/
}

/*func (bb *blockBuffer) Choice() []*bc.Block {
    bb.Lock()
    defer func() {
	    bb.blocks = bb.blocks[:0]
	    bb.Unlock()
    }()
    if bb.Len() > 0 {
        selectedIndex := rand.Intn(len(bb.blocks))
        bb.blocks[0], bb.blocks[selectedIndex] = bb.blocks[selectedIndex], bb.blocks[0]
    }
    return bb.blocks
}

func (bb *blockBuffer) ChoiceOne() *bc.Block {
    bb.Lock()
    defer func() {
	    bb.blocks = bb.blocks[:0]
	    bb.Unlock()
    }()
    if len(bb.blocks) == 0 {
        return nil
    }
    return bb.blocks[rand.Intn(len(bb.blocks))]
}*/
