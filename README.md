# Schwarzschild (BlackShield) ◼️🛡️

![123](./docs/images/image.png)

We believe that privacy should be a first-class citizen on Ethereum. Unfortunately, existing solutions like Tornado Cash, Railgun, and Privacy Pools struggle with user experience. That's why we decided to explore the promising EIP-7503, which enables faster finality for private transfers.

⚠️ **This repository builds on [jimjim's EIP-7503 implementation](https://github.com/jimjimvalkema/EIP7503-ERC20), extending it to support ERC-20 tokens, hardware wallets, and more.** However, the original implementation required users to maintain a separate private key for minting funds after the link between sender and recipient was broken which is a significant UX hurdle.

We focused on dramatically improving user experience by implementing a Diffie-Hellman-style key exchange inspired by Stealth Addresses. For users, this means that after a simple onboarding process, they can receive private transfers using their ENS name and seamlessly decide what to do with their funds once the link was broken.

## The Original EIP-7503 vs ERC-20 Compatible EIP-7503

The original EIP-7503 introduced important privacy features but required users to manage single-use addresses and lacked robust hardware wallet support, making secure and convenient interaction difficult for typical users. 

This repository improves usability and flexibility by:
- **Enabling address reuse** – No need to generate new addresses for each transaction
- **Simplifying balance tracking** – Users can easily monitor their private balances
- **Supporting hardware wallets** – Secure key management with familiar devices
- **Application-layer implementation** – Moving privacy features to smart contracts allows for upgrades and improvements while maintaining accessibility

By implementing these changes at the application layer, we can iterate on the protocol and enhance the user experience while preserving the core privacy guarantees.

## ERC-5564 (Stealth Addresses) and ERC-6538 (Stealth Meta-Address Registry) role 

This existing and not yet so widely adopted ERC may already allow users to do transfers (and most-importantly self-transfer) to break the link with source of funds. However it introduced new problems with neccesity to either fund these addresses from original wallet what break anonymity or introduce account abstraction which adds additional costs and techincal complexity. However mechanism of both parties accessing some shared secret (in original ERC it's ephemeral private key) to create a communication channel that protects funds reciever, but not the sender. 

## The Land in-between

Once we were discussing existing solutions we sparked that ZK-Wormholes and Stealth Addresses are nicely complimenting each other but require some additional adjustments to adopt to new specifics. Such as using shared secret as a salt for burn address generation and completely unnessesary mechanism of controling the private keys of burner addresses. 

Solution that we came up with includes not using generated stealth addresses to recieve funds but just discover a shared secret. The trick was to use stealth address as a shared secret itself and not post it onchain. In this case Alice can always derrive this secret as well as bob on his side. And the only modification is needed is to replace `stealthAddress` in the `Announcement` event with the `burnAddress` which Bob can derrive on his side and understand, that private note was meant to him and there is some balance to be claimed.

## WARNINGS
* The value `POW_DIFFICULTY` has been set to an arbitrary number and **IS LIKELY INSECURE**   
https://github.com/jimjimvalkema/EIP7503-ERC20/blob/7a485j0ddc6503442dfbd484ac3754a1bd0c02796/circuits/privateTransfer/src/main.nr#L13  
https://github.com/jimjimvalkema/EIP7503-ERC20/blob/7a4850ddc6503442dfbd484ac3754a1bd0c02796/src/constants.ts#L17  
* Compliance (not legal advice ofc): the viewing_key can be used to reveal transaction history but that use case needs more research. It should also be possible to make a PoI scheme work without modifying the circuits/contracts like on railgun. As of now this repo doesn't provide tools for compliance.  
* This is unaudited and experimental. The poseidon2 contract is also experimental and built with huff: https://github.com/zemse/poseidon2-evm  

## deploy
setup secrets:  
`yarn hardhat keystore set SEPOLIA_RPC_URL`  
`yarn hardhat keystore set SEPOLIA_PRIVATE_KEY`  
`yarn hardhat keystore set ETHERSCAN_API_KEY`  

deploy main contracts:  
```shell
yarn hardhat ignition deploy ignition/modules/wormtoken.ts --verify --network sepolia
yarn hardhat ignition deploy ignition/modules/stealth-burn-registry.ts --verify --network sepolia
```  

deploy poseidon2 hasher with create2 (if it's not deployed yet)
```shell
yarn hardhat run scripts/deployPoseidon2.ts --network sepolia
```

## deployed addresses
### sepolia 
WormholeToken - [0x67Cc5Ac2029aaA9FD56F7D036d61f2d80A034c10](https://sepolia.etherscan.io/address/0x67Cc5Ac2029aaA9FD56F7D036d61f2d80A034c10)  


PrivateTransferVerifier - [0x342149C7108bb2b0052624f61629f5813B9B9466](https://sepolia.etherscan.io/address/0x342149C7108bb2b0052624f61629f5813B9B9466)  
ZKTranscriptLib - [0x8F961e056967DD2A1170dBeCd9e5E51CA815B0D9](https://sepolia.etherscan.io/address/0x8F961e056967DD2A1170dBeCd9e5E51CA815B0D9)  
leanIMTPoseidon2 - [0xcbf45ce9650A8F4E51933A13857016B1A44c3d94](https://sepolia.etherscan.io/address/0xcbf45ce9650A8F4E51933A13857016B1A44c3d94)  

StealthBurnRegistry — [0x43208AA38104d08EC099d55F709dd46E52ea619C](https://sepolia.etherscan.io/address/0x43208AA38104d08EC099d55F709dd46E52ea619C)


## run ui in localhost

```shell
cd fronted/ ;
bun install ; 
bun run dev ;
```

## Authors

@TODO