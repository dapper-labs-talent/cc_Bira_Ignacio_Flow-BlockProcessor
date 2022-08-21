package chain

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// simple type redefinition to make testing less cluttered
type testBlockList []string

// Test notes
// Note that I've avoided table based tests explicitly
// so as to make tests easier to read at a glance

func TestProcessorStartsWithacceptedGenesisBlock(t *testing.T) {
	processor := NewBlockProcessor()

	assert.EqualValues(t, 0, processor.ProcessBlocks(1, []string{}))
}

func TestEmptyBlocksKeepHeight(t *testing.T) {
	height := processTestBlocks(t, []string{})
	assert.EqualValues(t, 0, height)
}

func TestDifferentBlocksNotaccepted(t *testing.T) {
	height := processTestBlocks(t, []string{"a", "b", "c"})
	assert.EqualValues(t, 0, height)
}

func TestEmptyStringBlockIdsNotAccepted(t *testing.T) {
	height := processTestBlocks(t,
		[]string{"", "a"},
		[]string{""},
		[]string{"", "b", "c"},
	)
	assert.EqualValues(t, 0, height)
}

func TestSameBlockDifferentHeightsNotAccepted(t *testing.T) {
	height := processTestBlocks(t, []string{"a", "a", "a"})
	assert.EqualValues(t, 0, height)
}

func TestFirstAcceptedBlock(t *testing.T) {
	// the first block height after the genesis block is accepted
	// the next accepted block is at height 1
	const expectedHeight uint64 = 1
	height := processTestBlocks(t,
		[]string{"accepted", "not accepted"},
		[]string{"accepted", "also not accepted"},
		[]string{"accepted", "not accepted again"},
	)

	assert.Equal(t, expectedHeight, height)
}

func TestThreeConfirmationsNeededForBlockToBeaccepted(t *testing.T) {
	const expectedHeight uint64 = 0
	height := processTestBlocks(t,
		[]string{"accepted", "not accepted"},
		[]string{"accepted", "also not accepted"},
	)

	assert.Equal(t, expectedHeight, height)
}

func TestMoreThanThreeConfirmationsAcceptBlock(t *testing.T) {
	const expectedHeight uint64 = 1
	height := processTestBlocks(t,
		[]string{"accepted", "1"},
		[]string{"aceppted"},
		[]string{"accepted", "b"},
		[]string{"accepted", "apple", "carrot"},
	)

	assert.Equal(t, expectedHeight, height)
}

func TestBlockSameIdDifferentHeightNotAccepted(t *testing.T) {
	height := processTestBlocks(t,
		[]string{"same id", "1"},
		[]string{"same id", "a", "2"},
		[]string{"3", "same id"},
	)

	assert.EqualValues(t, 0, height)
}

func TestBlocksAcceptedAsHeightIncreasesWhileProcessing(t *testing.T) {
	const expectedHeight uint64 = 3
	var height uint64 = 1

	// block 'a' will be accepted and even though 'b' has the same height 3 times
	// its height will be same as max height and not accepted.
	// block 'c' will have height 2 and will be accepted because 'a' was accepted with 1
	processor := NewBlockProcessor()
	blocks := []string{"a", "b", "c"}
	processor.ProcessBlocks(1, blocks)
	processor.ProcessBlocks(1, blocks)
	height = processor.ProcessBlocks(1, blocks)

	assert.Equal(t, expectedHeight, height)
}

// Test with a little more elaborate setup
func TestBlockIdSamePositionDifferentStartHeightNotAccepted(t *testing.T) {
	processor := NewBlockProcessor()
	processor.ProcessBlocks(1, []string{"a"})
	processor.ProcessBlocks(2, []string{"a"})
	height := processor.ProcessBlocks(1, []string{"a"})

	assert.EqualValues(t, 0, height)
}

