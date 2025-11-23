<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
  } from "$lib/components/ui/card";
  import {
    Alert,
    AlertTitle,
    AlertDescription,
  } from "$lib/components/ui/alert";
  import { setMetakey } from "$lib/utils/metakey";
  import { getEnsName } from "$lib/utils/wallet";
  import { signMessageWithWallet, recoverSignaturePublicKey } from "$lib/utils/signing";
  import {
    toHex,
    type Address,
    type Hex,
  } from "viem";
  import { getStealthMetaAddress } from "$lib/stealth/prepare-keys";
  import { currentProvider } from "$lib/providers";

  interface Props {
    connectedAddress: Address | null;
    onSuccess?: (metakey: string) => void;
    onError?: (error: string) => void;
    onBack?: () => void;
  }

  let { connectedAddress, onSuccess, onError, onBack }: Props = $props();

  let isCreatingMetakey = $state(false);
  let error: string | null = $state(null);

  async function createMetakey() {
    if (!connectedAddress || !window.ethereum) return;

    if (!$currentProvider) {
      console.error("No provider available");
      return;
    }
    let metaKey: string;
    isCreatingMetakey = true;
    error = null;

    try {
      const message = "Create MetaKey for Schwarzschild";

      // Sign the message using viem's signing function
      const signature = await signMessageWithWallet(message, connectedAddress, $currentProvider);

      if (!signature) throw new Error("Failed to sign message");

      // Recover the public key from the signature using viem
      const publicKey = await recoverSignaturePublicKey(message, signature as Hex);

      try {
        const result = await getStealthMetaAddress(toHex(signature), message);
        const viewingPublicKey = result.viewingPublicKey;
        const spendingPublicKey = result.spendingPublicKey.replace("0x", "");
        const publicKeyX = publicKey.replace("0x", "");

        metaKey = `st:eth:${viewingPublicKey}${spendingPublicKey}${publicKeyX}`;
      } catch (error) {
        console.error("Error generating stealth meta address:", error);
        throw error;
      }

      await setMetakey((await getEnsName(connectedAddress)) as string, metaKey);

      onSuccess?.(metaKey);
    } catch (err) {
      const errorMsg = `Failed to create MetaKey: ${err instanceof Error ? err.message : "Unknown error"}`;
      error = errorMsg;
      onError?.(errorMsg);
    } finally {
      isCreatingMetakey = false;
    }
  }
</script>

<Card>
  {#if error}
    <CardHeader>
      <CardTitle>Please register your ENS first</CardTitle>
      <CardDescription>
        You cannot proceed until you have an ENS
      </CardDescription>
    </CardHeader>
  {:else}
    <CardHeader>
      <CardTitle>Create MetaKey</CardTitle>
      <CardDescription>
        You need a MetaKey to receive payments securely
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <Alert>
        <AlertTitle>Sign to Create MetaKey</AlertTitle>
        <AlertDescription>
          Click the button below to sign a message with your wallet. This will
          generate your unique MetaKey.
        </AlertDescription>
      </Alert>
      <div class="flex gap-2">
        <Button
          variant="outline"
          class="flex-1"
          onclick={onBack}
          disabled={isCreatingMetakey}
        >
          Back
        </Button>
        <Button
          class="flex-1"
          onclick={createMetakey}
          disabled={isCreatingMetakey}
        >
          {isCreatingMetakey ? "Creating MetaKey..." : "Create MetaKey"}
        </Button>
      </div>
    </CardContent>
  {/if}
</Card>
