package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testBlockList []string

func TestProcessorStartsWithacceptedGenesisBlock(t *testing.T) {
	processor := NewBlockProcessor()

	assert.EqualValues(t, processor.MaxAcceptedHeight(), 0)
}

func TestDifferentBlocksNotaccepted(t *testing.T) {
	height := processTestBlocks(t, []string{"a", "b", "c"})
	assert.EqualValues(t, height, 0)
}

func TestSameBlockDifferentHeightsNotAccepted(t *testing.T) {
	height := processTestBlocks(t, []string{"a", "a", "a"})
	assert.EqualValues(t, height, 0)
}

func TestFirstAcceptedBlock(t *testing.T) {
	// the first block height after the genesis block is accepted
	// the next accepted block is at height 1
	var expectedHeight uint64 = 1
	height := processTestBlocks(t,
		[]string{"accepted", "not accepted"},
		[]string{"accepted", "also not accepted"},
		[]string{"accepted", "inaccepted again"},
	)

	assert.Equal(t, height, expectedHeight)
}

func TestThreeConfirmationsNeededForBlockToBeaccepted(t *testing.T) {
	var expectedHeight uint64 = 0
	height := processTestBlocks(t,
		[]string{"accepted", "not accepted"},
		[]string{"accepted", "also not accepted"},
	)

	assert.Equal(t, height, expectedHeight)
}

// test helpers
func processTestBlocks(t *testing.T, blocksList ...testBlockList) uint64 {
	var height uint64 = 0
	processor := NewBlockProcessor()

	for _, blocks := range blocksList {
		height = processor.ProcessBlocks(0, blocks)
	}

	return height
}
