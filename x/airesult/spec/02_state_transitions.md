<!--
order: 2
-->

# State Transitions

This document describes the state transitions of [AIRequestResult](./01_state.md#airequestresult)

## AIRequestResult

The state changes of this state are performed when there is a new report relating to the corresponding AI request.

### Pending to expired / new to expired

A request is `expired` when a report comes after the number of blocks allowed. When this happens:

- The request's result changes from `pending` to `expired` if the result object already exists.
- A new result is created with the `expired` state if no result is found yet on Oraichain.

### Pending to finished / new to finished

A request status is finished when a report comes within the number of blocks allowed. When this happens:

- The request's result changes from `pending` to `finished` if the result object already exists.
- A new result is created with the `finished` state if the request only needs one report.
