<!--
order: 0
title: AI Request Overview
parent:
  title: "airequest"
-->

# `airequest`

## Abstract

This paper specifies the AI Request module of Oraichain.

The module allows token holders to create an AI request that targets a specific AI service provided on Oraichain. The request and its result will be stored on-chain as proof. Other systems also can use the result for their purposes.

## Contents

1. **[State](01_state.md)**
    - [AIRequest](01_state.md#airequest)
2. **[State Transitions](02_state_transitions.md)**
    - [Validators](02_state_transitions.md#validators)
    - [Delegations](02_state_transitions.md#delegations)
    - [Slashing](02_state_transitions.md#slashing)
3. **[Messages](03_messages.md)**
    - [MsgCreateValidator](03_messages.md#msgcreatevalidator)
    - [MsgEditValidator](03_messages.md#msgeditvalidator)
    - [MsgDelegate](03_messages.md#msgdelegate)
    - [MsgBeginUnbonding](03_messages.md#msgbeginunbonding)
    - [MsgBeginRedelegate](03_messages.md#msgbeginredelegate)
4. **[Begin-Block](04_begin_block.md)**
    - [Historical Info Tracking](04_begin_block.md#historical-info-tracking)
4. **[End-Block ](05_end_block.md)**
    - [Validator Set Changes](05_end_block.md#validator-set-changes)
    - [Queues ](05_end_block.md#queues-)
5. **[Hooks](06_hooks.md)**
6. **[Events](07_events.md)**
    - [EndBlocker](07_events.md#endblocker)
    - [Handlers](07_events.md#handlers)
7. **[Parameters](08_params.md)**
