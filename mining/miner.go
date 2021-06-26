package mining

import (
    "sync"
    "time"

    bc "daolicloud/blockchain"

    "go.dedis.ch/onet/v3/log"
)

const (
    // maxNonce is the maximum value a nonce can be in a block header.
    maxNonce = ^uint32(0) // 2^32 - 1
)

type Listener func(block *bc.Block)

type Miner struct {
    sync.Mutex
    started           bool
    block	      *bc.Block
    quit              chan struct{}
    callback	      Listener
}

func (m *Miner) Start(block *bc.Block) {
    m.Lock()
    defer m.Unlock()

    // Nothing to do if the miner is already running
    if m.started {
        return
    }

    m.block = block
    m.quit = make(chan struct{})
    go m.miningWorker()
    m.started = true
    log.Infof("Miner started")
}

func (m *Miner) miningWorker() {
    quit := make(chan struct{})
    go m.generateBlocks(quit)
outer:
    for {
        select {
        case <-m.quit:
            close(quit)
            break outer
        }
    }
    log.Infof("Miner worker finished")
}

func (m *Miner) generateBlocks(quit chan struct{}) {
    log.Info("Starting generate blocks worker")
	newBlock := m.block
outer:
    for {
        select {
        case <-quit:
	        break outer
        default:
            // Non-blocking select to fall through
        }
	    if m.solveBlock(newBlock, quit) {
            m.callback(newBlock)
	        break outer
        }
    }
    log.Info("Generate blocks worker done")
}

func (m *Miner) solveBlock(block *bc.Block, quit chan struct{}) bool {
    header := block.BlockHeader
    //header.Index = header.Index + 1
    targetDifficulty := bc.CompactToBig(header.Bits)
    log.Lvlf3("Bits 0x%x convert to difficulty 0x%064x", header.Bits, targetDifficulty)
    // Search through the entire nonce range for a solution while
    // periodically checking for early quit and stale block
    // conditions along with updates to the speed monitor.
    for i := uint32(0); i <= maxNonce; i++ {
        select {
        case <-quit:
            log.Lvl3("solveBlock quit")
            return false
	    default:
            // Non-blocking select to fall through
        }
	    header.Nonce = uint64(i)
        header.Timestamp = uint64(time.Now().UnixNano())
	    hash, err := block.CalculateHash()
        if err != nil {
	        log.Warn(err)
            return false
        }
	    // The block is solved when the new block hash is less
        // than the target difficulty.  Yay!
        if bc.HashToBig(&hash).Cmp(targetDifficulty) <= 0 {
            //log.Lvlf3("Success Info: Times=%d,Height=%d,Hash=%064x", i, header.Index, hash)
            log.Lvlf3("Success Info: Times=%d,Hash=%064x", i, hash)
	    block.Hash = hash
            return true
        }
    }
    return false
}

// This function is safe for concurrent access.
func (m *Miner) Stop() {
    m.Lock()
    defer m.Unlock()

    // Nothing to do if the miner is not currently running
    if !m.started {
        return
    }

    log.Infof("Miner stopped")
    close(m.quit)
    m.started = false
}

func New(callback Listener) *Miner {
    return &Miner{
        callback: callback,
    }
}
