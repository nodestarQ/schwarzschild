<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Alert, AlertTitle, AlertDescription } from '$lib/components/ui/alert';
	import CopyableText from './CopyableText.svelte';
	import { fetchStealthBurnLogs } from '$lib/utils/contractLogs';
	import type { Address } from 'viem';
	import { ExternalLink, Loader2 } from '@lucide/svelte';

	interface Props {
		totalBalance: string;
		metakeyValue: string | null;
		connectedAddress: Address | null;
		onBack?: () => void;
	}

	let { totalBalance, metakeyValue, connectedAddress, onBack }: Props = $props();

	let copiedIndex = $state(-1);
	let logs = $state<any[]>([]);
	let logsLoading = $state(false);
	let logsError = $state<string | null>(null);

	const CONTRACT_ADDRESS = (import.meta.env.VITE_CONTRACT_ADDRESS ||
		'0x43208AA38104d08EC099d55F709dd46E52ea619C') as Address;

	function handleCopy(index: number) {
		copiedIndex = index;
		setTimeout(() => {
			copiedIndex = -1;
		}, 2000);
	}

	async function loadLogs() {
		try {
			logsLoading = true;
			logsError = null;
			logs = await fetchStealthBurnLogs(CONTRACT_ADDRESS, 50);
		} catch (error) {
			logsError = error instanceof Error ? error.message : 'Failed to load logs';
			console.error('Error loading logs:', error);
		} finally {
			logsLoading = false;
		}
	}

	function formatAddress(address: string) {
		if (!address) return '';
		return `${address.slice(0, 6)}...${address.slice(-4)}`;
	}

	function formatPublicKey(key: string) {
		if (!key) return '';
		return `${key.slice(0, 6)}...${key.slice(-4)}`;
	}

	function getExplorerUrl(txHash: string) {
		return `https://sepolia.etherscan.io/tx/${txHash}`;
	}

	// Load logs when component mounts
	$effect(() => {
		loadLogs();
	});
</script>

<div class="space-y-4">
	<Card>
		<CardHeader>
			<CardTitle>Your Balance</CardTitle>
			<CardDescription>
				Total balance across all your addresses
			</CardDescription>
		</CardHeader>
		<CardContent class="space-y-6">
			<!-- Total Balance -->
			<div
				class="p-6 bg-gradient-to-br from-primary/10 to-accent/10 rounded-lg border border-primary/20"
			>
				<p class="text-sm text-muted-foreground mb-2">Total Balance</p>
				<p class="text-4xl font-bold">{totalBalance} ETH</p>
			</div>

			<!-- MetaKey Display -->
			{#if metakeyValue}
				<CopyableText
					label="Your MetaKey"
					text={metakeyValue}
					isCopied={copiedIndex === 0}
					onCopy={() => handleCopy(0)}
				/>
				<p class="text-xs text-muted-foreground">
					Share this MetaKey with others to receive payments securely.
				</p>
			{/if}

			<!-- Connected Address -->
			{#if connectedAddress}
				<CopyableText
					label="Your Address"
					text={connectedAddress}
					isCopied={copiedIndex === 1}
					onCopy={() => handleCopy(1)}
				/>
			{/if}

			<Button class="w-full" onclick={onBack}>Back to Home</Button>
		</CardContent>
	</Card>

	<!-- Contract Logs Section -->
	<Card>
		<CardHeader>
			<CardTitle>Recent Stealth Burns</CardTitle>
			<CardDescription>
				Recent burn transactions from the StealthBurnRegistry contract
			</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			{#if logsLoading}
				<div class="flex items-center justify-center py-8">
					<Loader2 class="animate-spin mr-2" size={20} />
					<span>Loading logs...</span>
				</div>
			{:else if logsError}
				<Alert variant="destructive">
					<AlertTitle>Error Loading Logs</AlertTitle>
					<AlertDescription>{logsError}</AlertDescription>
				</Alert>
			{:else if logs.length === 0}
				<Alert>
					<AlertTitle>No Logs Found</AlertTitle>
					<AlertDescription>No stealth burn events have been recorded yet.</AlertDescription>
				</Alert>
			{:else}
				<div class="space-y-3">
					{#each logs as log, index}
						<div class="p-3 bg-card border border-border rounded-md hover:bg-accent/50 transition-colors">
							<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
								<!-- Block Number -->
								<div>
									<p class="text-xs text-muted-foreground mb-1">Block</p>
									<p class="font-mono text-sm font-semibold">{log.blockNumber.toString()}</p>
								</div>

								<!-- Burn Address -->
								<div>
									<p class="text-xs text-muted-foreground mb-1">Burn Address</p>
									<p class="font-mono text-sm break-all">{formatAddress(log.burnAddress)}</p>
								</div>

								<!-- Ephemeral Public Key -->
								<div class="md:col-span-2">
									<p class="text-xs text-muted-foreground mb-1">Ephemeral Public Key</p>
									<p class="font-mono text-xs break-all text-muted-foreground">{log.ephemeralPublicKey}</p>
								</div>

								<!-- Transaction Hash -->
								<div class="md:col-span-2">
									<p class="text-xs text-muted-foreground mb-1">Transaction</p>
									<div class="flex items-center gap-2">
										<p class="font-mono text-xs break-all">{formatAddress(log.transactionHash)}</p>
										<a
											href={getExplorerUrl(log.transactionHash)}
											target="_blank"
											rel="noreferrer"
											class="text-primary hover:underline flex-shrink-0"
										>
											<ExternalLink size={14} />
										</a>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>

				{#if logs.length > 0}
					<Button variant="outline" class="w-full mt-4" onclick={loadLogs} disabled={logsLoading}>
						Refresh Logs
					</Button>
				{/if}
			{/if}
		</CardContent>
	</Card>
</div>
