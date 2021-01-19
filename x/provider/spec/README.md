<!--
order: 0
title: Provider Overview
parent:
  title: "provider"
-->

# `provider`

## Abstract

This paper specifies the Provider module of Oraichain.

This module represents a digital store that allows AI and Test Case providers to publish their work from AI models to quality test cases onto the network. It also has Oracle Script - a form of smart contract that connects AI and Test Case providers together to build a complete AI service for the token holders.

The holders can choose an AI service through the Oracle Script, and they will receive results from different providers as well as an aggregated version of all results. The Oracle Script will take care of how the results are aggregated, and the whole process is public.

## Contents

1. **[State](01_state.md)**
2. **[Messages](02_messages.md)**
    - [MsgCreateOracleScript](02_messages.md#msgcreateoraclescript)
    - [MsgEditOracleScript](02_messages.md#msgeditoraclescript)
    - [MsgCreateAIDataSource](02_messages.md#msgcreateaidatasource)
    - [MsgEditAIDataSource](02_messages.md#msgeditaidatasource)
    - [MsgCreateTestCase](02_messages.md#msgcreatetestcase)
    - [MsgEditTestCase](02_messages.md#msgedittestcase)
3. **[Events](03_events.md)**
4. **[Parameters](04_params.md)**
