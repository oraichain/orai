<!--
order: 0
title: WebSocket Overview
parent:
  title: "websocket"
-->

# `websocket`

## Abstract

This paper specifies the WebSocket module of Oraichain.

After an AI Request is created, validators will execute the request before reporting back to Oraichain. The process of storing validator results and reports are handled by this module.

## Contents

1. **[State](01_state.md)**
2. **[Messages](03_messages.md)**
    - [MsgCreateReport](03_messages.md#MsgCreateReport)
    - [MsgCreateReporter](03_messages.md#MsgAddReporter)
    - [MsgRemoveReporter](03_messages.md#MsgRemoveReporter)
3. **[Events](07_events.md)**
4. **[Parameters](08_params.md)**
