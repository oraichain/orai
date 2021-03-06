use crate::error::ContractError;
use crate::msg::{DataSourceQueryMsg, HandleMsg, InitMsg, QueryMsg, SpecialQuery};
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
        } => to_binary(&classification_testcase(deps, &contract, input, output)?),
    }
}

fn classification_testcase<S: Storage, A: Api, Q: Querier>(
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
    // basic validation for the data source result
    // create specialquery with to handle the data source
    let req = SpecialQuery::Fetch {
        url: "http://192.168.1.47:5001/testcase_classification".to_string(),
        body: data_source.to_string(),
        method: "POST".to_string(),
    }
    .into();
    let response: Binary = deps.querier.custom_query(&req)?;
    let response_str = String::from_utf8(response.to_vec()).unwrap();
    Ok(response_str)
}
