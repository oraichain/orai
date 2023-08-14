import { GasPrice } from "@cosmjs/stargate";
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";
import { spawn } from "child_process";
import { __dirname, deployMultiExecuteContract } from "./utils.mjs";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import "dotenv/config";
import path from "path";

const scriptPath = path.join(__dirname, "..", "scripts");

async function executeSpawn(command, params) {
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
