import { fileURLToPath } from "url";
import path from "path";
import { readFileSync } from "fs";
import { spawn } from "child_process";
export const __filename = fileURLToPath(import.meta.url);
export const __dirname = path.dirname(__filename);

export async function executeSpawn(command, params) {
  const proc = spawn(command, params);
  proc.stderr.pipe(process.stdout);

  return new Promise((resolve, reject) => {
    proc.stdout.on("data", function (data) {
      console.log(data.toString());
    });

    // proc.stderr.on("data", function (data) {
    //   //   console.log("stderr: " + data.toString());
    //   //   reject(data.toString());
    // });

    proc.on("exit", function (code) {
      // console.log("child process exited with code ", code);
      if (code === 0) resolve(code);
      reject(code);
    });
  });
}

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
