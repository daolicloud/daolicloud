package  blockchain

// blockStatus is a bitfield representing the validation state of the block.
type blockStatus byte

const (
    // statusNone indicates that the block has no validation state flags set.
    //
    // NOTE: This must be defined last in order to avoid influencing iota.
    statusNone blockStatus = 0
)
