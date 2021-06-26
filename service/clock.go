package service

import (
    "container/list"
    "sync"

    "golang.org/x/xerrors"
)

type PrivateClock struct {
    *list.List
    sync.Mutex
    Size int
}

func newPrivateClock(size int) *PrivateClock{
    return &PrivateClock{
        List: list.New(),
	    Size: size,
    }
}

func (pc *PrivateClock) lockLen() int {
    pc.Lock()
    defer pc.Unlock()
    return pc.Len()
}

func (pc *PrivateClock) Push(val uint64) {
    pc.Lock()
    defer pc.Unlock()
    if pc.lockLen() >= pc.Size {
        for i := 0; i <= pc.lockLen() - pc.Size; i++ {
            pc.Remove(pc.Front())
	}
    }
    pc.PushBack(val)
}

func (pc *PrivateClock) Clock() uint64{
    pc.Lock()
    defer pc.Unlock()
    if pc.lockLen() == 0 {
        return 0
    }
    return (pc.Back().Value.(uint64) - pc.Front().Value.(uint64)) / uint64(pc.lockLen())
}

// GC_Cycle(TB)=timestamp(TB)-timestamp((TB-w+1-i)/w+i
// PrivateCycles_N(TB)=PrivateClock_N(TB)-PrivateClock_N((TB-w+1-i)/w+i
// \Delta_N(TB,Evt)=GC_Cycle(TB)/PrivateCycles_N(TB)(PrivateClock_N(Evt)-PrivateClock_N(TB))
func (s *Service) calcDelta() (uint64, error) {
    latestBlock := s.db.GetLatest()
    preIndex := latestBlock.Index - uint64(s.privateClock.Size)
    div := uint64(s.privateClock.Size)
    if div == 0 {
        return 0, xerrors.New("size could not be zero")
    }
    if preIndex < 0 {
        preIndex = 0
	div = latestBlock.Index
    }
    block, err := s.db.GetBlockByIndex(preIndex)
    if err != nil {
        return 0, err
    }
    globalClock := (latestBlock.Timestamp - block.Timestamp) / div
    privateClock := s.privateClock.Clock()
    if privateClock == 0 {
        return 0, xerrors.New("private clock could not be zero")
    }
    return globalClock / privateClock, nil
}

