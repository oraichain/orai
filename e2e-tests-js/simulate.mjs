import { SimulateCosmWasmClient } from "@oraichain/cw-simulate";
import { __dirname, deployMultiExecuteContract } from "./utils.mjs";
import { coins } from "@cosmjs/proto-signing";

const sender = "orai18hr8jggl3xnrutfujy2jwpeu0l76azprlvgrwt";

function initSimulateCosmWasmClient() {
  return new SimulateCosmWasmClient({
    chainId: "Oraichain",
    bech32Prefix: "orai",
  });
}

// simulate multi-execute
async function simulate() {
  // fixture
  const cosmwasmClient = initSimulateCosmWasmClient();
  const { multiExecute } = await deployMultiExecuteContract(
    cosmwasmClient,
    sender
  );
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
