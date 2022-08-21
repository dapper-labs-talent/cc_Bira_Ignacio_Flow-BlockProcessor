## Assumptions made for the solution

- Two blocks cannot be accepted for the same height
- Calls to `ProcessBlocks` have `startHeight > 0`
- There's no need to track previouslly accepted block and we can use max height for that, given the above
- This also mean we cannot accept blocks with height bellow
- Max height can change inside ProcessBlocks (due to concurrency) preventing other blocks from being accepted
- Max height is always returned even if no block is accepted
- As per clarification email, blocks IDs are unique transaction IDs
- Blocks with empty strings for IDs are not valid

>  2b. Yes, in fact, we recently realized that the description is confusing here. It's more accurate to think of the IDs as transaction IDs for a given height.