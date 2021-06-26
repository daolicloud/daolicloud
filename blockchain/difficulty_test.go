package blockchain

import (
    "fmt"
    "math/big"
    "time"
    "testing"
)


func TestCompactToBigPrint(t *testing.T) {
    fmt.Printf("New target %08x (%064x)", DEFAULT_BITS, CompactToBig(DEFAULT_BITS))
}

func TestHashToBig(t *testing.T) {
    block := GetGenesisBlock().Copy()
    for i := uint32(0); i <= (^uint32(0)); i++ {
        block.Nonce = uint64(i)
        block.Timestamp = uint64(time.Now().UnixNano())
	    hash, err := block.CalculateHash()
        if err != nil {
            t.Errorf("%+v", err)
        }
        hashValue := HashToBig(&hash)
        diffValue := CompactToBig(DEFAULT_BITS)
        //t.Logf("\nHash=%064x\nDiff=%064x", hashValue, diffValue)
        if hashValue.Cmp(diffValue) <= 0 {
            fmt.Printf("^^^^^%v times^^^^^^^", i)
            break
        }
    }
}

// This example demonstrates how to convert the compact "bits" in a block header
// which represent the target difficulty to a big integer and display it using
// the typical hex notation.
func TestCompactToBig(t *testing.T) {
    // Convert the bits from block 300000 in the main block chain.
    bits := uint32(DEFAULT_BITS)
    targetDifficulty := CompactToBig(bits)

    // Display it in hex.
    fmt.Printf("%064x\n", targetDifficulty.Bytes())

    // Output:
    // 0000861000000000896c00000000000000000000000000000000000000000000
    // 0000016800000000896c00000000000000000000000000000000000000000000
}

// This example demonstrates how to convert a target difficulty into the compact
// "bits" in a block header which represent that target difficulty .
func TestBigToCompact(t *testing.T) {
    // Convert the target difficulty from block 300000 in the main block
    // chain to compact form.
    //bigT := "0000000000000000896c00000000000000000000000000000000000000000000"
    bigT := "0000016800000000896c00000000000000000000000000000000000000000000"
    targetDifficulty, success := new(big.Int).SetString(bigT, 16)
    if !success {
        fmt.Println("invalid target difficulty")
        return
    }
    bits := BigToCompact(targetDifficulty)

    fmt.Println(bits)

    // Output:
    // 503408640
}
