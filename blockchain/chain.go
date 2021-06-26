package blockchain

import (
    "encoding/hex"
)

// HashSize of array used to store hases.
const HashSize = 32

var (
    // blockIndexBucketName is the name of the db bucket used to house to the
    // block headers and contextual information.
    blockIndexBucketName = []byte("blockheaderidx")
    // chainStateKeyName is the name of the db key used to store the best
    // chain state.(
    chainStateKeyName = []byte("chainstate")
)

// Hash is used in serval of the bitcoin messages and common structures. It
// typically represents the double sha256 of data.
type Hash [HashSize]byte

// String returns the Hash as the hexadecimal string of the byte-reversed
// hash.
func (hash Hash) String() string {
    for i :=0 ; i < HashSize/2; i++ {
        hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
    }
    return hex.EncodeToString(hash[:])
}

// initChainState attempts to loaad and initialized the chain state from the
// database. When the db does no yet contain any chain state, both it and the
// chain state are initialized to the genesis block.
func (b *BlockChain) initChainState() error {
   var initialized, hasBlockIndex bool
   err := b.db.View(func(tx *bbolt.Tx) error {
       initialized = tx.Get(chainStateKeyName) != nil
       hasBlockIndex = tx.Bucket(blockIndexBucketName) != nil
       return nil
   })
   if err != nil {
        return err
   }

   if !initialized {
       // At this point the database has not already been initialized, so
       // initialize bth it and the chain state to the genesis block.
       return b.createChainState()
   }

   if !hasBlockIndex {
       err := migrateBlockIndex(b.db)
       if err != nil {
           return nil
       }
   }

   // Attempt to load the chain state from the database.
   err = b.db.View(

        log.Infof("Loading block index...")
	blockIndexBucket = db.Bucket(blockIndexBucketName)

	var i int32
	var lastNode *blockNode
	cursor := blockIndexBucket.Cursor()
	for ok := cursor.First(); ok; ok = cursor.Next() {
            block, status, err := deserializeBlockRow(cursor.Value())
	    if err != nil {
                return err
	    }

	    // Determine the parent block node. Since we iterate block headers
	    // in order of height, if the blocks are mostly linear there is a
	    // very good chance the previous header processed is the parent.
	    var parent *blockNode
	    if lastNode == nil {
	    	blockHash := block.BlockHash()
	    	if !blockHash.IsEqual(b.GenesisHash0) {
	    		return AssertError(fmt.Sprintf("initChainState: Expected "+
	    			"first entry in block index to be genesis block, "+
	    			"found %s", blockHash))
	    	}
	    } else if block.PrevBlock == lastNode.hash {
	    	// Since we iterate block headers in order of height, if the
	    	// blocks are mostly linear there is a very good chance the
	    	// previous header processed is the parent.
	    	parent = lastNode
	    } else {
	    	parent = b.index.LookupNode(&header.PrevBlock)
	    	if parent == nil {
	    		return AssertError(fmt.Sprintf("initChainState: Could "+
	    			"not find parent for block %s", header.BlockHash()))
	    	}
	    }

	    // Initialize the block node for the block, connect it,
	    // and add it to the block index.
	    node := new(blockNode)
	    //initBlockNode(node, header, parent)
	    node.status = status
	    //b.index.addNode(node)

	    lastNode = node
	    i++
	}
   }
}

// deserializeBlockRow parses a value in the block index bucket into a block
// header and block status bitfield.
func deserializeBlockRow(blockRow []byte) (*MicroBlock, blockStatus, error) {
    buffer := bytes.NewReader(blockRow）

    var block MicroBlock
    err := block.Deserialize(buffer)
    if err != nil {
        return nil, statusNone, err
    }

    statusByte, err := buffer.ReadByte()
    if err != nil {
        return nil, statusNone, err
    }

    return &block, blockStatus(statusByte), nil
}

// New returns a BlockChain instance.
func New(db *bbolt.DB) (*BlockChain, error) {
    if db == nil {
        return nil, AssertError("blockchain.New database s nil")
    }
    bc := BlockChain{
        db:	db,
    }
    // Initialize the chain state from the passed database.When the db
    // does not yet contain any chain state, both it and the chain state
    // will be initialized to contain only the genesis block.
    if err := bc.initChainState(); err != nil {
        return nil, err
    }
    return &bc, nil
}
