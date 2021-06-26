package service

import "time"

func getCurrentTimestamp() uint64 {
    return uint64(time.Now().UnixNano())
}
