<script lang="ts">
import { ERC20_WORMHOLE_TOKEN } from '$lib/constants';
import WormholeTokenAbi from '$lib/abis/WormholeToken.json';
import { writeContract } from 'viem/actions';
import type { Address } from 'viem';
import { getPublicClient, getWalletClient } from '$lib/utils/wallet';
import { sepolia } from 'viem/chains';

const client = getPublicClient();
const walletClient = getWalletClient();


async function mintFaucet() {
    const [account] = await walletClient.getAddresses()

    const availableTokens = await client.readContract({
        abi: WormholeTokenAbi.abi as any,
        address: ERC20_WORMHOLE_TOKEN,
        functionName: 'amountFreeTokens',
        args: [],
    })

    console.log({ availableTokens })

    const request = await walletClient.writeContract({
        abi: WormholeTokenAbi.abi as any,
        address: ERC20_WORMHOLE_TOKEN,
        functionName: 'getFreeTokens',
        args: [account],
        chain: sepolia,
        account: account,
    })

    if (!request) {
        console.error('Failed to request transaction');
        return;
    }

    const tx = await writeContract(client, request)
}

</script>

<div onclick={mintFaucet}>
    <h1>~GET DAT BAG~</h1>
</div>