<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { walletModal } from '$lib/store/walletModal';
	import { navigateSend, navigateReceive } from '$lib/utils/navigation';
	import WalletConnectModal from '$lib/components/WalletConnectModal.svelte';
	import WalletHeader from '$lib/components/WalletHeader.svelte';
	import TransactionHistory from '$lib/components/TransactionHistory.svelte';

	function handleWalletConnect() {
		walletModal.close();
	}
</script>

<main class="min-h-screen flex flex-col items-center justify-center p-4 bg-background">
	<WalletHeader onOpenWalletModal={() => walletModal.open()} />

	<div class="w-full max-w-2xl">
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold mb-2">Schwarzschild ◼️🛡️</h1>
			<p class="text-lg text-muted-foreground">Send and receive crypto with MetaKeys</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
			<!-- Send Card -->
			<Card class="flex flex-col h-full hover:shadow-lg transition-shadow">
				<CardHeader class="flex-1">
					<CardTitle class="text-2xl">Send Money</CardTitle>
					<CardDescription>
						Transfer crypto to any address using MetaKeys
					</CardDescription>
				</CardHeader>
				<CardContent class="pt-0">
					<p class="text-sm text-muted-foreground mb-4">
						Search for a wallet address or ENS name, verify their MetaKey, and send them crypto securely.
					</p>
					<Button
						onclick={navigateSend}
						class="w-full bg-primary hover:bg-primary/90 text-primary-foreground"
					>
						Send Crypto
					</Button>
				</CardContent>
			</Card>

			<!-- Receive Card -->
			<Card class="flex flex-col h-full hover:shadow-lg transition-shadow">
				<CardHeader class="flex-1">
					<CardTitle class="text-2xl">Receive Money</CardTitle>
					<CardDescription>
						Set up your MetaKey and monitor your balance
					</CardDescription>
				</CardHeader>
				<CardContent class="pt-0">
					<p class="text-sm text-muted-foreground mb-4">
						Connect your wallet, create a MetaKey if needed, and view your total balance across all tokens.
					</p>
					<Button
						onclick={navigateReceive}
						class="w-full bg-secondary hover:bg-secondary/90 text-secondary-foreground"
					>
						Receive Crypto
					</Button>
				</CardContent>
			</Card>
		</div>

		<!-- Transaction History -->
		<TransactionHistory />
	</div>
</main>

<WalletConnectModal
	open={$walletModal}
	onOpenChange={(open) => {
		if (open) walletModal.open();
		else walletModal.close();
	}}
	onConnect={handleWalletConnect}
/>

<style>
	:global(body) {
		background-color: var(--background);
	}
</style>
