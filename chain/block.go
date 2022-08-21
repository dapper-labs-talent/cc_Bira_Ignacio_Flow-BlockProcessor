package chain

// BlockProcessor validates and accepts blockchain transactions and tracks accepted block height
type BlockProcessor struct {
}

// NewBlockProcessor creates a new block processor with an accepted genesis block
func NewBlockProcessor() *BlockProcessor {
	return &BlockProcessor{}
}

// ProcessBlocks consumes a sequence of block transactions ids of certain height and
// returns the maximum accepted height.
// If none of the transactions can be validated, the method returns the most current
// maximum accepted height
func (p *BlockProcessor) ProcessBlocks(startHeight uint64, blocks []string) uint64 {

	return 0
}

// MaxAcceptedHeight returns the last maximum accepted height
func (p *BlockProcessor) MaxAcceptedHeight() uint64 {
	return 0
}
