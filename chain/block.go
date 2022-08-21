package chain

import (
	"sync"
	"sync/atomic"
)

/*
	Instances of BlockProcessor simply count the number of occurences of each block with the same height.
	If 3 or more blocks have the same height and the height immediately prior to it, stored in `maxHeight`
	is accepted, then the block is considered accepted for that height and `maxHeight` is increased by 1.
	The number of occurrences of a block with the same height is tracked in `blockTracker`.

	See SOLUTION.md for details on assumptions made.

	There are a few ways to track counts of blocks per height, the simplest would be to have a hash map
	with block ID + height as a compound key, mapped to the number of occurrences of that combination.

	However, in a closer-to-real-world scenario, the map would continuously increase in size until we run out of resources.

	I decided then track blocks and heights at 2 levels. First with a hash map from height -> blocks, where blocks is another
	hash map of blocks -> occurrences

	This allows us to implement a process to remove entire blocks from `blockTracker` if they cannot be accepted anymore
	for being below `maxHeight` (assumed different blocks cannot share the same height and IDs are unique).
*/

const minAcceptedBlockCount = 3

// blockTracker counts how many times we've seen a block
type blockCounter sync.Map

// BlockProcessor validates and accepts blockchain transactions and tracks accepted block height
type BlockProcessor struct {
	// atomically stores max height
	maxHeight uint64

	// concurrently tracks heights and blocks
	blockTracker sync.Map
}

// NewBlockProcessor creates a new block processor with an accepted genesis block
func NewBlockProcessor() *BlockProcessor {
	return &BlockProcessor{
		maxHeight:    0, // genesis block is accepted
		blockTracker: sync.Map{},
	}
}

// ProcessBlocks consumes a sequence of block transactions ids of certain height and
// returns the maximum accepted height.
// If none of the transactions can be validated, the method returns the most current
// maximum accepted height
func (p *BlockProcessor) ProcessBlocks(startHeight uint64, blocks []string) uint64 {
	blockHeight := startHeight
	for _, block := range blocks {

		// Here I try optimistic improvement that skips heights that are less or equal to p.maxHeight
		// This is more to show that further optiization  might be possible given the behaviour of the system in real world
		// In reality in a high concurrency environment this might not have any benefit, atomic calls can be expensive
		// and the only way to determine that would be to properly benchmark different implementations
		currentMaxHeight := atomic.LoadUint64(&p.maxHeight)
		if blockHeight > currentMaxHeight {
			p.processBlock(blockHeight, block)
		}
		blockHeight++
	}
	return atomic.LoadUint64(&p.maxHeight)
}

func (p *BlockProcessor) processBlock(height uint64, block string) {
	// this could be optimized, we're always creating a sync.Map here but might not always
	// need it and it would add memory pressure
	bc := &sync.Map{}

	// find all blocks that could have this same height. If none exists yet, add a new entry
	if bcValue, bcLoaded := p.blockTracker.LoadOrStore(height, bc); bcLoaded {
		bc = bcValue.(*sync.Map)
	}

	// we need to store pointers to values so we can increment them atomically
	// without having to modify the map again
	currentCount := startBlockCountPtr()
	// when the block coutner entry is created, it starts at 1
	if counterValue, counterLoaded := bc.LoadOrStore(block, currentCount); counterLoaded {
		currentCount = counterValue.(*uint64)

		// if we loaded, we need to increment the counter
		atomic.AddUint64(currentCount, 1)
	}

	// A block can be accepted now so we can update maxHeight
	if *currentCount >= minAcceptedBlockCount {
		p.updateMaxHeight(height)
	}
}

func (p *BlockProcessor) updateMaxHeight(height uint64) {
	// Try to update maxHeight but it's possible another block was accepted for the same height in another thread
	// and the block that got us here cannot be accepted anymore for this height

	for {
		// Given the assumption different blocks cannot be accepted with the same height,
		// height must be exactly equal to p.maxHeight + 1
		curMaxHeight := atomic.LoadUint64(&p.maxHeight)
		if curMaxHeight > height || (height-curMaxHeight) != 1 {
			return
		}

		// if the height we're trying to set as maximum maxHeight was incremented by another thread, we cannot use this height anymore
		// and need to retry. Otherwise, we're good to return with a new maxHeight set
		if atomic.CompareAndSwapUint64(&p.maxHeight, p.maxHeight, height) {
			return
		}
	}
}

func startBlockCountPtr() *uint64 {
	var startCount uint64 = 1
	return &startCount
}
