import type { Address } from "viem";

export interface RelayRequest {
  ephemeralPublicKey: string;
  burnAddress: string;
}

export interface RelayResponse {
  success: boolean;
  message: string;
  txHash?: string;
  error?: string;
}

export async function callRelay(
  ephemeralPublicKey: string,
  burnAddress: Address,
): Promise<RelayResponse> {
  const relayUrl = import.meta.env.VITE_RELAY_URL;

  if (!relayUrl) {
    throw new Error("VITE_RELAY_URL environment variable is not set");
  }

  try {
    const response = await fetch(`${relayUrl}/relay`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        ephemeralPublicKey,
        burnAddress,
      }),
    });

    if (!response.ok) {
      throw new Error(
        `Relay request failed with status ${response.status}`,
      );
    }

    const data: RelayResponse = await response.json();
    return data;
  } catch (error) {
    const errorMessage =
      error instanceof Error ? error.message : "Unknown relay error";
    throw new Error(`Failed to call relay: ${errorMessage}`);
  }
}
