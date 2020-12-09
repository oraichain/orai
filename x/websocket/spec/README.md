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
2. **[Messages](02_messages.md)**
    - [MsgCreateReport](02_messages.md#MsgCreateReport)
    - [MsgCreateReporter](02_messages.md#MsgAddReporter)
    - [MsgRemoveReporter](02_messages.md#MsgRemoveReporter)