func TestBlockNotAcceptedIfPreviousHeightNotAccepted(t *testing.T) {
	// BlockProcess only knows about the genesis block at this point
	processor := NewBlockProcessor()
	processor.ProcessBlocks(2, []string{"a"})
	processor.ProcessBlocks(2, []string{"a"})
	height := processor.ProcessBlocks(2, []string{"a"})

	assert.EqualValues(t, 0, height)
}

func TestTwoBlocksNotAcceptedForSameHeight(t *testing.T) {
	var height uint64 = 1
	processor := NewBlockProcessor()

	// accepted one block
	processor.ProcessBlocks(1, []string{"a"})
	processor.ProcessBlocks(1, []string{"a"})
	height = processor.ProcessBlocks(1, []string{"a"})
	assert.EqualValues(t, 1, height)

	// and lets accept another
	processor.ProcessBlocks(2, []string{"b"})
	processor.ProcessBlocks(2, []string{"b"})
	height = processor.ProcessBlocks(2, []string{"b"})
	assert.EqualValues(t, 2, height)

	// but if we try again with other blocks, such that any valid height would have been previously recorded,
	// the blocks is not be accepted
	processor.ProcessBlocks(1, []string{"c", "d"})
	processor.ProcessBlocks(1, []string{"c", "d"})
	height = processor.ProcessBlocks(1, []string{"c", "d"})
	assert.EqualValues(t, 2, height)
}

func TestAcceptBlockWithDifferentStartHeightInDifferentCalls(t *testing.T) {
	// accept 'a' block first while setting up 'b' block to be accepted later
	var height uint64 = 0
	processor := NewBlockProcessor()

	// accepted 'a' block with height 1 and setup block 'b' to be accepted later
	processor.ProcessBlocks(1, []string{"a", "b"})
	processor.ProcessBlocks(1, []string{"a", "b"})
	height = processor.ProcessBlocks(1, []string{"a", "-", "b"})

	height = processor.ProcessBlocks(2, []string{"b"})
	assert.EqualValues(t, 2, height)
}

// concurrency tests
func TestConcurrentSingleBlockAccepted(t *testing.T) {
	const (
		concurrencyLevel        = 12
		expectedHeight   uint64 = 1
	)
	processor := NewBlockProcessor()

	wg := sync.WaitGroup{}
	wg.Add(concurrencyLevel)
	blocks := []string{"single"}

	for i := 0; i < concurrencyLevel; i++ {
		go func() {
			processor.ProcessBlocks(1, blocks)
			wg.Done()
		}()
	}

	wg.Wait()

	// same trick to get max height
	height := processor.ProcessBlocks(1, []string{})
	assert.Equal(t, expectedHeight, height)
}

func TestConcurrentBlocksAcceptedAsHeightIncreasesWhileProcessing(t *testing.T) {
	const (
		concurrencyLevel        = 4
		expectedHeight   uint64 = 3
	)
	processor := NewBlockProcessor()

	wg := sync.WaitGroup{}
	wg.Add(concurrencyLevel)
	blocks := []string{"a", "b", "c"}

	// accept 2 blocks out of 3 concurrently
	for i := 0; i < concurrencyLevel; i++ {
		go func() {
			processor.ProcessBlocks(1, blocks)
			wg.Done()
		}()
	}

	wg.Wait()
	// after processing is done, 2 blocks should have been accepted at height 1 and 2
	// same trick to get max height
	height := processor.ProcessBlocks(1, []string{})

	assert.EqualValues(t, expectedHeight, height)
}

// Test helpers

// processTestBlocks makes makes as many ProcessBlocks calls as the length of the blocksList
// which is a list of blocks (slice of slice of blocks)
// All calls have the same start height
func processTestBlocks(t *testing.T, blocksList ...testBlockList) uint64 {
	var height uint64 = 0
	processor := NewBlockProcessor()

	for _, blocks := range blocksList {
		height = processor.ProcessBlocks(1, blocks)
	}

	return height
}
