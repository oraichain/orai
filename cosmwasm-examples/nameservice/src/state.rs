use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{CanonicalAddr, Coin, ReadonlyStorage, Storage};
use cosmwasm_storage::{
    bucket, bucket_read, singleton, singleton_read, Bucket, ReadonlyBucket, ReadonlySingleton,
    Singleton,
};

pub static NAME_RESOLVER_KEY: &[u8] = b"nameresolver";
pub static CONFIG_KEY: &[u8] = b"config";

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct Config {
    pub purchase_price: Option<Coin>,
    pub transfer_price: Option<Coin>,
}

pub fn config<S: Storage>(storage: &mut S) -> Singleton<S, Config> {
    singleton(storage, CONFIG_KEY)
}

pub fn config_read<S: ReadonlyStorage>(storage: &S) -> ReadonlySingleton<S, Config> {
    singleton_read(storage, CONFIG_KEY)
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct NameRecord {
    pub owner: CanonicalAddr,
}

pub fn resolver<S: Storage>(storage: &mut S) -> Bucket<S, NameRecord> {
    bucket(NAME_RESOLVER_KEY, storage)
}

pub fn resolver_read<S: Storage>(storage: &S) -> ReadonlyBucket<S, NameRecord> {
    bucket_read(NAME_RESOLVER_KEY, storage)
}
