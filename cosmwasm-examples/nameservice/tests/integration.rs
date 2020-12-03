//! This integration test tries to run and call the generated wasm.
//! It depends on a Wasm build being available, which you can create with `cargo wasm`.
//! Then running `cargo integration-test` will validate we can properly call into that generated Wasm.
//!
//! You can easily convert unit tests to integration tests as follows:
//! 1. Copy them over verbatim
//! 2. Then change
//!      let mut deps = mock_dependencies(20, &[]);
//!    to
//!      let mut deps = mock_instance(WASM, &[]);
//! 3. If you access raw storage, where ever you see something like:
//!      deps.storage.get(CONFIG_KEY).expect("no data stored");
//!    replace it with:
//!      deps.with_storage(|store| {
//!          let data = store.get(CONFIG_KEY).expect("no data stored");
//!          //...
//!      });
//! 4. Anywhere you see query(&deps, ...) you must replace it with query(&mut deps, ...)

use cosmwasm_std::{
    coin, coins, from_binary, Coin, HandleResponse, HandleResult, HumanAddr, InitResponse, StdError,
};
use cosmwasm_storage::to_length_prefixed;
use cosmwasm_vm::testing::{
    handle, init, mock_env, mock_instance, query, MockApi, MockQuerier, MockStorage,
};
use cosmwasm_vm::{from_slice, Instance, Storage};

use cw_nameservice::msg::{HandleMsg, InitMsg, QueryMsg, ResolveRecordResponse};
use cw_nameservice::state::{Config, CONFIG_KEY};

// This line will test the output of cargo wasm
static WASM: &[u8] = include_bytes!("../target/wasm32-unknown-unknown/release/cw_nameservice.wasm");
// You can uncomment this line instead to test productionified build from rust-optimizer
// static WASM: &[u8] = include_bytes!("../contract.wasm");

fn assert_name_owner(
    mut deps: &mut Instance<MockStorage, MockApi, MockQuerier>,
    name: &str,
    owner: &str,
) {
    let res = query(
        &mut deps,
        QueryMsg::ResolveRecord {
            name: name.to_string(),
        },
    )
    .unwrap();

    let value: ResolveRecordResponse = from_binary(&res).unwrap();
    assert_eq!(Some(HumanAddr::from(owner)), value.address);
}

fn mock_init_with_price(
    mut deps: &mut Instance<MockStorage, MockApi, MockQuerier>,
    purchase_price: Coin,
    transfer_price: Coin,
) {
    let msg = InitMsg {
        purchase_price: Some(purchase_price),
        transfer_price: Some(transfer_price),
    };

    let params = mock_env("creator", &coins(2, "token"));
    // unwrap: contract successfully handles InitMsg
    let _res: InitResponse = init(&mut deps, params, msg).unwrap();
}

fn mock_init_no_price(mut deps: &mut Instance<MockStorage, MockApi, MockQuerier>) {
    let msg = InitMsg {
        purchase_price: None,
        transfer_price: None,
    };

    let params = mock_env("creator", &coins(2, "token"));
    // unwrap: contract successfully handles InitMsg
    let _res: InitResponse = init(&mut deps, params, msg).unwrap();
}

fn mock_alice_registers_name(
    mut deps: &mut Instance<MockStorage, MockApi, MockQuerier>,
    sent: &[Coin],
) {
    // alice can register an available name
    let params = mock_env("alice_key", sent);
    let msg = HandleMsg::Register {
        name: "alice".to_string(),
    };
    // unwrap: contract successfully handles Register message
    let _res: HandleResponse = handle(&mut deps, params, msg).unwrap();
}

#[test]
fn proper_init_no_fees() {
    let mut deps = mock_instance(WASM, &[]);

    mock_init_no_price(&mut deps);

    deps.with_storage(|storage| {
        let key = to_length_prefixed(CONFIG_KEY);
        let data = storage.get(&key).0.unwrap().unwrap();
        let config_state: Config = from_slice(&data).unwrap();

        assert_eq!(
            config_state,
            Config {
                purchase_price: None,
                transfer_price: None
            }
        );
        Ok(())
    })
    .unwrap();
}

#[test]
fn proper_init_with_prices() {
    let mut deps = mock_instance(WASM, &[]);

    mock_init_with_price(&mut deps, coin(3, "token"), coin(4, "token"));

    deps.with_storage(|storage| {
        let key = to_length_prefixed(CONFIG_KEY);
        let data = storage.get(&key).0.unwrap().unwrap();
        let config_state: Config = from_slice(&data).unwrap();

        assert_eq!(
            config_state,
            Config {
                purchase_price: Some(coin(3, "token")),
                transfer_price: Some(coin(4, "token")),
            }
        );

        Ok(())
    })
    .unwrap();
}

#[test]
fn register_available_name_and_query_works_with_prices() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_with_price(&mut deps, coin(2, "token"), coin(2, "token"));
    mock_alice_registers_name(&mut deps, &coins(2, "token"));

    // anyone can register an available name with more fees than needed
    let params = mock_env("bob_key", &coins(5, "token"));
    let msg = HandleMsg::Register {
        name: "bob".to_string(),
    };

    // unwrap: contract successfully handles Register message
    let _res: HandleResponse = handle(&mut deps, params, msg).unwrap();

    // querying for name resolves to correct address
    assert_name_owner(&mut deps, "alice", "alice_key");
    assert_name_owner(&mut deps, "bob", "bob_key");
}

#[test]
fn register_available_name_and_query_works() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_no_price(&mut deps);
    mock_alice_registers_name(&mut deps, &[]);

    // querying for name resolves to correct address
    assert_name_owner(&mut deps, "alice", "alice_key");
}

