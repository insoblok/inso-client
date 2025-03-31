// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 private count;

    event Increment(address indexed sender, uint256 newValue);

    constructor() {
        count = 0;
    }

    function increment() public {
        count += 1;
        emit Increment(msg.sender, count);
    }

    function get() public view returns (uint256) {
        return count;
    }
}
