// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CounterV2 {
    uint256 private count;

    event Incremented(address indexed by, uint256 newValue);
    event Decremented(address indexed by, uint256 newValue);
    event Reset(address indexed by);

    event Debug(string message); // Add this at the top

    constructor() {
        emit Debug("INIT");
        count = 0;
        emit Debug("DONE INIT");
    }

    function increment() external {
        count += 1;
        emit Incremented(msg.sender, count);
    }

    function decrement() external {
        require(count > 0, "Counter is already zero");
        count -= 1;
        emit Decremented(msg.sender, count);
    }

    function reset() external {
        count = 0;
        emit Reset(msg.sender);
    }

    function getCount() external view returns (uint256) {
        return count;
    }
}
