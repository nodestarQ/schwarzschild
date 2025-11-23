<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
  } from "$lib/components/ui/card";
  import { Alert, AlertTitle, AlertDescription } from "$lib/components/ui/alert";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import ProgressBar from "$lib/components/ProgressBar.svelte";
  import WalletSearch from "$lib/components/flows/WalletSearch.svelte";
  import SendAmount from "$lib/components/flows/SendAmount.svelte";
  import SendTransaction from "$lib/components/flows/SendTransaction.svelte";
  import WalletHeader from "$lib/components/WalletHeader.svelte";
  import WalletConnectModal from "$lib/components/WalletConnectModal.svelte";
  import { walletModal } from "$lib/store/walletModal";
  import { navigateHome } from "$lib/utils/navigation";
  import type { Address, Hex } from "viem";
  import { CheckCircle2 } from "@lucide/svelte";

  type SendStep = "search" | "amount" | "confirm" | "success";
  import { getMetakey } from "$lib/utils/metakey";
  import { getBurnAddress } from "$lib/burn";

  let selectedAddress: Address | null = $state(null);
  let ensName: string | null = $state(null);
  let sendAmount: string | null = $state(null);
  let transactionHash: string | null = $state(null);

  const STEPS: SendStep[] = ["search", "amount", "confirm", "success"];
  const STEP_TITLES: Record<SendStep, string> = {
    search: "Find Recipient",
    amount: "Send Amount",
    confirm: "Confirm Transaction",
    success: "Transaction Complete",
  };
  const STEP_DESCRIPTIONS: Record<SendStep, string> = {
    search: "Enter a wallet address or ENS name. We'll verify they have a MetaKey.",
    amount: "Enter the amount you want to send to this recipient.",
    confirm: "Review the transaction details before sending.",
    success: "Your transaction has been submitted to the network.",
  };
  let transactionError: string | null = $state(null);
  let currentStep = $state<"search" | "amount" | "confirm" | "success">("search");
  let stealthMetaAddress: string | null = $state(null);
  let burnAddress: Address | null = $state(null);

  $effect(() => {
    if (!selectedAddress) return
    getMetakey(selectedAddress).then(metakey => {
      stealthMetaAddress = metakey

      // skip first 141 chars as they are:
      // "st:eth:"
      // 0xViewingPublicKey
      // 0xSpendingPublicKey
      // and only get:
      // 0xPubKeyX

      const pubKeyX = metakey?.slice(141, 141 + 64)

      burnAddress = getBurnAddress(pubKeyX as Hex, selectedAddress as Address)
      console.log({ burnAddress })
    })
  })


  function handleBack() {
    router.navigate("home");
  }

  function handleWalletSelect(
    address: Address,
    name: string,
    _hasKey: boolean,
  ) {
    selectedAddress = address;
    ensName = name;
    currentStep = "amount";
  }

  function handleAmountSubmit(amount: string) {
    sendAmount = amount;
    currentStep = "confirm";
  }

  function handleBack2Steps() {
    currentStep = "search";
    selectedAddress = null;
    sendAmount = null;
  }

  function handleBackFromAmount() {
    currentStep = "search";
    selectedAddress = null;
  }

  function handleTransactionSuccess(hash: string) {
    transactionHash = hash;
    currentStep = "success";
  }

  function resetFlow() {
    currentStep = "search";
    selectedAddress = null;
    ensName = null;
    sendAmount = null;
    transactionHash = null;
  }
</script>

<main class="min-h-screen bg-background p-4">
  <WalletHeader onOpenWalletModal={() => walletModal.open()} />
  
  <div class="max-w-2xl mx-auto">
    <PageHeader
      title="Send Crypto"
      onBack={navigateHome}
    />

    <ProgressBar {currentStep} steps={STEPS} />

    <!-- Content -->
    <Card>
       <CardHeader>
         <CardTitle>{STEP_TITLES[currentStep]}</CardTitle>
         <CardDescription>
           {STEP_DESCRIPTIONS[currentStep]}
         </CardDescription>
       </CardHeader>

       <CardContent>
         {#if currentStep === "search"}
           <WalletSearch onSelect={handleWalletSelect} />
         {:else if currentStep === "amount"}
           {#if selectedAddress}
             <SendAmount
               recipientAddress={selectedAddress}
               recipientName={ensName}
               burnAddress={burnAddress}
               onSubmit={handleAmountSubmit}
               onBack={handleBackFromAmount}
             />
           {/if}
         {:else if currentStep === "confirm"}
           {#if selectedAddress && sendAmount}
             <SendTransaction
               recipientAddress={selectedAddress}
               recipientName={ensName}
               amount={sendAmount}
               burnAddress={burnAddress}
               stealthMetaAddress={stealthMetaAddress}
               onSuccess={handleTransactionSuccess}
               onError={() => {}}
               onBack={handleBack2Steps}
             />
           {/if}
         {:else}
           <!-- Success State -->
           <div class="space-y-4">
             <Alert class="border-green-500/50 bg-green-500/10">
               <CheckCircle2 size={16} class="text-green-600" />
               <AlertTitle class="text-green-600">Transaction Submitted!</AlertTitle>
               <AlertDescription class="text-green-700">
                 Your transaction has been submitted to the network. You can check its status on the block explorer.
               </AlertDescription>
             </Alert>

             <div class="p-4 bg-card border border-border rounded-md">
               <p class="text-xs text-muted-foreground mb-1">Recipient</p>
               <p class="font-mono text-sm break-all">{ensName || selectedAddress}</p>
             </div>

             <div class="p-4 bg-card border border-border rounded-md">
               <p class="text-xs text-muted-foreground mb-1">Amount</p>
               <p class="text-2xl font-bold">{sendAmount} ETH</p>
             </div>

             {#if transactionHash}
               <div class="p-4 bg-card border border-border rounded-md">
                 <p class="text-xs text-muted-foreground mb-1">Transaction Hash</p>
                 <p class="font-mono text-xs break-all">{transactionHash}</p>
               </div>
             {/if}

             <div class="flex gap-2">
               <Button variant="outline" class="flex-1" onclick={resetFlow}>
                 Send Another
               </Button>
               <Button class="flex-1" onclick={navigateHome}>
                 Back to Home
               </Button>
             </div>
           </div>
         {/if}
       </CardContent>
    </Card>
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
