# Block processing problem

In this coding challenge we leverage some basic blockchain concepts. Specifically we imitate how a blockchain is an
append-only linked list containing elements called blocks. Blocks are produced and accepted into the chain consecutively
such that the first block has height 0, its child has height 1, and so on. In real blockchains the mechanism for
accepting blocks into the chain is achieved through complex distributed consensus algorithms. However, in this challenge
we employ a simplified mechanism to accept a block inside the `BlockProcessor` component. We consider a given height to be
**accepted** only once we've received at least `3` responses agreeing on the same block at that height. In other words, we
expect 3 blocks of the same ID to be observed to accept it.

You may assume that the underlying representation of a `Block` is simply a string representing the block ID. It's also
guaranteed that blocks with different heights will have different IDs.

This coding challenge does not require any other detailed understanding or research into blockchain systems or design,
it is a discrete coding problem inspired by how blockchains work and which can be completed based solely on the
requirements set out below.

### Task:

Implement a `BlockProcessor` module that has a single method `ProcessBlocks` which takes a startHeight and a list of
`block`. It returns uint64 value as the max accepted height.

With the returned max accepted height, consider that there will be an external component calling from time to time for
the next range of blocks. When blocks are received by the upstream component `ProcessBlocks` will be called with the
blocks. Note, you don’t need to implement this external component.

`func (p *BlockProcessor) ProcessBlocks(startHeight uint64, blocks []string) uint64`

* The input `blocks` represents a consecutive range of blocks. The string value in the array is the block’s ID. The height of each block can be inferred from its index in the array and the `startHeight`. The first block in `blocks` has height `startHeight`, the second block in `blocks` has height `startHeight` + 1, and so on
* A block is accepted if at least `3` blocks of the same ID at a given height are received and its parent block (block at height - 1) is also accepted.
* Max accepted height is the height of the highest accepted block.
* You can assume that block IDs (the string) are unique in relation to the block height they are accepted for.
* The block at height `0` is the Genesis block, which is already accepted when `BlockProcessor` is created. In other words, max accepted height starts from `0`.
* Consider efficiency in your solution, call rates from node peers will be high and optimal time complexity is essential.

## Guidance

* We prefer candidates to use GoLang - our in house language, however, please feel free to use any language of your preference.
* You should ensure that your implementation is concurrency-safe.
* Please approach this as you would a real-world problem. We are not only assessing your ability to solve the problem but also trade-offs/edge cases considered and your holistic approach to quality.
* If anything remains unclear about this problem don’t hesitate to ask your Talent team associate who can get follow-ups from engineering. We will always respond to questions over email. However, we do not meet candidates to discuss questions over a call to ensure candidates are not given an unfair advantage.
* When submitting, please include a section in your response, or within your code, to summarize any assumptions, or other matters to share to us in considering your submission.