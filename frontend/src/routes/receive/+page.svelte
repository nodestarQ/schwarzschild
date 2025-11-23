<script lang="ts">
  import { Card, CardContent } from "$lib/components/ui/card";
  import {
    Alert,
    AlertTitle,
    AlertDescription,
  } from "$lib/components/ui/alert";
  import { Button } from "$lib/components/ui/button";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import WalletHeader from "$lib/components/WalletHeader.svelte";
  import WalletConnectModal from "$lib/components/WalletConnectModal.svelte";
  import MetakeyFlow from "$lib/components/flows/MetakeyFlow.svelte";
  import BalanceDisplay from "$lib/components/BalanceDisplay.svelte";
  import { walletModal } from "$lib/store/walletModal";
  import { navigateHome } from "$lib/utils/navigation";
  import { getBalanceForERC20 } from "$lib/utils/balance";
  import { getMetakey } from "$lib/utils/metakey";
  import { formatAddressForDisplay } from "$lib/utils/wallet";
  import type { Address } from "viem";
  import { ERC20_WORMHOLE_TOKEN } from "$lib/constants";

  type ReceiveStep = "setup" | "create-metakey" | "display";

  let currentStep = $state<ReceiveStep>("setup");
  let connectedAddress: Address | null = $state(null);
  let metakeyValue: string | null = $state(null);
  let totalBalance = $state("0");
  let isLoading = $state(false);
  let error: string | null = $state(null);

  async function loadReceiveData() {
    isLoading = true;
    error = null;

    try {
      if (!window.ethereum) {
        error = "No wallet connected";
        return;
      }

      const accounts = await window.ethereum.request({
        method: "eth_accounts",
      });
      if (!accounts || accounts.length === 0) {
        error = "No wallet connected";
        return;
      }

      const address = accounts[0] as Address;
      connectedAddress = address;

      // Check if address has ENS name
      const ensName = await formatAddressForDisplay(address, true);
      if (ensName === address) {
        error = "Please register this address on ENS to receive payments";
        currentStep = "setup";
        return;
      }

      // Get balance
      const balanceData = await getBalanceForERC20(address, ERC20_WORMHOLE_TOKEN);
      if (!balanceData) {
        error = "Could not fetch balance. Please try again.";
        return;
      }

      totalBalance = balanceData?.formatted || "0";

      // Check if ENS already has MetaKey record
      metakeyValue = await getMetakey(ensName);

      // If MetaKey exists, show display. If not, go to creation step
      currentStep = metakeyValue ? "display" : "create-metakey";
    } catch (err) {
      error = `Error: ${err instanceof Error ? err.message : "Unknown error"}`;
    } finally {
      isLoading = false;
    }
  }

  function handleMetakeyCreated(metakey: string) {
    metakeyValue = metakey;
    currentStep = "display";
  }

  // Initialize on component mount
  $effect(() => {
    loadReceiveData();
  });
</script>

<main class="min-h-screen bg-background p-4">
  <WalletHeader onOpenWalletModal={() => walletModal.open()} />

  <div class="max-w-2xl mx-auto">
    <PageHeader
      title="Receive Crypto"
      subtitle={isLoading ? "Loading..." : "View your balance and MetaKey"}
      onBack={navigateHome}
    />

    <!-- Content -->
    {#if error}
      <Alert variant="destructive" class="mb-4">
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>{error}</AlertDescription>
      </Alert>
      <Card>
        <CardContent class="py-8 text-center">
          <Button class="w-full" onclick={navigateHome}>Back to Home</Button>
        </CardContent>
      </Card>
    {:else if isLoading}
      <Card>
        <CardContent class="py-8 text-center">
          <p class="text-muted-foreground">Loading your balance...</p>
        </CardContent>
      </Card>
    {:else if currentStep === "create-metakey" && connectedAddress}
      <MetakeyFlow
        {connectedAddress}
        onSuccess={handleMetakeyCreated}
        onBack={navigateHome}
      />
    {:else if currentStep === "display" && connectedAddress}
      <BalanceDisplay
        {totalBalance}
        {metakeyValue}
        {connectedAddress}
        onBack={navigateHome}
      />
    {:else}
      <Card>
        <CardContent class="py-8 text-center">
          <p class="text-muted-foreground">No wallet connected</p>
          <p class="text-sm text-muted-foreground mt-2">
            Please connect a wallet using the header button
          </p>
        </CardContent>
      </Card>
    {/if}
  </div>

  <WalletConnectModal
    open={$walletModal}
    onOpenChange={(open) => {
      if (open) walletModal.open();
      else walletModal.close();
    }}
    onConnect={() => {}}
  />
</main>
