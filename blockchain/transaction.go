package blockchain

type Transaction struct {
    Version uint32
    Signature []byte
    PubKeyIn string
    PubKeyOut string
    Value uint64
    Timestamp uint64
    Nonce uint64
    // Store any data.
    Data []byte `protobuf:"bytes,opt"`
}
