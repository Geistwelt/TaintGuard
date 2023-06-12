// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract Lib {
    address private _owner;

    constructor() {
        _owner = msg.sender;
    }

    function owner() public view returns (address) {
        return _owner;
    }
}

contract Hack {
    Lib private lib;

    constructor(Lib _lib) {
        lib = _lib;
    }

    function call() public {
        address(lib).delegatecall(msg.data);
    }
}