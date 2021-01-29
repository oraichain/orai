use crate::error::ContractError;
use crate::msg::{HandleMsg, InitMsg, QueryMsg, SpecialQuery};
use cosmwasm_std::{
    to_binary, Api, Binary, Env, Extern, HandleResponse, InitResponse, MessageInfo, Querier,
    StdResult, Storage,
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
        QueryMsg::GetPrice {} => to_binary(&query_price(deps)?),
    }
}

fn query_price<S: Storage, A: Api, Q: Querier>(deps: &Extern<S, A, Q>) -> StdResult<String> {
    // create specialquery with default empty string
    let req = SpecialQuery::Fetch {
        url: "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
            .to_string(),
    }
    .into();
    // because not support f32, we need to do it manually
    // dont use String because it will deserialize bytes to base 64
    let response: Binary = deps.querier.custom_query(&req)?;
    let data = String::from_utf8(response.to_vec()).unwrap();
    let first = data.find(r#""usd":"#).unwrap() + 6;
    let last = first + data.get(first..).unwrap().find("}").unwrap();
    let price = data.get(first..last).unwrap().to_string();
    Ok(price)
}
