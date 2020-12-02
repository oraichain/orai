use cosmwasm_std::{
    to_binary, Api, Binary, Env, Extern, HandleResponse, HandleResult, HumanAddr, InitResponse,
    InitResult, Querier, StdError, StdResult, Storage,
};

use crate::coin_helpers::assert_sent_sufficient_coin;
use crate::msg::{HandleMsg, InitMsg, QueryMsg, ResolveRecordResponse};
use crate::state::{config, config_read, resolver, resolver_read, Config, NameRecord};

const MIN_NAME_LENGTH: usize = 3;
const MAX_NAME_LENGTH: usize = 64;

pub fn init<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    _env: Env,
    msg: InitMsg,
) -> InitResult {
    let config_state = Config {
        purchase_price: msg.purchase_price,
        transfer_price: msg.transfer_price,
    };

    config(&mut deps.storage).save(&config_state)?;

    Ok(InitResponse::default())
}

pub fn handle<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    msg: HandleMsg,
) -> HandleResult {
    match msg {
        HandleMsg::Register { name } => try_register(deps, env, name),
        HandleMsg::Transfer { name, to } => try_transfer(deps, env, name, to),
    }
}

pub fn try_register<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    name: String,
) -> HandleResult {
    // we only need to check here - at point of registration
    validate_name(&name)?;
    let config_state = config(&mut deps.storage).load()?;
    assert_sent_sufficient_coin(&env.message.sent_funds, config_state.purchase_price)?;

    let key = name.as_bytes();
    let record = NameRecord {
        owner: deps.api.canonical_address(&env.message.sender)?,
    };

    if (resolver(&mut deps.storage).may_load(key)?).is_some() {
        // name is already taken
        return Err(StdError::generic_err("Name is already taken"));
    }

    // name is available
    resolver(&mut deps.storage).save(key, &record)?;

    Ok(HandleResponse::default())
}

pub fn try_transfer<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    name: String,
    to: HumanAddr,
) -> HandleResult {
    let api = deps.api;
    let config_state = config(&mut deps.storage).load()?;
    assert_sent_sufficient_coin(&env.message.sent_funds, config_state.transfer_price)?;

    let key = name.as_bytes();
    let new_owner = deps.api.canonical_address(&to)?;

    resolver(&mut deps.storage).update(key, |record| {
        if let Some(mut record) = record {
            if api.canonical_address(&env.message.sender)? != record.owner {
                return Err(StdError::unauthorized());
            }

            record.owner = new_owner.clone();
            Ok(record)
        } else {
            Err(StdError::generic_err("Name does not exist"))
        }
    })?;
    Ok(HandleResponse::default())
}

pub fn query<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    msg: QueryMsg,
) -> StdResult<Binary> {
    match msg {
        QueryMsg::ResolveRecord { name } => query_resolver(deps, name),
        QueryMsg::Config {} => to_binary(&config_read(&deps.storage).load()?),
    }
}

fn query_resolver<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    name: String,
) -> StdResult<Binary> {
    let key = name.as_bytes();

    let address = match resolver_read(&deps.storage).may_load(key)? {
        Some(record) => Some(deps.api.human_address(&record.owner)?),
        None => None,
    };
    let resp = ResolveRecordResponse { address };

    to_binary(&resp)
}

// let's not import a regexp library and just do these checks by hand
fn invalid_char(c: char) -> bool {
    let is_valid =
        (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c == '.' || c == '-' || c == '_');
    !is_valid
}

/// validate_name returns an error if the name is invalid
/// (we require 3-64 lowercase ascii letters, numbers, or . - _)
fn validate_name(name: &str) -> StdResult<()> {
    if name.len() < MIN_NAME_LENGTH {
        Err(StdError::generic_err("Name too short"))
    } else if name.len() > MAX_NAME_LENGTH {
        Err(StdError::generic_err("Name too long"))
    } else {
        match name.find(invalid_char) {
            None => Ok(()),
            Some(bytepos_invalid_char_start) => {
                let c = name[bytepos_invalid_char_start..].chars().next().unwrap();
                Err(StdError::generic_err(format!("Invalid character: '{}'", c)))
            }
        }
    }
}
