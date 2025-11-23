import {
	createWalletClient,
	custom,
	recoverMessageAddress,
	hashMessage,
	recoverPublicKey,
	type Address,
	type Hex,
} from 'viem';

/**
 * Sign a message using viem's walletClient.signMessage()
 * This follows the EIP-191 standard and uses the wallet provider
 * @param message - The plain text message to sign
 * @param account - The user's wallet address
 * @param provider - The EIP-1193 provider (window.ethereum)
 * @returns The signature as a hex string
 */
export async function signMessageWithWallet(
	message: string,
	account: Address,
	provider: any,
): Promise<Hex> {
	try {
		// Create a wallet client using the user's provider
		const walletClient = createWalletClient({
			account,
			transport: custom(provider),
		});

		// Sign the message using viem's signMessage
		// This automatically handles EIP-191 formatting and hashing
		const signature = await walletClient.signMessage({
			account,
			message,
		});

		console.log('Message signed successfully:', signature);
		return signature;
	} catch (error) {
		console.error('Failed to sign message:', error);
		throw error;
	}
}

/**
 * Recover the signer's address from a signed message
 * @param message - The original message that was signed
 * @param signature - The signature hex string
 * @returns The recovered address
 */
export async function recoverMessageSigner(
	message: string,
	signature: Hex,
): Promise<Address> {
	try {
		// Use viem's recoverMessageAddress to recover the signer address
		const address = await recoverMessageAddress({
			message,
			signature,
		});

		return address;
	} catch (error) {
		console.error('Failed to recover message signer:', error);
		throw error;
	}
}

/**
 * Recover the public key from a signed message
 * @param message - The original message that was signed
 * @param signature - The signature hex string
 * @returns The recovered public key as a hex string
 */
export async function recoverSignaturePublicKey(
	message: string,
	signature: Hex,
): Promise<string> {
	try {
		// Hash the message using EIP-191 format
		const messageHash = hashMessage(message);

		// Recover the public key from the signature
		const publicKey = await recoverPublicKey({
			hash: messageHash,
			signature,
		});

		return publicKey;
	} catch (error) {
		console.error('Failed to recover public key:', error);
		throw error;
	}
}
