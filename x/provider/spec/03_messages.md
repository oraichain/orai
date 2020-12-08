<!--
order: 3
-->

# Messages

All messages relating to the provider module have been specified carefully in the official [document](https://docs.orai.io/docs/WhitePaper/ProtocolMessages) website of Oraichain. Below are the provider module messages that point to the above website:

## MsgCreateOracleScript

[MsgCreateOracleScript](https://docs.orai.io/docs/WhitePaper/ProtocolMessages#msgcreateoraclescript)

If one of the below conditions occurs, the message will not be accepted by the network:

- Oracle Script name is empty or the name has been taken.
- Oracle script source code is empty.
- AI Data Sources and Test Cases cannot be found.
- Fees provided by the user is lower than the required fees of the Oracle Script.
- The fee type is invalid.

## MsgEditOracleScript

[MsgEditOracleScript](https://docs.orai.io/docs/WhitePaper/ProtocolMessages#msgeditoraclescript)

If one of the below conditions occurs, the message will not be accepted by the network:

- Cannot find the Oracle Script given its unique identifier.
- Oracle Script old name or source code or new name is empty.
- Wrong owner of the Oracle Script.
- AI Data Sources and Test Cases of the new Oracle Script source code cannot be found.
- Fees provided by the user is lower than the required fees of the Oracle Script.
- The fee type is invalid.

## MsgCreateAIDataSource

[MsgCreateAIDataSource](https://docs.orai.io/docs/WhitePaper/ProtocolMessages#msgcreateaidatasource)

If one of the below conditions occurs, the message will not be accepted by the network:

- AI Data Source name is empty or the name has been taken.
- AI Data Source source code is empty.
- The fee type is invalid.

## MsgEditAIDataSource

[MsgEditAIDataSource](https://docs.orai.io/docs/WhitePaper/ProtocolMessages#msgeditaidatasource)

- Cannot find the AI Data Source given its unique identifier.
- AI Data Source old name or source code or new name is empty.
- Wrong owner of the AI Data Source.
- The fee type is invalid.

## MsgCreateTestCase

[MsgCreateTestCase](https://docs.orai.io/docs/WhitePaper/ProtocolMessages#msgcreatetestcase)

- Test Case name is empty or the name has been taken.
- Test Case source code is empty.
- The fee type is invalid.

## MsgEditTestCase

[MsgEditTestCase](https://docs.orai.io/docs/WhitePaper/ProtocolMessages#msgedittestcase)

- Cannot find the Test Case given its unique identifier.
- Test Case old name or source code or new name is empty.
- Wrong owner of the Test Case.
- The fee type is invalid.