#[test]
fn fails_on_register_already_taken_name() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_no_price(&mut deps);
    mock_alice_registers_name(&mut deps, &[]);

    // bob can't register the same name
    let params = mock_env("bob_key", &coins(2, "token"));
    let msg = HandleMsg::Register {
        name: "alice".to_string(),
    };
    let res: HandleResult = handle(&mut deps, params, msg);
    match res.unwrap_err() {
        StdError::GenericErr { msg, .. } => assert_eq!(msg, "Name is already taken"),
        _ => panic!("Unexpected error type"),
    }

    // alice can't register the same name again
    let params = mock_env("alice_key", &coins(2, "token"));
    let msg = HandleMsg::Register {
        name: "alice".to_string(),
    };
    let res: HandleResult = handle(&mut deps, params, msg);
    match res.unwrap_err() {
        StdError::GenericErr { msg, .. } => assert_eq!(msg, "Name is already taken"),
        _ => panic!("Unexpected error type"),
    }
}

#[test]
fn fails_on_register_insufficient_fees() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_with_price(&mut deps, coin(2, "token"), coin(2, "token"));

    // anyone can register an available name with sufficient fees
    let params = mock_env("alice_key", &[]);
    let msg = HandleMsg::Register {
        name: "alice".to_string(),
    };

    let res: HandleResult = handle(&mut deps, params, msg);
    match res.unwrap_err() {
        StdError::GenericErr { msg, .. } => assert_eq!(msg, "Insufficient funds sent"),
        _ => panic!("Unexpected error type"),
    }
}

#[test]
fn fails_on_register_wrong_fee_denom() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_with_price(&mut deps, coin(2, "token"), coin(2, "token"));

    // anyone can register an available name with sufficient fees
    let params = mock_env("alice_key", &coins(2, "earth"));
    let msg = HandleMsg::Register {
        name: "alice".to_string(),
    };

    let res: HandleResult = handle(&mut deps, params, msg);
    match res.unwrap_err() {
        StdError::GenericErr { msg, .. } => assert_eq!(msg, "Insufficient funds sent"),
        _ => panic!("Unexpected error type"),
    }
}

#[test]
fn transfer_works() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_no_price(&mut deps);
    mock_alice_registers_name(&mut deps, &[]);

    // alice can transfer her name successfully to bob
    let params = mock_env("alice_key", &[]);
    let msg = HandleMsg::Transfer {
        name: "alice".to_string(),
        to: HumanAddr::from("bob_key"),
    };

    let _res: HandleResponse = handle(&mut deps, params, msg).unwrap();
    // querying for name resolves to correct address (bob_key)
    assert_name_owner(&mut deps, "alice", "bob_key");
}

#[test]
fn transfer_works_with_fees() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_with_price(&mut deps, coin(2, "token"), coin(2, "token"));
    mock_alice_registers_name(&mut deps, &coins(2, "token"));

    // alice can transfer her name successfully to bob
    let params = mock_env("alice_key", &vec![coin(1, "earth"), coin(2, "token")]);
    let msg = HandleMsg::Transfer {
        name: "alice".to_string(),
        to: HumanAddr::from("bob_key"),
    };

    let _res: HandleResponse = handle(&mut deps, params, msg).unwrap();
    // querying for name resolves to correct address (bob_key)
    assert_name_owner(&mut deps, "alice", "bob_key");
}

#[test]
fn fails_on_transfer_from_nonowner() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_no_price(&mut deps);
    mock_alice_registers_name(&mut deps, &[]);

    // alice can transfer her name successfully to bob
    let params = mock_env("frank_key", &coins(2, "token"));
    let msg = HandleMsg::Transfer {
        name: "alice".to_string(),
        to: HumanAddr::from("bob_key"),
    };

    let res: HandleResult = handle(&mut deps, params, msg);
    match res.unwrap_err() {
        StdError::Unauthorized { .. } => {}
        _ => panic!("Unexpected error type"),
    }

    // querying for name resolves to correct address (alice_key)
    assert_name_owner(&mut deps, "alice", "alice_key");
}

#[test]
fn fails_on_transfer_insufficient_fees() {
    let mut deps = mock_instance(WASM, &[]);
    mock_init_with_price(&mut deps, coin(2, "token"), coin(5, "token"));
    mock_alice_registers_name(&mut deps, &coins(2, "token"));

    // alice can transfer her name successfully to bob
    let params = mock_env("alice_key", &vec![coin(1, "earth"), coin(2, "token")]);
    let msg = HandleMsg::Transfer {
        name: "alice".to_string(),
        to: HumanAddr::from("bob_key"),
    };

    let res: HandleResult = handle(&mut deps, params, msg);
    match res.unwrap_err() {
        StdError::GenericErr { msg, .. } => assert_eq!(msg, "Insufficient funds sent"),
        _ => panic!("Unexpected error type"),
    }

    // querying for name resolves to correct address (bob_key)
    assert_name_owner(&mut deps, "alice", "alice_key");
}

#[test]
fn returns_empty_on_query_unregistered_name() {
    let mut deps = mock_instance(WASM, &[]);

    mock_init_no_price(&mut deps);

    // querying for unregistered name results in NotFound error
    let res = query(
        &mut deps,
        QueryMsg::ResolveRecord {
            name: "alice".to_string(),
        },
    )
    .unwrap();
    let value: ResolveRecordResponse = from_binary(&res).unwrap();
    assert_eq!(None, value.address);
}
