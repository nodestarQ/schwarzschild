import {
  createPublicClient,
  http,
  type Address,
  defineChain,
  parseAbiItem,
} from "viem";

export interface StealthBurnLog {
  ephemeralPublicKey: string;
  burnAddress: Address;
  blockNumber: bigint;
  transactionHash: string;
  logIndex: number;
}

const sepolia = /*#__PURE__*/ defineChain({
  id: 11_155_111,
  name: "Sepolia",
  nativeCurrency: { name: "Sepolia Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: {
    default: {
      http: ["https://ethereum-sepolia-rpc.publicnode.com"],
    },
  },
  blockExplorers: {
    default: {
      name: "Etherscan",
      url: "https://sepolia.etherscan.io",
      apiUrl: "https://api-sepolia.etherscan.io/api",
    },
  },
  testnet: true,
});

// StealthBurn event ABI - parameters are NOT indexed in the actual contract
const STEALTH_BURN_EVENT = parseAbiItem(
  "event StealthBurn(bytes32 ephemeralPublicKey, address burnAddress)",
);

function getPublicClient() {
  return createPublicClient({
    chain: sepolia,
    transport: http(),
  });
}

export async function fetchStealthBurnLogs(
  contractAddress: Address,
  limit: number = 100,
): Promise<StealthBurnLog[]> {
  try {
    const publicClient = getPublicClient();

    // Get the latest block
    const latestBlock = await publicClient.getBlockNumber();

    // Fetch logs for the StealthBurn event
    const logs = await publicClient.getLogs({
      address: contractAddress,
      event: STEALTH_BURN_EVENT,
      fromBlock:
        latestBlock >= BigInt(1000) ? latestBlock - BigInt(1000) : BigInt(0),
      toBlock: latestBlock,
    });

    console.log(`Found ${logs.length} StealthBurn events`, logs);

    // Parse the logs - the logs are already decoded by viem
    const parsedLogs: StealthBurnLog[] = logs.map((log) => {
      // Handle both cases: indexed and non-indexed parameters
      const ephemeralKey = log.args.ephemeralPublicKey;
      const burnAddr = log.args.burnAddress;

      return {
        ephemeralPublicKey: ephemeralKey
          ? typeof ephemeralKey === "string"
            ? ephemeralKey
            : ephemeralKey.toString()
          : "0x",
        burnAddress:
          burnAddr && typeof burnAddr === "string"
            ? (burnAddr as Address)
            : ("0x" as Address),
        blockNumber: log.blockNumber,
        transactionHash: log.transactionHash,
        logIndex: log.logIndex,
      };
    });

    // Sort by block number descending (most recent first)
    return parsedLogs.sort((a, b) => Number(b.blockNumber - a.blockNumber));
  } catch (error) {
    console.error("Failed to fetch contract logs:", error);
    return [];
  }
}
