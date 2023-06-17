pragma solidity ^ 0.6.0;
interface IProxy {
    function core() external view returns (address);
}
pragma solidity ^ 0.6.0;
contract Proxy is IProxy {
    address public core;
    modifier onlyCore() {
        require(core == msg.sender, "PR01");
        _;
    }
    constructor(address _core) public {
        core = _core;
    }
    function updateCore(address _core) public onlyCore returns (bool) {
        core = _core;
        return true;
    }
    function staticCallUint256() internal view returns (uint256 result) {
        (bool status, bytes memory value) = core.staticcall(msg.data);
        require(status, "PR02");
        assembly {
            result := mload(add(value, 0x20))
        }
    }
}
pragma solidity ^ 0.6.0;
contract IAccessDefinitions {
    bytes32 constant ALL_PRIVILEGES = bytes32("AllPrivileges");
    address constant ALL_PROXIES = address(0x416c6C50726F78696573);
    bytes32 constant FACTORY_CORE_ROLE = bytes32("FactoryCoreRole");
    bytes32 constant FACTORY_PROXY_ROLE = bytes32("FactoryProxyRole");
    bytes4 constant DEFINE_ROLE_PRIV = bytes4(keccak256("defineRole(bytes32,bytes4[])"));
    bytes4 constant ASSIGN_OPERATORS_PRIV = bytes4(keccak256("assignOperators(bytes32,address[])"));
    bytes4 constant REVOKE_OPERATORS_PRIV = bytes4(keccak256("revokeOperators(address[])"));
    bytes4 constant ASSIGN_PROXY_OPERATORS_PRIV = bytes4(keccak256("assignProxyOperators(address,bytes32,address[])"));
}
pragma solidity ^ 0.6.0;
abstract contract IOperableStorage is IAccessDefinitions {
    struct RoleData{
        mapping (bytes4 => bool) privileges;
    }
    struct OperatorData{
        bytes32 coreRole;
        mapping (address => bytes32) proxyRoles;
    }
    function coreRole(address _address) public view virtual returns (bytes32);
    function proxyRole(address _proxy, address _address) public view virtual returns (bytes32);
    function rolePrivilege(bytes32 _role, bytes4 _privilege) public view virtual returns (bool);
    function roleHasPrivilege(bytes32 _role, bytes4 _privilege) public view virtual returns (bool);
    function hasCorePrivilege(address _address, bytes4 _privilege) public view virtual returns (bool);
    function hasProxyPrivilege(address _address, address _proxy, bytes4 _privilege) public view virtual returns (bool);
    event RoleDefined(bytes32 role);
    event OperatorAssigned(bytes32 role, address operator);
    event ProxyOperatorAssigned(address proxy, bytes32 role, address operator);
    event OperatorRevoked(address operator);
}
pragma solidity ^ 0.6.0;
abstract contract IOperableCore is IOperableStorage {
    function defineRole(bytes32 _role, bytes4[] memory _privileges) public virtual returns (bool);
    function assignOperators(bytes32 _role, address[] memory _operators) public virtual returns (bool);
    function assignProxyOperators(address _proxy, bytes32 _role, address[] memory _operators) public virtual returns (bool);
    function revokeOperators(address[] memory _operators) public virtual returns (bool);
}
pragma solidity ^ 0.6.0;
contract Ownable {
    address public owner;
    event OwnershipRenounced(address previousOwner);
    event OwnershipTransferred(address previousOwner, address newOwner);
    constructor() public {
        owner = msg.sender;
    }
    modifier onlyOwner() {
        require(msg.sender == owner, "OW01");
        _;
    }
    function renounceOwnership() public onlyOwner {
        emit OwnershipRenounced(owner);
        owner = address(0);
    }
    function transferOwnership(address _newOwner) public onlyOwner {
        _transferOwnership(_newOwner);
    }
    function _transferOwnership(address _newOwner) internal {
        require(_newOwner != address(0), "OW02");
        emit OwnershipTransferred(owner, _newOwner);
        owner = _newOwner;
    }
}
pragma solidity ^ 0.6.0;
contract Storage {
    mapping (address => uint256) public proxyDelegateIds;
    mapping (uint256 => address) public delegates;
}
pragma solidity ^ 0.6.0;
contract OperableStorage is IOperableStorage, Ownable, Storage {
    mapping (address => OperatorData) operators;
    mapping (bytes32 => RoleData) roles;
    function coreRole(address _address) public view override returns (bytes32) {
        return operators[_address].coreRole;
    }
    function proxyRole(address _proxy, address _address) public view override returns (bytes32) {
        return operators[_address].proxyRoles[_proxy];
    }
    function rolePrivilege(bytes32 _role, bytes4 _privilege) public view override returns (bool) {
        return roles[_role].privileges[_privilege];
    }
    function roleHasPrivilege(bytes32 _role, bytes4 _privilege) public view override returns (bool) {
        return (_role == ALL_PRIVILEGES) || roles[_role].privileges[_privilege];
    }
    function hasCorePrivilege(address _address, bytes4 _privilege) public view override returns (bool) {
        bytes32 role = operators[_address].coreRole;
        return (role == ALL_PRIVILEGES) || roles[role].privileges[_privilege];
    }
    function hasProxyPrivilege(address _address, address _proxy, bytes4 _privilege) public view override returns (bool) {
        OperatorData storage data = operators[_address];
        bytes32 role = (data.proxyRoles[_proxy] != bytes32(0))?data.proxyRoles[_proxy]:data.proxyRoles[ALL_PROXIES];
        return (role == ALL_PRIVILEGES) || roles[role].privileges[_privilege];
    }
}
pragma solidity ^ 0.6.0;
library BytesConvert {
    function toUint256(bytes memory _source) internal pure returns (uint256 result) {
        require(_source.length == 32, "BC01");
        assembly {
            result := mload(add(_source, 0x20))
        }
    }
    function toBytes32(bytes memory _source) internal pure returns (bytes32 result) {
        require(_source.length <= 32, "BC02");
        assembly {
            result := mload(add(_source, 0x20))
        }
    }
}
pragma solidity ^ 0.6.0;
contract Core is Storage {
    using BytesConvert for bytes;
    modifier onlyProxy() {
        require(delegates[proxyDelegateIds[msg.sender]] != address(0), "CO01");
        _;
    }
    function delegateCall(address _proxy) internal returns (bool status) {
        uint256 delegateId = proxyDelegateIds[_proxy];
        address delegate = delegates[delegateId];
        require(delegate != address(0), "CO02");
        (status, ) = delegate.delegatecall(msg.data);
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(status, "CO03");
    }
    function delegateCallUint256(address _proxy) internal returns (uint256) {
        return delegateCallBytes(_proxy).toUint256();
    }
    function delegateCallBytes(address _proxy) internal returns (bytes memory result) {
        bool status;
        uint256 delegateId = proxyDelegateIds[_proxy];
        address delegate = delegates[delegateId];
        require(delegate != address(0), "CO02");
        (status, result) = delegate.delegatecall(msg.data);
        assert(xxx_track_mapping_[xxx_track_] == xxx_track_func_());
        require(status, "CO03");
    }
    function defineDelegateInternal(uint256 _delegateId, address _delegate) internal returns (bool) {
        require(_delegateId != 0, "CO04");
        delegates[_delegateId] = _delegate;
        return true;
    }
    function defineProxyInternal(address _proxy, uint256 _delegateId) internal virtual returns (bool) {
        require(delegates[_delegateId] != address(0), "CO02");
        require(_proxy != address(0), "CO05");
        proxyDelegateIds[_proxy] = _delegateId;
        return true;
    }
    function migrateProxyInternal(address _proxy, address _newCore) internal returns (bool) {
        require(proxyDelegateIds[_proxy] != 0, "CO06");
        require(Proxy(_proxy).updateCore(_newCore), "CO07");
        return true;
    }
    function removeProxyInternal(address _proxy) internal returns (bool) {
        require(proxyDelegateIds[_proxy] != 0, "CO06");
        delete proxyDelegateIds[_proxy];
        return true;
    }
}
pragma solidity ^ 0.6.0;
contract OperableCore is IOperableCore, Core, OperableStorage {
    constructor(address[] memory _sysOperators) public {
        assignOperators(ALL_PRIVILEGES, _sysOperators);
        assignProxyOperators(ALL_PROXIES, ALL_PRIVILEGES, _sysOperators);
    }
    modifier onlySysOp() {
        require(msg.sender == owner || hasCorePrivilege(msg.sender, msg.sig), "OC01");
        _;
    }
    modifier onlyCoreOp() {
        require(hasCorePrivilege(msg.sender, msg.sig), "OC02");
        _;
    }
    modifier onlyProxyOp(address _proxy) {
        require(hasProxyPrivilege(msg.sender, _proxy, msg.sig), "OC03");
        _;
    }
    function defineRole(bytes32 _role, bytes4[] memory _privileges) public onlySysOp override returns (bool) {
        require(_role != ALL_PRIVILEGES, "OC04");
        delete roles[_role];
        for (uint256 i = 0; i < _privileges.length; i++) {
            roles[_role].privileges[_privileges[i]] = true;
        }
        emit RoleDefined(_role);
        return true;
    }
    function assignOperators(bytes32 _role, address[] memory _operators) public onlySysOp override returns (bool) {
        for (uint256 i = 0; i < _operators.length; i++) {
            operators[_operators[i]].coreRole = _role;
            emit OperatorAssigned(_role, _operators[i]);
        }
        return true;
    }
    function defineProxyInternal(address _proxy, uint256 _delegateId) internal override returns (bool) {
        require(_proxy != ALL_PROXIES, "OC05");
        return super.defineProxyInternal(_proxy, _delegateId);
    }
    function assignProxyOperators(address _proxy, bytes32 _role, address[] memory _operators) public onlySysOp override returns (bool) {
        for (uint256 i = 0; i < _operators.length; i++) {
            operators[_operators[i]].proxyRoles[_proxy] = _role;
            emit ProxyOperatorAssigned(_proxy, _role, _operators[i]);
        }
        return true;
    }
    function revokeOperators(address[] memory _operators) public onlySysOp override returns (bool) {
        for (uint256 i = 0; i < _operators.length; i++) {
            delete operators[_operators[i]];
            emit OperatorRevoked(_operators[i]);
        }
        return true;
    }
}
pragma solidity ^ 0.6.0;
contract OperableProxy is Proxy {
    constructor(address _core) Proxy(_core) public {

    }
    modifier onlyOperator() {
        require(OperableCore(core).hasProxyPrivilege(msg.sender, address(this), msg.sig), "OP01");
        _;
    }
}
pragma solidity ^ 0.6.0;
interface IERC20 {
    event Transfer(address from, address to, uint256 value);
    event Approval(address owner, address spender, uint256 value);
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint256);
    function totalSupply() external view returns (uint256);
    function balanceOf(address _owner) external view returns (uint256);
    function transfer(address _to, uint256 _value) external returns (bool);
    function allowance(address _owner, address _spender) external view returns (uint256);
    function transferFrom(address _from, address _to, uint256 _value) external returns (bool);
    function approve(address _spender, uint256 _value) external returns (bool);
    function increaseApproval(address _spender, uint256 _addedValue) external returns (bool);
    function decreaseApproval(address _spender, uint256 _subtractedValue) external returns (bool);
}
pragma solidity ^ 0.6.0;
abstract contract ITokenProxy is IERC20, Proxy {
    function canTransfer(address, address, uint256) public view virtual returns (uint256);
    function emitTransfer(address _from, address _to, uint256 _value) public virtual returns (bool);
    function emitApproval(address _owner, address _spender, uint256 _value) public virtual returns (bool);
}
pragma solidity ^ 0.6.0;
library SafeMath {
    function mul(uint256 a, uint256 b) internal pure returns (uint256 c) {
        if(a == 0) {
            return 0;
        }
        c = a * b;
        assert(c / a == b);
        return c;
    }
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        return a / b;
    }
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        assert(b <= a);
        return a - b;
    }
    function add(uint256 a, uint256 b) internal pure returns (uint256 c) {
        c = a + b;
        assert(c >= a);
        return c;
    }
}
pragma solidity ^ 0.6.0;
interface IRule {
    function isAddressValid(address _address) external view returns (bool);
    function isTransferValid(address _from, address _to, uint256 _amount) external view returns (bool);
}
pragma solidity ^ 0.6.0;
abstract contract IUserRegistry {
    enum KeyCode {
        KYC_LIMIT_KEY,
        RECEPTION_LIMIT_KEY,
        EMISSION_LIMIT_KEY
    }
    event UserRegistered(uint256 userId, address address_, uint256 validUntilTime);
    event AddressAttached(uint256 userId, address address_);
    event AddressDetached(uint256 userId, address address_);
    event UserSuspended(uint256 userId);
    event UserRestored(uint256 userId);
    event UserValidity(uint256 userId, uint256 validUntilTime);
    event UserExtendedKey(uint256 userId, uint256 key, uint256 value);
    event UserExtendedKeys(uint256 userId, uint256[] values);
    event ExtendedKeysDefinition(uint256[] keys);
    function registerManyUsersExternal(address[] calldata _addresses, uint256 _validUntilTime) external virtual returns (bool);
    function registerManyUsersFullExternal(address[] calldata _addresses, uint256 _validUntilTime, uint256[] calldata _values) external virtual returns (bool);
    function attachManyAddressesExternal(uint256[] calldata _userIds, address[] calldata _addresses) external virtual returns (bool);
    function detachManyAddressesExternal(address[] calldata _addresses) external virtual returns (bool);
    function suspendManyUsersExternal(uint256[] calldata _userIds) external virtual returns (bool);
    function restoreManyUsersExternal(uint256[] calldata _userIds) external virtual returns (bool);
    function updateManyUsersExternal(uint256[] calldata _userIds, uint256 _validUntilTime, bool _suspended) external virtual returns (bool);
    function updateManyUsersExtendedExternal(uint256[] calldata _userIds, uint256 _key, uint256 _value) external virtual returns (bool);
    function updateManyUsersAllExtendedExternal(uint256[] calldata _userIds, uint256[] calldata _values) external virtual returns (bool);
    function updateManyUsersFullExternal(uint256[] calldata _userIds, uint256 _validUntilTime, bool _suspended, uint256[] calldata _values) external virtual returns (bool);
    function name() public view virtual returns (string memory);
    function currency() public view virtual returns (bytes32);
    function userCount() public view virtual returns (uint256);
    function userId(address _address) public view virtual returns (uint256);
    function validUserId(address _address) public view virtual returns (uint256);
    function validUser(address _address, uint256[] memory _keys) public view virtual returns (uint256, uint256[] memory);
    function validity(uint256 _userId) public view virtual returns (uint256, bool);
    function extendedKeys() public view virtual returns (uint256[] memory);
    function extended(uint256 _userId, uint256 _key) public view virtual returns (uint256);
    function manyExtended(uint256 _userId, uint256[] memory _key) public view virtual returns (uint256[] memory);
    function isAddressValid(address _address) public view virtual returns (bool);
    function isValid(uint256 _userId) public view virtual returns (bool);
    function defineExtendedKeys(uint256[] memory _extendedKeys) public virtual returns (bool);
    function registerUser(address _address, uint256 _validUntilTime) public virtual returns (bool);
    function registerUserFull(address _address, uint256 _validUntilTime, uint256[] memory _values) public virtual returns (bool);
    function attachAddress(uint256 _userId, address _address) public virtual returns (bool);
    function detachAddress(address _address) public virtual returns (bool);
    function detachSelf() public virtual returns (bool);
    function detachSelfAddress(address _address) public virtual returns (bool);
    function suspendUser(uint256 _userId) public virtual returns (bool);
    function restoreUser(uint256 _userId) public virtual returns (bool);
    function updateUser(uint256 _userId, uint256 _validUntilTime, bool _suspended) public virtual returns (bool);
    function updateUserExtended(uint256 _userId, uint256 _key, uint256 _value) public virtual returns (bool);
    function updateUserAllExtended(uint256 _userId, uint256[] memory _values) public virtual returns (bool);
    function updateUserFull(uint256 _userId, uint256 _validUntilTime, bool _suspended, uint256[] memory _values) public virtual returns (bool);
}
pragma solidity ^ 0.6.0;
abstract contract IRatesProvider {
    function defineRatesExternal(uint256[] calldata _rates) external virtual returns (bool);
    function name() public view virtual returns (string memory);
    function rate(bytes32 _currency) public view virtual returns (uint256);
    function currencies() public view virtual returns (bytes32[] memory, uint256[] memory, uint256);
    function rates() public view virtual returns (uint256, uint256[] memory);
    function convert(uint256 _amount, bytes32 _fromCurrency, bytes32 _toCurrency) public view virtual returns (uint256);
    function defineCurrencies(bytes32[] memory _currencies, uint256[] memory _decimals, uint256 _rateOffset) public virtual returns (bool);
    function defineRates(uint256[] memory _rates) public virtual returns (bool);
    event RateOffset(uint256 rateOffset);
    event Currencies(bytes32[] currencies, uint256[] decimals);
    event Rate(bytes32 currency, uint256 rate);
}
pragma solidity ^ 0.6.0;
abstract contract ITokenStorage {
    enum TransferCode {
        UNKNOWN,
        OK,
        INVALID_SENDER,
        NO_RECIPIENT,
        INSUFFICIENT_TOKENS,
        LOCKED,
        FROZEN,
        RULE,
        INVALID_RATE,
        NON_REGISTRED_SENDER,
        NON_REGISTRED_RECEIVER,
        LIMITED_EMISSION,
        LIMITED_RECEPTION
    }
    enum Scope {
        DEFAULT
    }
    enum AuditStorageMode {
        ADDRESS,
        USER_ID,
        SHARED
    }
    enum AuditMode {
        NEVER,
        ALWAYS,
        ALWAYS_TRIGGERS_EXCLUDED,
        WHEN_TRIGGERS_MATCHED
    }
    event OracleDefined(IUserRegistry userRegistry, IRatesProvider ratesProvider, address currency);
    event TokenDelegateDefined(uint256 delegateId, address delegate, uint256[] configurations);
    event TokenDelegateRemoved(uint256 delegateId);
    event AuditConfigurationDefined(uint256 configurationId, uint256 scopeId, AuditMode mode, uint256[] senderKeys, uint256[] receiverKeys, IRatesProvider ratesProvider, address currency);
    event AuditTriggersDefined(uint256 configurationId, address[] triggers, bool[] tokens, bool[] senders, bool[] receivers);
    event AuditsRemoved(address scope, uint256 scopeId);
    event SelfManaged(address holder, bool active);
    event Minted(address token, uint256 amount);
    event MintFinished(address token);
    event Burned(address token, uint256 amount);
    event RulesDefined(address token, IRule[] rules);
    event LockDefined(address token, uint256 startAt, uint256 endAt, address[] exceptions);
    event Seize(address token, address account, uint256 amount);
    event Freeze(address address_, uint256 until);
    event ClaimDefined(address token, address claim, uint256 claimAt);
    event TokenDefined(address token, uint256 delegateId, string name, string symbol, uint256 decimals);
    event TokenMigrated(address token, address newCore);
    event TokenRemoved(address token);
    event LogTransferData(address token, address caller, address sender, address receiver, uint256 senderId, uint256[] senderKeys, bool senderFetched, uint256 receiverId, uint256[] receiverKeys, bool receiverFetched, uint256 value, uint256 convertedValue);
    event LogTransferAuditData(uint256 auditConfigurationId, uint256 scopeId, address currency, IRatesProvider ratesProvider, bool senderAuditRequired, bool receiverAuditRequired);
    event LogAuditData(uint64 createdAt, uint64 lastTransactionAt, uint256 cumulatedEmission, uint256 cumulatedReception);
}
pragma solidity ^ 0.6.0;
contract TokenStorage is ITokenStorage, OperableStorage {
    using SafeMath for uint256;
    struct Lock{
        uint256 startAt;
        uint256 endAt;
        mapping (address => bool) exceptions;
        address[] exceptionsList;
    }
    struct TokenData{
        string name;
        string symbol;
        uint256 decimals;
        uint256 totalSupply;
        mapping (address => uint256) balances;
        mapping (address => mapping (address => uint256)) allowances;
        bool mintingFinished;
        uint256 allTimeMinted;
        uint256 allTimeBurned;
        uint256 allTimeSeized;
        mapping (address => uint256) frozenUntils;
        Lock lock;
        IRule[] rules;
    }
    struct AuditData{
        uint64 createdAt;
        uint64 lastTransactionAt;
        uint256 cumulatedEmission;
        uint256 cumulatedReception;
    }
    struct AuditStorage{
        address currency;
        AuditData sharedData;
        mapping (uint256 => AuditData) userData;
        mapping (address => AuditData) addressData;
    }
    struct AuditConfiguration{
        uint256 scopeId;
        AuditMode mode;
        uint256[] senderKeys;
        uint256[] receiverKeys;
        IRatesProvider ratesProvider;
        mapping (address => bool) triggerSenders;
        mapping (address => bool) triggerReceivers;
        mapping (address => bool) triggerTokens;
    }
    mapping (uint256 => AuditConfiguration) auditConfigurations;
    mapping (uint256 => uint256[]) delegatesConfigurations_;
    mapping (address => TokenData) tokens;
    mapping (address => mapping (uint256 => AuditStorage)) audits;
    mapping (address => bool) selfManaged;
    IUserRegistry userRegistry_;
    IRatesProvider ratesProvider_;
    address currency_;
    string name_;
    function currentTime() internal view returns (uint64) {
        return uint64(now);
    }
}
pragma solidity ^ 0.6.0;
abstract contract ITokenCore is ITokenStorage, IOperableCore {
    function name() public view virtual returns (string memory);
    function oracle() public view virtual returns (IUserRegistry userRegistry, IRatesProvider ratesProvider, address currency);
    function auditConfiguration(uint256 _configurationId) public view virtual returns (uint256 scopeId, AuditMode mode, uint256[] memory senderKeys, uint256[] memory receiverKeys, IRatesProvider ratesProvider, address currency);
    function auditTriggers(uint256 _configurationId, address[] memory _triggers) public view virtual returns (bool[] memory senders, bool[] memory receivers, bool[] memory tokens);
    function delegatesConfigurations(uint256 _delegateId) public view virtual returns (uint256[] memory);
    function auditCurrency(address _scope, uint256 _scopeId) external view virtual returns (address currency);
    function audit(address _scope, uint256 _scopeId, AuditStorageMode _storageMode, bytes32 _storageId) external view virtual returns (uint64 createdAt, uint64 lastTransactionAt, uint256 cumulatedEmission, uint256 cumulatedReception);
    function tokenName() external view virtual returns (string memory);
    function tokenSymbol() external view virtual returns (string memory);
    function decimals() external virtual returns (uint256);
    function totalSupply() external virtual returns (uint256);
    function balanceOf(address) external virtual returns (uint256);
    function allowance(address, address) external virtual returns (uint256);
    function transfer(address, address, uint256) external virtual returns (bool status);
    function transferFrom(address, address, address, uint256) external virtual returns (bool status);
    function approve(address, address, uint256) external virtual returns (bool status);
    function increaseApproval(address, address, uint256) external virtual returns (bool status);
    function decreaseApproval(address, address, uint256) external virtual returns (bool status);
    function token(address _token) external view virtual returns (bool mintingFinished, uint256 allTimeMinted, uint256 allTimeBurned, uint256 allTimeSeized, uint256[2] memory lock, address[] memory lockExceptions, uint256 freezedUntil, IRule[] memory);
    function canTransfer(address, address, uint256) external virtual returns (uint256);
    function mint(address, address[] calldata, uint256[] calldata) external virtual returns (bool);
    function finishMinting(address) external virtual returns (bool);
    function burn(address, uint256) external virtual returns (bool);
    function seize(address _token, address, uint256) external virtual returns (bool);
    function freezeManyAddresses(address _token, address[] calldata _addresses, uint256 _until) external virtual returns (bool);
    function defineLock(address, uint256, uint256, address[] calldata) external virtual returns (bool);
    function defineRules(address, IRule[] calldata) external virtual returns (bool);
    function defineToken(address _token, uint256 _delegateId, string memory _name, string memory _symbol, uint256 _decimals) external virtual returns (bool);
    function migrateToken(address _token, address _newCore) external virtual returns (bool);
    function removeToken(address _token) external virtual returns (bool);
    function defineOracle(IUserRegistry _userRegistry, IRatesProvider _ratesProvider, address _currency) external virtual returns (bool);
    function defineTokenDelegate(uint256 _delegateId, address _delegate, uint256[] calldata _configurations) external virtual returns (bool);
    function defineAuditConfiguration(uint256 _configurationId, uint256 _scopeId, AuditMode _mode, uint256[] calldata _senderKeys, uint256[] calldata _receiverKeys, IRatesProvider _ratesProvider, address _currency) external virtual returns (bool);
    function removeAudits(address _scope, uint256 _scopeId) external virtual returns (bool);
    function defineAuditTriggers(uint256 _configurationId, address[] calldata _triggerAddresses, bool[] calldata _triggerSenders, bool[] calldata _triggerReceivers, bool[] calldata _triggerTokens) external virtual returns (bool);
    function isSelfManaged(address _owner) external view virtual returns (bool);
    function manageSelf(bool _active) external virtual returns (bool);
}
pragma solidity ^ 0.6.0;
abstract contract ITokenDelegate is ITokenStorage {
    function decimals() public view virtual returns (uint256);
    function totalSupply() public view virtual returns (uint256);
    function balanceOf(address _owner) public view virtual returns (uint256);
    function allowance(address _owner, address _spender) public view virtual returns (uint256);
    function transfer(address _sender, address _receiver, uint256 _value) public virtual returns (bool);
    function transferFrom(address _caller, address _sender, address _receiver, uint256 _value) public virtual returns (bool);
    function canTransfer(address _sender, address _receiver, uint256 _value) public view virtual returns (TransferCode);
    function approve(address _sender, address _spender, uint256 _value) public virtual returns (bool);
    function increaseApproval(address _sender, address _spender, uint _addedValue) public virtual returns (bool);
    function decreaseApproval(address _sender, address _spender, uint _subtractedValue) public virtual returns (bool);
    function checkConfigurations(uint256[] memory _auditConfigurationIds) public virtual returns (bool);
}
pragma solidity ^ 0.6.0;
contract TokenCore is ITokenCore, OperableCore, TokenStorage {
    constructor(string memory _name, address[] memory _sysOperators) OperableCore(_sysOperators) public {
        name_ = _name;
    }
    function name() public view override returns (string memory) {
        return name_;
    }
    function oracle() public view override returns (IUserRegistry userRegistry, IRatesProvider ratesProvider, address currency) {
        return (userRegistry_, ratesProvider_, currency_);
    }
    function auditConfiguration(uint256 _configurationId) public view override returns (uint256 scopeId, AuditMode mode, uint256[] memory senderKeys, uint256[] memory receiverKeys, IRatesProvider ratesProvider, address currency) {
        AuditConfiguration storage auditConfiguration_ = auditConfigurations[_configurationId];
        return (auditConfiguration_.scopeId, auditConfiguration_.mode, auditConfiguration_.senderKeys, auditConfiguration_.receiverKeys, auditConfiguration_.ratesProvider, audits[address(this)][auditConfiguration_.scopeId].currency);
    }
    function auditTriggers(uint256 _configurationId, address[] memory _triggers) public view override returns (bool[] memory senders, bool[] memory receivers, bool[] memory tokens) {
        AuditConfiguration storage auditConfiguration_ = auditConfigurations[_configurationId];
        senders = new bool[](_triggers.length);
        receivers = new bool[](_triggers.length);
        tokens = new bool[](_triggers.length);
        for (uint256 i = 0; i < _triggers.length; i++) {
            senders[i] = auditConfiguration_.triggerSenders[_triggers[i]];
            receivers[i] = auditConfiguration_.triggerReceivers[_triggers[i]];
            tokens[i] = auditConfiguration_.triggerTokens[_triggers[i]];
        }
    }
    function delegatesConfigurations(uint256 _delegateId) public view override returns (uint256[] memory) {
        return delegatesConfigurations_[_delegateId];
    }
    function auditCurrency(address _scope, uint256 _scopeId) external view override returns (address currency) {
        return audits[_scope][_scopeId].currency;
    }
    function audit(address _scope, uint256 _scopeId, AuditStorageMode _storageMode, bytes32 _storageId) external view override returns (uint64 createdAt, uint64 lastTransactionAt, uint256 cumulatedEmission, uint256 cumulatedReception) {
        AuditData memory auditData;
        if(_storageMode == AuditStorageMode.SHARED) {
            auditData = audits[_scope][_scopeId].sharedData;
        }
        if(_storageMode == AuditStorageMode.ADDRESS) {
            auditData = audits[_scope][_scopeId].addressData[address(bytes20(_storageId))];
        }
        if(_storageMode == AuditStorageMode.USER_ID) {
            auditData = audits[_scope][_scopeId].userData[uint256(_storageId)];
        }
        createdAt = auditData.createdAt;
        lastTransactionAt = auditData.lastTransactionAt;
        cumulatedEmission = auditData.cumulatedEmission;
        cumulatedReception = auditData.cumulatedReception;
    }
    function tokenName() external view override returns (string memory) {
        return tokens[msg.sender].name;
    }
    function tokenSymbol() external view override returns (string memory) {
        return tokens[msg.sender].symbol;
    }
    function decimals() external onlyProxy override returns (uint256) {
        return delegateCallUint256(msg.sender);
    }
    function totalSupply() external onlyProxy override returns (uint256) {
        return delegateCallUint256(msg.sender);
    }
    function balanceOf(address) external onlyProxy override returns (uint256) {
        return delegateCallUint256(msg.sender);
    }
    function allowance(address, address) external onlyProxy override returns (uint256) {
        return delegateCallUint256(msg.sender);
    }
    function transfer(address, address, uint256) external onlyProxy override returns (bool status) {
        return delegateCall(msg.sender);
    }
    function transferFrom(address, address, address, uint256) external onlyProxy override returns (bool status) {
        return delegateCall(msg.sender);
    }
    function approve(address, address, uint256) external onlyProxy override returns (bool status) {
        return delegateCall(msg.sender);
    }
    function increaseApproval(address, address, uint256) external onlyProxy override returns (bool status) {
        return delegateCall(msg.sender);
    }
    function decreaseApproval(address, address, uint256) external onlyProxy override returns (bool status) {
        return delegateCall(msg.sender);
    }
    function token(address _token) external view override returns (bool mintingFinished, uint256 allTimeMinted, uint256 allTimeBurned, uint256 allTimeSeized, uint256[2] memory lock, address[] memory lockExceptions, uint256 frozenUntil, IRule[] memory rules) {
        TokenData storage tokenData = tokens[_token];
        mintingFinished = tokenData.mintingFinished;
        allTimeMinted = tokenData.allTimeMinted;
        allTimeBurned = tokenData.allTimeBurned;
        allTimeSeized = tokenData.allTimeSeized;
        lock = (tokenData.lock.startAt, tokenData.lock.endAt);
        lockExceptions = tokenData.lock.exceptionsList;
        frozenUntil = tokenData.frozenUntils[_token];
        rules = tokenData.rules;
    }
    function canTransfer(address, address, uint256) external onlyProxy override returns (uint256) {
        return delegateCallUint256(msg.sender);
    }
    function mint(address _token, address[] calldata, uint256[] calldata) external onlyProxyOp(_token) override returns (bool) {
        return delegateCall(_token);
    }
    function finishMinting(address _token) external onlyProxyOp(_token) override returns (bool) {
        return delegateCall(_token);
    }
    function burn(address _token, uint256) external onlyProxyOp(_token) override returns (bool) {
        return delegateCall(_token);
    }
    function seize(address _token, address, uint256) external onlyProxyOp(_token) override returns (bool) {
        return delegateCall(_token);
    }
    function freezeManyAddresses(address _token, address[] calldata, uint256) external onlyProxyOp(_token) override returns (bool) {
        return delegateCall(_token);
    }
    function defineLock(address _token, uint256, uint256, address[] calldata) external onlyProxyOp(_token) override returns (bool) {
        return delegateCall(_token);
    }
    function defineRules(address _token, IRule[] calldata) external onlyProxyOp(_token) override returns (bool) {
        return delegateCall(_token);
    }
    function defineToken(address _token, uint256 _delegateId, string calldata _name, string calldata _symbol, uint256 _decimals) external onlyCoreOp override returns (bool) {
        require(_token != ALL_PROXIES, "TC01");
        defineProxyInternal(_token, _delegateId);
        TokenData storage tokenData = tokens[_token];
        tokenData.name = _name;
        tokenData.symbol = _symbol;
        tokenData.decimals = _decimals;
        emit TokenDefined(_token, _delegateId, _name, _symbol, _decimals);
        return true;
    }
    function migrateToken(address _token, address _newCore) external onlyCoreOp override returns (bool) {
        migrateProxyInternal(_token, _newCore);
        emit TokenMigrated(_token, _newCore);
        return true;
    }
    function removeToken(address _token) external onlyCoreOp override returns (bool) {
        removeProxyInternal(_token);
        delete tokens[_token];
        emit TokenRemoved(_token);
        return true;
    }
    function defineOracle(IUserRegistry _userRegistry, IRatesProvider _ratesProvider, address _currency) external onlyCoreOp override returns (bool) {
        userRegistry_ = _userRegistry;
        ratesProvider_ = _ratesProvider;
        currency_ = _currency;
        emit OracleDefined(userRegistry_, _ratesProvider, _currency);
        return true;
    }
    function defineTokenDelegate(uint256 _delegateId, address _delegate, uint256[] calldata _auditConfigurations) external onlyCoreOp override returns (bool) {
        require(_delegate == address(0) || ITokenDelegate(_delegate).checkConfigurations(_auditConfigurations), "TC03");
        defineDelegateInternal(_delegateId, _delegate);
        if(_delegate != address(0)) {
            delegatesConfigurations_[_delegateId] = _auditConfigurations;
            emit TokenDelegateDefined(_delegateId, _delegate, _auditConfigurations);
        } else {
            delete delegatesConfigurations_[_delegateId];
            emit TokenDelegateRemoved(_delegateId);
        }
        return true;
    }
    function defineAuditConfiguration(uint256 _configurationId, uint256 _scopeId, AuditMode _mode, uint256[] calldata _senderKeys, uint256[] calldata _receiverKeys, IRatesProvider _ratesProvider, address _currency) external onlyCoreOp override returns (bool) {
        AuditStorage storage auditStorage = audits[address(this)][_scopeId];
        if(auditStorage.currency == address(0)) {
            auditStorage.currency = _currency;
        } else {
            require(auditStorage.currency == _currency, "TC04");
        }
        AuditConfiguration storage auditConfiguration_ = auditConfigurations[_configurationId];
        auditConfiguration_.mode = _mode;
        auditConfiguration_.scopeId = _scopeId;
        auditConfiguration_.senderKeys = _senderKeys;
        auditConfiguration_.receiverKeys = _receiverKeys;
        auditConfiguration_.ratesProvider = _ratesProvider;
        emit AuditConfigurationDefined(_configurationId, _scopeId, _mode, _senderKeys, _receiverKeys, _ratesProvider, _currency);
        return true;
    }
    function removeAudits(address _scope, uint256 _scopeId) external onlyCoreOp override returns (bool) {
        delete audits[_scope][_scopeId];
        emit AuditsRemoved(_scope, _scopeId);
        return true;
    }
    function defineAuditTriggers(uint256 _configurationId, address[] calldata _triggerAddresses, bool[] calldata _triggerTokens, bool[] calldata _triggerSenders, bool[] calldata _triggerReceivers) external onlyCoreOp override returns (bool) {
        require(_triggerAddresses.length == _triggerSenders.length && _triggerAddresses.length == _triggerReceivers.length && _triggerAddresses.length == _triggerTokens.length, "TC05");
        AuditConfiguration storage auditConfiguration_ = auditConfigurations[_configurationId];
        for (uint256 i = 0; i < _triggerAddresses.length; i++) {
            auditConfiguration_.triggerSenders[_triggerAddresses[i]] = _triggerSenders[i];
            auditConfiguration_.triggerReceivers[_triggerAddresses[i]] = _triggerReceivers[i];
            auditConfiguration_.triggerTokens[_triggerAddresses[i]] = _triggerTokens[i];
        }
        emit AuditTriggersDefined(_configurationId, _triggerAddresses, _triggerTokens, _triggerSenders, _triggerReceivers);
        return true;
    }
    function isSelfManaged(address _owner) external view override returns (bool) {
        return selfManaged[_owner];
    }
    function manageSelf(bool _active) external override returns (bool) {
        selfManaged[msg.sender] = _active;
        emit SelfManaged(msg.sender, _active);
    }
}
pragma solidity ^ 0.6.0;
contract TokenProxy is ITokenProxy, OperableProxy {
    constructor(address _core) OperableProxy(_core) public {

    }
    function name() public view override returns (string memory) {
        return TokenCore(core).tokenName();
    }
    function symbol() public view override returns (string memory) {
        return TokenCore(core).tokenSymbol();
    }
    function decimals() public view override returns (uint256) {
        return staticCallUint256();
    }
    function totalSupply() public view override returns (uint256) {
        return staticCallUint256();
    }
    function balanceOf(address) public view override returns (uint256) {
        return staticCallUint256();
    }
    function allowance(address, address) public view override returns (uint256) {
        return staticCallUint256();
    }
    function transfer(address _to, uint256 _value) public override returns (bool status) {
        return TokenCore(core).transfer(msg.sender, _to, _value);
    }
    function transferFrom(address _from, address _to, uint256 _value) public override returns (bool status) {
        return TokenCore(core).transferFrom(msg.sender, _from, _to, _value);
    }
    function approve(address _spender, uint256 _value) public override returns (bool status) {
        return TokenCore(core).approve(msg.sender, _spender, _value);
    }
    function increaseApproval(address _spender, uint256 _addedValue) public override returns (bool status) {
        return TokenCore(core).increaseApproval(msg.sender, _spender, _addedValue);
    }
    function decreaseApproval(address _spender, uint256 _subtractedValue) public override returns (bool status) {
        return TokenCore(core).decreaseApproval(msg.sender, _spender, _subtractedValue);
    }
    function canTransfer(address, address, uint256) public view override returns (uint256) {
        return staticCallUint256();
    }
    function emitTransfer(address _from, address _to, uint256 _value) public onlyCore override returns (bool) {
        emit Transfer(_from, _to, _value);
        return true;
    }
    function emitApproval(address _owner, address _spender, uint256 _value) public onlyCore override returns (bool) {
        emit Approval(_owner, _spender, _value);
        return true;
    }
}
