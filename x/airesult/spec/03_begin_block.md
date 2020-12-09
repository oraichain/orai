<!--
order: 3
-->

# Begin-Block

When the system is at abci begin block, an allocation token process is activated. Validators that report successfully will receive rewards when a new block comes. The system also allocates tokens to Data Source & Test Case owners included in the reports.

## Reasons behind a custom allocation token.

The fee collector, which is responsible for collecting fees from the Oraichain transactions, is a module account from the `supply` module. After each block, it gathers all the transaction fees and allocates them as an incentive to the bonded validators that participate in the consensus process. Nevertheless, this can lead to validators taking all the rewards if there is an AI request from the previous block as this process called by the `distribution` module is repetitive. Of course, we can stop using it and replace it with a custom function. However, it is built by the Cosmos-SDK team, which is quality enough not to re-build a new one.

## Allocating process

As a result, each time there is an AI request, the system removes the request's fees from the fee collector. Since the total transaction fees required for the AI request is already stored, we do not need to worry about losing information.

When a block contains one or more reports, the system will add the corresponding fees that have been removed from the fee collector.

Next, each Data Source & Test Case owner will receive his reward according to the fees specified in the Data Source & Test Case. Finally, the remaining transaction fees are allocated to the validators executing the request proportionally to their voting powers.

## BankKeeper bug

There is a bug when setting coins for the fee collector, where it does not update the bank keeper coins to send to other modules. Hence, we need to subtract coins from the bank keeper manually when removing the request fees. When we want to allocate tokens, we add the same amount that is removed to reward the entities.
