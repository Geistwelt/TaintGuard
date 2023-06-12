/**
 *Submitted for verification at Etherscan.io on 2021-04-08
*/

// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract Ownable {
    address private _owner;
    address private _test;

    // constructor() {
    //     _owner = msg.sender;
    // }

    function owner() internal view returns (address) {
        return _owner;
    }

    function test() public view returns (address) {
        return _test;
    }
}

contract A is Ownable {
    address private _b;

    // constructor(address _lib) {
    //     _b = _lib;
    // }

    // function b() public view returns (address) {
    //     return  _b;
    // }

    // function call1() public {
    //     _b.delegatecall(abi.encodeWithSignature("modi1()"));
    // }

    // function call2() public {
    //     _b.delegatecall(abi.encodeWithSignature("modi2()"));
    // }
}

contract B is A {
    address private _owner;
    bytes private _test;

    mapping(bytes=>address) private xxx_track_mapping_owner_;

    function func() public {
        xxx_track_mapping_owner_["haha"] = address(this);
        _test = "haha";
        
        _owner.delegatecall(msg.data);
        assert(xxx_track_mapping_owner_[_test] == owner());
    }

    // function modi1() public {
    //     owner = msg.sender;
    // }

    // function modi2() public {
    //     test = msg.sender;
    // }
}