// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IERC20} from "forge-std/interfaces/IERC20.sol";

interface AggregatorV3Interface {
    function decimals() external view returns (uint8);

    function description() external view returns (string memory);

    function version() external view returns (uint256);

    function getRoundData(uint80 roundId_)
        external
        view
        returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);

    function latestRoundData()
        external
        view
        returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);
}

interface AggregatorProxy is AggregatorV3Interface {
    function aggregator() external view returns (address);

    function phaseId() external view returns (uint16);
}

interface IUniswapV3Factory {
    function feeAmountTickSpacing(uint24 fee) external view returns (int24);

    function getPool(address tokenA, address tokenB, uint24 fee) external view returns (address);
}

interface IUniswapV3Pool {
    function token0() external view returns (address);

    function token1() external view returns (address);

    function fee() external view returns (uint24);

    function liquidity() external view returns (uint128);

    function slot0()
        external
        view
        returns (
            uint160 sqrtPriceX96,
            int24 tick,
            uint16 observationIndex,
            uint16 observationCardinality,
            uint16 observationCardinalityNext,
            uint8 feeProtocol,
            bool unlocked
        );

    function observations(uint256 index)
        external
        view
        returns (
            uint32 blockTimestamp,
            int56 tickCumulative,
            uint160 secondsPerLiquidityCumulativeX128,
            bool initialized
        );

    function observe(uint32[] calldata secondsAgos)
        external
        view
        returns (int56[] memory tickCumulatives, uint160[] memory secondsPerLiquidityCumulativeX128s);
}

interface ArbSys {
    function arbBlockNumber() external view returns (uint256);

    function arbOSVersion() external view returns (uint256);
}

interface ArbGasInfo {
    function getPricesInWei() external view returns (uint256, uint256, uint256, uint256, uint256, uint256);
}

interface IStockToken is IERC20 {
    event TransferWithScaledUI(address indexed from, address indexed to, uint256 value, uint256 uiValue);
    event UIMultiplierUpdated(uint256 oldMultiplier, uint256 newMultiplier, uint256 effectiveAtTimestamp);
    event Paused();
    event Unpaused();
    event OraclePaused();
    event OracleUnpaused();

    error IsPaused();
    error Blocked(address account);
    error AccessControlUnauthorizedAccount(address account, bytes32 neededRole);

    function ACCESS_CONTROLLED_REGISTRY() external view returns (address);
    function DOMAIN_SEPARATOR() external view returns (bytes32);
    function nonces(address owner) external view returns (uint256);
    function eip712Domain()
        external
        view
        returns (
            bytes1 fields,
            string memory name,
            string memory version,
            uint256 chainId,
            address verifyingContract,
            bytes32 salt,
            uint256[] memory extensions
        );
    function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external;
    function paused() external view returns (bool);
    function tokenPaused() external view returns (bool);
    function pause() external;
    function unpause() external;
    function adminBurn(address from, uint256 amount) external;
    function uiMultiplier() external view returns (uint256);
    function newUIMultiplier() external view returns (uint256);
    function effectiveAt() external view returns (uint256);
    function balanceOfUI(address account) external view returns (uint256);
    function totalSupplyUI() external view returns (uint256);
    function supportsInterface(bytes4 interfaceId) external view returns (bool);
    function updateMultiplier(uint256 newMultiplier, uint256 effectiveAt_) external;
    function oraclePaused() external view returns (bool);
    function pauseOracle() external;
    function unpauseOracle() external;
}

interface IStockTokenRegistry {
    event Blocked(address indexed account);
    event Unblocked(address indexed account);
    event Paused();
    event Unpaused();
    event Upgraded(address indexed implementation);

    function implementation() external view returns (address);
    function hasRole(bytes32 role, address account) external view returns (bool);
    function getRoleAdmin(bytes32 role) external view returns (bytes32);
    function isBlocked(address account) external view returns (bool);
    function blockAccounts(address[] calldata accounts) external;
    function unblockAccounts(address[] calldata accounts) external;
    function paused() external view returns (bool);
    function pause() external;
    function unpause() external;
}

interface IUSDG is IERC20 {
    error ContractPaused();
    error AddressFrozen();
    error FacetNotFound();
    error AddressNotFrozen();

    function owner() external view returns (address);
    function defaultAdmin() external view returns (address);
    function defaultAdminDelay() external view returns (uint48);
    function hasRole(bytes32 role, address account) external view returns (bool);
    function DOMAIN_SEPARATOR() external view returns (bytes32);
    function nonces(address owner) external view returns (uint256);
    function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external;
    function paused() external view returns (bool);
    function pause() external;
    function unpause() external;
    function isFrozen(address addr) external view returns (bool);
    function freeze(address addr) external;
    function wipeFrozenAddress(address addr) external;
    function upgradeToAndCall(address newImplementation, bytes calldata data) external payable;
}

interface ITimelockController {
    function getMinDelay() external view returns (uint256);
    function hasRole(bytes32 role, address account) external view returns (bool);
}

interface IMorpho {
    struct MarketParams {
        address loanToken;
        address collateralToken;
        address oracle;
        address irm;
        uint256 lltv;
    }

    function isIrmEnabled(address irm) external view returns (bool);
    function isLltvEnabled(uint256 lltv) external view returns (bool);
    function createMarket(MarketParams memory marketParams) external;
    function idToMarketParams(bytes32 id) external view returns (MarketParams memory);
    function market(bytes32 id)
        external
        view
        returns (
            uint128 totalSupplyAssets,
            uint128 totalSupplyShares,
            uint128 totalBorrowAssets,
            uint128 totalBorrowShares,
            uint128 lastUpdate,
            uint128 fee
        );
    function position(bytes32 id, address user)
        external
        view
        returns (uint256 supplyShares, uint128 borrowShares, uint128 collateral);
    function supply(
        MarketParams memory marketParams,
        uint256 assets,
        uint256 shares,
        address onBehalf,
        bytes memory data
    ) external returns (uint256, uint256);
    function withdraw(
        MarketParams memory marketParams,
        uint256 assets,
        uint256 shares,
        address onBehalf,
        address receiver
    ) external returns (uint256, uint256);
    function borrow(
        MarketParams memory marketParams,
        uint256 assets,
        uint256 shares,
        address onBehalf,
        address receiver
    ) external returns (uint256, uint256);
    function repay(
        MarketParams memory marketParams,
        uint256 assets,
        uint256 shares,
        address onBehalf,
        bytes memory data
    ) external returns (uint256, uint256);
    function supplyCollateral(MarketParams memory marketParams, uint256 assets, address onBehalf, bytes memory data)
        external;
    function withdrawCollateral(MarketParams memory marketParams, uint256 assets, address onBehalf, address receiver)
        external;
    function liquidate(
        MarketParams memory marketParams,
        address borrower,
        uint256 seizedAssets,
        uint256 repaidShares,
        bytes memory data
    ) external returns (uint256, uint256);
    function accrueInterest(MarketParams memory marketParams) external;
}
