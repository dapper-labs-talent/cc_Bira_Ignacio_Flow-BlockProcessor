package chain

import (
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

	assert.EqualValues(t, processor.MaxAcceptedHeight(), 0)
}

func TestEmptyBlocksKeepHeight(t *testing.T) {
	height := processTestBlocks(t, []string{})
	assert.EqualValues(t, height, 0)
}

func TestDifferentBlocksNotaccepted(t *testing.T) {
	height := processTestBlocks(t, []string{"a", "b", "c"})
	assert.EqualValues(t, height, 0)
}

func TestEmptyStringBlockIdsNotAccepted(t *testing.T) {
	height := processTestBlocks(t,
		[]string{"", "a"},
		[]string{""},
		[]string{"", "b", "c"},
	)
	assert.EqualValues(t, height, 0)
}

func TestSameBlockDifferentHeightsNotAccepted(t *testing.T) {
	height := processTestBlocks(t, []string{"a", "a", "a"})
	assert.EqualValues(t, height, 0)
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

func TestMultipleAcceptBlocksIncreaseHeight(t *testing.T) {
	const expectedHeight uint64 = 1

	// blocks with IDs 'a' and '1' are accepted after 4 calls to ProcessBlocks
	height := processTestBlocks(t,
		[]string{"1", "a", "c"},
		[]string{"1", "a", "f"},
		[]string{"1", "b", "strawberry", "apple"},
		[]string{"b", "a", "c"},
	)

	assert.Equal(t, expectedHeight, height)
}

// test helpers
// processTestBlocks makes makes as many ProcessBlocks calls as the length of the blocksList
// which is a list of blocks (slice of slice of blocks)
func processTestBlocks(t *testing.T, blocksList ...testBlockList) uint64 {
	var height uint64 = 0
	processor := NewBlockProcessor()

	for _, blocks := range blocksList {
		height = processor.ProcessBlocks(0, blocks)
	}

	return height
}
