package mining

import (
    "testing"

    bc "daolicloud/blockchain"
    "go.dedis.ch/onet/v3/log"
)

func TestSolveBlock(t *testing.T) {
    m := New(nil)
    log.SetDebugVisible(1)
    block := bc.GetGenesisBlock().Copy()
    block.Bits = bc.DEFAULT_BITS
    for i := uint32(0); i <= 0x100; i++ {
        t.Logf("Bits: %v", block.Bits)
        log.Infof("Bits: %v", block.Bits)
        m.solveBlock(block, make(chan struct{}))
    }
}
