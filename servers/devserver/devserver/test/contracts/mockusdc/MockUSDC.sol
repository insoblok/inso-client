// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockUSDC {
    string public name = "Mock USD Coin";
    string private symbol = "USDC";
    uint8 private decimals = 6;
    uint256 private totalSupply;

    mapping(address => uint256) private balanceOf;
    mapping(address => mapping(address => uint256)) private allowance;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    // 🧪 Debug events
    event Debug(string message);

    constructor() {
        emit Debug("Start constructor");

        uint256 initialSupply = 1_000_000 * 1e6;
        emit Debug("Initial supply calculated");

        balanceOf[msg.sender] = initialSupply;
        emit Debug("Balance assigned to msg.sender");

        totalSupply = initialSupply;
        emit Debug("Total supply set");

        emit Transfer(address(0), msg.sender, initialSupply);
        emit Debug("Transfer event emitted");

        emit Debug("Constructor done");
    }

//    function transfer(address to, uint256 amount) public returns (bool) {
//        require(balanceOf[msg.sender] >= amount, "Insufficient balance");
//        balanceOf[msg.sender] -= amount;
//        balanceOf[to] += amount;
//        emit Transfer(msg.sender, to, amount);
//        return true;
//    }

//    function approve(address spender, uint256 amount) public returns (bool) {
//        allowance[msg.sender][spender] = amount;
//        emit Approval(msg.sender, spender, amount);
//        return true;
//    }
//
//    function transferFrom(address from, address to, uint256 amount) public returns (bool) {
//        require(balanceOf[from] >= amount, "Insufficient balance");
//        require(allowance[from][msg.sender] >= amount, "Not allowed");
//        balanceOf[from] -= amount;
//        balanceOf[to] += amount;
//        allowance[from][msg.sender] -= amount;
//        emit Transfer(from, to, amount);
//        return true;
//    }
}
