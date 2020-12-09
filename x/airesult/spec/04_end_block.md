<!--
order: 4
-->

# End-Block

## Reward

When the system is at abci end block, it retrieves a list of reports broadcast in that block to prepare for the next block reward round. If it cannot find any report, it skips and moves on to the next block. Else, it verifies those reports and stores all valid information in a `reward` object state.

## AIRequestResult

In addition to the `reward`, the request's result is also resolved depending on the requirement of the request. Indeed, if the result does not have enough reports from the validator, its status will be `pending`, whereas the status turns `expired` after a fixed number of blocks have passed. The status is only `success` when the correct number of reports have been collected within the number of allowed blocks.