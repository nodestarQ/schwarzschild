import { hexToBigInt, getAddress, toHex, type Hex, type Address} from "viem";
import { poseidon2Hash } from "@zkpassport/poseidon2"
import { PRIVATE_ADDRESS_TYPE } from "./constants"

const hashAddress = (eoaPubKeyX: Hex) => {
    const pubKeyField = hexToBigInt("0x" + eoaPubKeyX.slice(2 + 2) as Hex)
    return pubKeyField
}

export const getBurnAddress = (eoaPubKeyX: Hex, stealthAddress: Address): Address => {
    const pubKeyField = hashAddress(eoaPubKeyX)
    const stealthAddressField = hashAddress(stealthAddress)

    const hash = poseidon2Hash([pubKeyField, stealthAddressField, PRIVATE_ADDRESS_TYPE])

    const _x = toHex(hash, { size: 32 })

    const _v = _x.slice(2, 40 + 2)

    return getAddress('0x' + _v)  //slice off bytes and make it the address type in viem

}