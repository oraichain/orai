import { fileURLToPath } from "url";
import path from "path";
import { readFileSync } from "fs";
export const __filename = fileURLToPath(import.meta.url);
export const __dirname = path.dirname(__filename);

function getContractDir(contractDir, name) {
  return path.join(contractDir, name + ".wasm");
}

export const deployContract = async (
  client,
  senderAddress,
  msg,
  contractName,
  label,
  contractDir = undefined
) => {
  // upload and instantiate the contract
  const wasmBytecode = readFileSync(
    getContractDir(
      contractDir ?? path.join(__dirname, "..", "scripts", "wasm_file"),
      contractName
    )
  );
  const uploadRes = await client.upload(senderAddress, wasmBytecode, "auto");
  const initRes = await client.instantiate(
    senderAddress,
    uploadRes.codeId,
    msg ?? {},
    label ?? contractName,
    "auto"
  );
  return { ...uploadRes, ...initRes };
};

export async function deployMultiExecuteContract(client, sender) {
  const multiExecute = await deployContract(
    client,
    sender,
    {},
    "multi-execute",
    "multi-execute"
  );
  return { multiExecute };
}
