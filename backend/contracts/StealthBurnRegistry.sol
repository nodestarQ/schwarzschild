// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.0.0
pragma solidity ^0.8.3;

event StealthBurn(bytes32 ephemeralPublicKey, address burnAddress);

contract StealthBurnRegistry {

    constructor() {
    }

    function emitBurn(bytes32 ephemeralPublicKey, address burnAddress) public {
        emit StealthBurn(ephemeralPublicKey, burnAddress);
    }
}