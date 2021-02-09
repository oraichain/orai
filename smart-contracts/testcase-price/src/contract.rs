use crate::error::ContractError;
use crate::msg::{DataSourceQueryMsg, HandleMsg, InitMsg, QueryMsg};
use cosmwasm_std::{
    to_binary, Api, Binary, Env, Extern, HandleResponse, HumanAddr, InitResponse, MessageInfo,
    Querier, StdResult, Storage,
};

// Note, you can use StdResult in some functions where you do not
// make use of the custom errors
pub fn init<S: Storage, A: Api, Q: Querier>(
    _deps: &mut Extern<S, A, Q>,
    _env: Env,
    _info: MessageInfo,
    _: InitMsg,
) -> StdResult<InitResponse> {
    Ok(InitResponse::default())
}

// And declare a custom Error variant for the ones where you will want to make use of it
pub fn handle<S: Storage, A: Api, Q: Querier>(
    _: &mut Extern<S, A, Q>,
    _env: Env,
    _: MessageInfo,
    _: HandleMsg,
) -> Result<HandleResponse, ContractError> {
    Ok(HandleResponse::default())
}

pub fn query<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    _env: Env,
    msg: QueryMsg,
) -> StdResult<Binary> {
    match msg {
        QueryMsg::Test {
            input,
            output,
            contract,
        } => to_binary(&test_price(deps, &contract, input, output)?),
    }
}

// using string pointer
fn parse_i32(input: &str) -> i32 {
    // get first item from iterator
    let number = input.split('.').next();
    // will panic instead for forward error with ?
    number.unwrap().parse().unwrap()
}

fn test_price<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    contract: &HumanAddr,
    input: String,
    output: String,
) -> StdResult<String> {
    // check output if empty then do nothing
    if output.is_empty() {
        return Ok(String::new());
    }
    let msg = DataSourceQueryMsg::Get { input };
    let data_source: String = deps.querier.query_wasm_smart(contract, &msg)?;
    // positive using unwrap, otherwise rather panic than return default value
    let data_source_result = parse_i32(&data_source);
    let expected_result = parse_i32(&output);
    let deviation = data_source_result - expected_result;
    if deviation < 10000 {
        Ok(data_source)
    } else {
        Ok(String::new())
    }
}
