pragma solidity ^ 0.5.0;
contract Ownable {
    bytes32 constant private masterPosition = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;
    constructor(address masterAddress) public {
        setMaster(masterAddress);
    }
    function requireMaster(address _address) internal view {
        require(_address == getMaster(), "oro11");
    }
    function getMaster() public view returns (address master) {
        bytes32 position = masterPosition;
        assembly { master := sload(position) }
    }
    function setMaster(address _newMaster) internal {
        bytes32 position = masterPosition;
        assembly { sstore(position, _newMaster) }
    }
    function transferMastership(address _newMaster) external {
        requireMaster(msg.sender);
        require(_newMaster != address(0), "otp11");
        setMaster(_newMaster);
    }
}
pragma solidity ^ 0.5.0;
interface Upgradeable {
    function upgradeTarget(address newTarget, bytes calldata newTargetInitializationParameters) external;
}
pragma solidity ^ 0.5.0;
interface UpgradeableMaster {
    function getNoticePeriod() external returns (uint);
    function upgradeNoticePeriodStarted() external;
    function upgradePreparationStarted() external;
    function upgradeCanceled() external;
    function upgradeFinishes() external;
    function isReadyForUpgrade() external returns (bool);
}
pragma solidity ^ 0.5.0;
contract Proxy is Upgradeable, UpgradeableMaster, Ownable {
    bytes32 constant private targetPosition = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    constructor(address target, bytes memory targetInitializationParameters) Ownable(msg.sender) public {
        setTarget(target);
        (bool initializationSuccess, ) = getTarget().delegatecall(abi.encodeWithSignature("initialize(bytes)", targetInitializationParameters));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(initializationSuccess, "uin11");
    }
    function initialize(bytes calldata) external pure {
        revert("ini11");
    }
    function upgrade(bytes calldata) external pure {
        revert("upg11");
    }
    function getTarget() public view returns (address target) {
        bytes32 position = targetPosition;
        assembly { target := sload(position) }
    }
    function setTarget(address _newTarget) internal {
        bytes32 position = targetPosition;
        assembly { sstore(position, _newTarget) }
    }
    function upgradeTarget(address newTarget, bytes calldata newTargetUpgradeParameters) external {
        requireMaster(msg.sender);
        setTarget(newTarget);
        (bool upgradeSuccess, ) = getTarget().delegatecall(abi.encodeWithSignature("upgrade(bytes)", newTargetUpgradeParameters));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(upgradeSuccess, "ufu11");
    }
    function() external payable {
        address _target = getTarget();
        assembly {
            let ptr := mload(0x40)
            calldatacopy(ptr, 0x0, calldatasize())
            let result := delegatecall(gas(), _target, ptr, calldatasize(), 0x0, 0)
            let size := returndatasize()
            returndatacopy(ptr, 0x0, size)
            switch result
            case 0 { revert(ptr, size) }
            default { return(ptr, size) }
        }
    }
    function getNoticePeriod() external returns (uint) {
        (bool success, bytes memory result) = getTarget().delegatecall(abi.encodeWithSignature("getNoticePeriod()"));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(success, "unp11");
        return abi.decode(result, (uint));
    }
    function upgradeNoticePeriodStarted() external {
        requireMaster(msg.sender);
        (bool success, ) = getTarget().delegatecall(abi.encodeWithSignature("upgradeNoticePeriodStarted()"));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(success, "nps11");
    }
    function upgradePreparationStarted() external {
        requireMaster(msg.sender);
        (bool success, ) = getTarget().delegatecall(abi.encodeWithSignature("upgradePreparationStarted()"));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(success, "ups11");
    }
    function upgradeCanceled() external {
        requireMaster(msg.sender);
        (bool success, ) = getTarget().delegatecall(abi.encodeWithSignature("upgradeCanceled()"));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(success, "puc11");
    }
    function upgradeFinishes() external {
        requireMaster(msg.sender);
        (bool success, ) = getTarget().delegatecall(abi.encodeWithSignature("upgradeFinishes()"));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(success, "puf11");
    }
    function isReadyForUpgrade() external returns (bool) {
        (bool success, bytes memory result) = getTarget().delegatecall(abi.encodeWithSignature("isReadyForUpgrade()"));
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(success, "rfu11");
        return abi.decode(result, (bool));
    }
}
