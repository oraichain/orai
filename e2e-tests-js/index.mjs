import { GasPrice } from "@cosmjs/stargate";
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";
import {
  __dirname,
  deployMultiExecuteContract,
  executeSpawn,
} from "./utils.mjs";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import "dotenv/config";
import path from "path";

const scriptPath = path.join(__dirname, "..", "scripts");

async function initCosmWasmClient() {
  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
    process.env.MNEMONIC,
    { prefix: "orai" }
  );
  const cosmwasmClient = await SigningCosmWasmClient.connectWithSigner(
    "http://localhost:26657",
    wallet,
    { gasPrice: GasPrice.fromString("0.001orai") }
  );
  return { wallet, cosmwasmClient };
}

async function startMultinodeLocal() {
  await executeSpawn(path.join(scriptPath, "multinode-local-testnet.sh"));
}

async function testIbcHooks() {
  const { wallet, cosmwasmClient } = await initCosmWasmClient();
  const senderAddress = (await wallet.getAccounts())[0].address;
  console.log("preparing to deploy the multi execute contract...");
  const { multiExecute } = await deployMultiExecuteContract(
    cosmwasmClient,
    senderAddress
  );
  const msgs = [
    {
      bank: {
        send: {
          to_address: senderAddress,
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
          to_address: senderAddress,
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

  console.log("Running a batch of multi execute contract messages...");
  const result = await cosmwasmClient.execute(
    senderAddress,
    multiExecute.contractAddress,
    {
      execute: {
        msg: Buffer.from(JSON.stringify(msgs)).toString("base64"),
      },
    },
    "auto",
    undefined,
    [{ denom: "orai", amount: "3" }]
  );

  console.log("result: ", result);
}

async function cleanNetwork() {
  await executeSpawn(path.join(scriptPath, "clean-multinode-local-testnet.sh"));
}

async function test() {
  // start the network with three validator nodes. Wait til the network is up
  await startMultinodeLocal();

  // execute
  await testIbcHooks();

  // finalize the test and remove everything
  await cleanNetwork();
}

test();
