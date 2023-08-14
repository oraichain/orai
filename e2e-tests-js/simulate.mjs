import { SimulateCosmWasmClient } from "@oraichain/cw-simulate";
import { __dirname } from "./utils.mjs";
import path from "path";
const contractDir = path.join(__dirname, "..", "scripts", "wasm_file");
import { readFileSync } from "fs";
import { coins } from "@cosmjs/proto-signing";

const sender = "orai18hr8jggl3xnrutfujy2jwpeu0l76azprlvgrwt";

function getContractDir(name) {
  return path.join(contractDir, name + ".wasm");
}

export const deployContract = async (
  client,
  senderAddress,
  msg,
  contractName,
  label
) => {
  // upload and instantiate the contract
  const wasmBytecode = readFileSync(getContractDir(contractName));
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

function initSimulateCosmWasmClient() {
  return new SimulateCosmWasmClient({
    chainId: "Oraichain",
    bech32Prefix: "orai",
  });
}

async function deployMultiExecuteContract(client) {
  const multiExecute = await deployContract(
    client,
    sender,
    {},
    "multi-execute",
    "multi-execute"
  );
  return { multiExecute };
}

// simulate multi-execute
async function simulate() {
  // fixture
  const cosmwasmClient = initSimulateCosmWasmClient();
  const { multiExecute } = await deployMultiExecuteContract(cosmwasmClient);
  const multiExecuteAddress = multiExecute.contractAddress;
  cosmwasmClient.app.bank.setBalance(multiExecuteAddress, coins(3, "orai"));
  const msgs = [
    {
      bank: {
        send: {
          to_address: "receiver",
          amount: [
            {
              denom: "orai",
              amount: "1",
            },
          ],
        },
      },
    },
    {
      bank: {
        send: {
          to_address: "receiver2",
          amount: [
            {
              denom: "orai",
              amount: "2",
            },
          ],
        },
      },
    },
  ];
  // action
  const result = await cosmwasmClient.execute(sender, multiExecuteAddress, {
    execute: {
      msg: Buffer.from(JSON.stringify(msgs)).toString("base64"),
    },
  });

  // assert
  if (result.events.filter((ev) => ev.type === "transfer").length !== 2) {
    throw "Error simulating the multi-execute contract";
  }
  console.log("Simulation passed!");
}

simulate();
