<!--
order: 4
-->

# Begin-Block

When the system is at abci begin block, a new seed for random validator algorithm is generated to make sure that no validator can manipulate its value.

In the beginning, a completely 64-byte new seed is generated. After each begin block, one byte is removed, and a new one from the block hash is appended to the byte array. By doing this, only the block proposer, which is chosen randomly by the consensus algorithm, will be able to know the new byte being added. 
