pub mod contract;
pub mod error;
pub mod msg;

#[cfg(target_arch = "wasm32")]
cosmwasm_std::create_entry_points!(contract);
