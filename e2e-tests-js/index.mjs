import { StargateClient } from "@cosmjs/stargate";
import { spawn } from "child_process";
import { fileURLToPath } from "url";
import path from "path";
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

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
      console.log("child process exited with code ", code);
      if (code === 0) resolve(code);
      reject(code);
    });
  });
}

async function startMultinodeLocal() {
  await executeSpawn(path.join(scriptPath, "multinode-local-testnet.sh"));
}

async function test() {
  // start the network with three validator nodes. Wait til the network is up
  await startMultinodeLocal();

  const client = await StargateClient.connect("http://localhost:26657");
  const result = await client.getChainId();
  console.log("chain id: ", result);

  // finalize the test and remove everything
  await cleanNetwork();
}

async function cleanNetwork() {
  await executeSpawn(path.join(scriptPath, "clean-multinode-local-testnet.sh"));
}

test();
