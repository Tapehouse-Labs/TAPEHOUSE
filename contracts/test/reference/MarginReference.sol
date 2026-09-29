// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {Pools} from "./Pools.sol";
import {Pool, Requirements} from "./Requirements.sol";
import {Parameters, Scenarios, Set} from "./Scenarios.sol";
import {Session} from "./Session.sol";

interface IDecimals {
    function decimals() external view returns (uint8);
}

/// @notice The margin engine in Solidity: the Stylus program's views, bit for bit, over the same configuration.
/// It holds that configuration as given, without the program's checks, its parameter updates or its ownership
/// transfers.
contract MarginReference {
    struct Asset {
        bytes32 symbol;
        uint32 volatilityFloor;
        uint32 volatility;
        uint32 gapFloor;
        uint32 gap;
        uint32 sellingDepth;
        uint32 buyingDepth;
        address pool;
        address token;
    }

    struct Source {
        address pool;
        bool stockIsToken0;
        bool quoteIsWeth;
        uint8 stockDecimals;
        uint8 quoteDecimals;
        uint32 fee;
    }

    uint256 internal constant SCENARIO_SIZE = 256;

    bytes32[] internal _assets;
    mapping(bytes32 => uint256) internal _positions;
    uint32[] internal _volatilities;
    uint32[] internal _volatilityFloors;
    uint16[] internal _correlations;
    uint16[] internal _correlationFloors;
    uint32[] internal _gaps;
    uint32[] internal _gapFloors;
    uint32[] internal _depths;
    uint256 internal _market;
    mapping(bytes32 => Source) internal _sources;
    /// @notice The Chainlink ETH/USD feed that prices pools quoted in WETH.
    address public immutable ethUsdFeed;
    /// @notice The band whose session sets the regime.
    address public immutable band;
    /// @notice The owner the configuration was given for.
    address public immutable owner;

    /// @notice `symbol` is not one of the assets.
    error UnknownAsset(bytes32 symbol);
    /// @notice The arrays do not hold one entry per asset or per pair.
    error LengthMismatch();
    /// @notice `size` is not a lattice size of the set: 32, 64, 128 or 256.
    error UnsupportedScenarioSize(uint16 size);
    /// @notice `index` lies past the scenario set.
    error ScenarioOutOfRange(uint16 index);
    /// @notice `pool` does not trade `symbol`'s Stock Token against USDG or WETH with 30 minutes of history.
    error InvalidPool(bytes32 symbol, address pool);
    /// @notice `feed` does not report with 8 decimals.
    error InvalidFeed(address feed);
    /// @notice A position in `symbol` is worth more than the engine margins.
    error ExposureTooLarge(bytes32 symbol);
    /// @notice The band's session call ran out of the gas this call left it.
    error InsufficientGas();

    constructor(
        Asset[] memory assets_,
        uint16[] memory correlationFloors,
        uint16[] memory correlations,
        bytes32 market_,
        address usdg,
        address weth,
        address ethUsd,
        address band_,
        address initialOwner
    ) {
        for (uint256 i; i < assets_.length; ++i) {
            Asset memory asset = assets_[i];
            _assets.push(asset.symbol);
            _positions[asset.symbol] = i + 1;
            _volatilities.push(asset.volatility);
            _volatilityFloors.push(asset.volatilityFloor);
            _gaps.push(asset.gap);
            _gapFloors.push(asset.gapFloor);
            _depths.push(asset.sellingDepth);
            _depths.push(asset.buyingDepth);
            if (asset.pool != address(0)) _sources[asset.symbol] = _source(asset, usdg, weth, ethUsd);
        }
        _correlations = correlations;
        _correlationFloors = correlationFloors;
        if (market_ != bytes32(0)) _market = _position(market_) + 1;
        ethUsdFeed = ethUsd;
        band = band_;
        owner = initialOwner;
    }

    /// @notice The assets, in the order of every parameter list.
    function assets() external view returns (bytes32[] memory) {
        return _assets;
    }

    /// @notice The daily volatility of `symbol` and its floor, in centi-basis-points.
    function volatility(bytes32 symbol) external view returns (uint32, uint32) {
        uint256 i = _position(symbol);
        return (_volatilities[i], _volatilityFloors[i]);
    }

    /// @notice The correlation of `symbol` and `other` and its floor, in basis points.
    function correlation(bytes32 symbol, bytes32 other) external view returns (uint16, uint16) {
        (uint256 i, uint256 j) = (_position(symbol), _position(other));
        if (i == j) return (10_000, 10_000);
        uint256 k = Scenarios.pair(_assets.length, i < j ? i : j, i < j ? j : i);
        return (_correlations[k], _correlationFloors[k]);
    }

    /// @notice The weekend gap of `symbol` and its floor, in centi-basis-points.
    function weekendGap(bytes32 symbol) external view returns (uint32, uint32) {
        uint256 i = _position(symbol);
        return (_gaps[i], _gapFloors[i]);
    }

    /// @notice The USD a liquidation of `symbol` can sell, then buy, within a 10% move of its pool, and the
    /// ceilings over both, which are the first values.
    function depth(bytes32 symbol) external view returns (uint32, uint32, uint32, uint32) {
        uint256 i = _position(symbol);
        return (_depths[2 * i], _depths[2 * i + 1], _depths[2 * i], _depths[2 * i + 1]);
    }

    /// @notice When the parameters were last updated: never, since the reference takes no updates.
    function lastUpdate() external pure returns (uint64) {
        return 0;
    }

    /// @notice The account an ownership transfer awaits; none, since the reference takes no transfers.
    function pendingOwner() external pure returns (address) {
        return address(0);
    }

    /// @notice The Uniswap v3 pool `symbol` is liquidated in; zero for none.
    function pool(bytes32 symbol) external view returns (address) {
        _position(symbol);
        return _sources[symbol].pool;
    }

    /// @notice The asset that stands for the market; zero for the equal-weighted portfolio of the assets.
    function market() external view returns (bytes32) {
        return _market == 0 ? bytes32(0) : _assets[_market - 1];
    }

    /// @notice The most gross exposure the current requirement allows per unit of margin across a closure, in
    /// basis points.
    function weekendLeverage() external pure returns (uint32) {
        return uint32(Requirements.WEEKEND_LEVERAGE);
    }

    /// @notice Scenario `index` of the engine's set over a horizon of `horizon` seconds.
    function scenario(uint16 index, uint64 horizon) external view returns (int32[] memory returns_) {
        uint256 n = _assets.length;
        if (index >= Scenarios.count(SCENARIO_SIZE, n)) revert ScenarioOutOfRange(index);
        Set memory set = Scenarios.create(_parameters(), SCENARIO_SIZE, horizon);
        int256[] memory row = new int256[](n);
        Scenarios.fill(set, index, row, new int256[](n));
        returns_ = new int32[](n);
        for (uint256 i; i < n; ++i) {
            returns_[i] = int32(row[i]);
        }
    }

    /// @notice Keccak-256 of every scenario of a set with `size` lattice points over `horizon` seconds.
    function scenarioDigest(uint16 size, uint64 horizon) external view returns (bytes32) {
        if (size != 32 && size != 64 && size != 128 && size != 256) revert UnsupportedScenarioSize(size);
        return keccak256(Scenarios.encoded(Scenarios.create(_parameters(), size, horizon)));
    }

    /// @notice The margin a portfolio needs over `horizon` seconds, and a bit for every asset whose liquidity
    /// input is missing.
    function requirement(int256[] calldata quantities, uint256[] calldata prices, uint64 horizon, bool spansClosure)
        external
        view
        returns (uint256, uint8)
    {
        (uint256 open, uint256 closed,, uint8 missing) = _requirements(quantities, prices, horizon);
        return (spansClosure ? closed : open, missing);
    }

    /// @notice The margin a portfolio needs now, the bits of `requirement`, and the regime the band's session
    /// puts it in: 0 unknown, 1 closed, 2 open, 3 closing.
    function currentRequirement(int256[] calldata quantities, uint256[] calldata prices)
        external
        view
        returns (uint256, uint8, uint8)
    {
        (uint256 open, uint256 closed, uint256 floor, uint8 missing) =
            _requirements(quantities, prices, Session.HORIZON);
        uint256 before = gasleft();
        (bool known, uint256 state,,,, uint256 boundaryMs) = Session.read(band);
        if (!known && gasleft() * 64 <= before) revert InsufficientGas();
        (uint8 code, uint256 elapsedMs) = Session.regime(known, state, boundaryMs, block.timestamp * 1000);
        return (Session.current(open, closed, floor, code, elapsedMs), missing, code);
    }

    function _requirements(int256[] calldata quantities, uint256[] calldata prices, uint64 horizon)
        internal
        view
        returns (uint256 open, uint256 closed, uint256 floor, uint8 missing)
    {
        uint256 n = _assets.length;
        if (quantities.length != n || prices.length != n) revert LengthMismatch();
        (int256[] memory values, uint256 bad) = Requirements.exposures(quantities, prices);
        if (bad != 0) revert ExposureTooLarge(_assets[bad - 1]);
        bool held;
        for (uint256 i; i < n; ++i) {
            if (values[i] != 0) held = true;
        }
        if (!held) return (0, 0, 0, 0);
        Set memory set = Scenarios.create(_parameters(), SCENARIO_SIZE, horizon);
        return Requirements.requirements(set, values, prices, _depths, _pools(values));
    }

    function _pools(int256[] memory values) internal view returns (Pool[] memory pools) {
        pools = new Pool[](values.length);
        bool ethRead;
        bool ethOk;
        uint256 ethUsd;
        for (uint256 i; i < values.length; ++i) {
            if (values[i] == 0) continue;
            Source memory source = _sources[_assets[i]];
            if (source.pool == address(0)) continue;
            pools[i].kind = Requirements.UNREAD;
            pools[i].fee = source.fee;
            (bool ok, int24 tick, uint128 liquidity) = Pools.consult(source.pool);
            if (!ok) continue;
            uint256 usd = 1e8;
            if (source.quoteIsWeth) {
                if (!ethRead) {
                    (ethOk, ethUsd) = Pools.answer(ethUsdFeed, block.timestamp);
                    ethRead = true;
                }
                if (!ethOk) continue;
                usd = ethUsd;
            }
            (bool fits, uint256 price, uint256 selling, uint256 buying) =
                Pools.terms(source.stockIsToken0, source.stockDecimals, source.quoteDecimals, tick, liquidity, usd);
            if (!fits) continue;
            pools[i] = Pool(Requirements.READ, price, selling, buying, source.fee);
        }
    }

    function _source(Asset memory asset, address usdg, address weth, address ethUsd)
        internal
        view
        returns (Source memory source)
    {
        IUniswapV3Pool uniswap = IUniswapV3Pool(asset.pool);
        (address token0, address token1, uint24 fee) = (uniswap.token0(), uniswap.token1(), uniswap.fee());
        (,,,, uint16 cardinality,,) = uniswap.slot0();
        bool stockIsToken0 = token0 == asset.token;
        address quote = stockIsToken0 ? token1 : token0;
        if (
            cardinality <= Pools.WINDOW || (token0 == asset.token) == (token1 == asset.token) || quote == address(0)
                || (quote != usdg && quote != weth)
        ) {
            revert InvalidPool(asset.symbol, asset.pool);
        }
        uint8 stockDecimals = IDecimals(asset.token).decimals();
        uint8 quoteDecimals = IDecimals(quote).decimals();
        if (stockDecimals > 36 || quoteDecimals > 18) revert InvalidPool(asset.symbol, asset.pool);
        if (quote == weth && (ethUsd == address(0) || IDecimals(ethUsd).decimals() != 8)) revert InvalidFeed(ethUsd);
        return Source(asset.pool, stockIsToken0, quote == weth, stockDecimals, quoteDecimals, fee);
    }

    function _parameters() internal view returns (Parameters memory) {
        return Parameters(_assets, _volatilities, _correlations, _gaps, _market);
    }

    function _position(bytes32 symbol) internal view returns (uint256) {
        uint256 p = _positions[symbol];
        if (p == 0) revert UnknownAsset(symbol);
        return p - 1;
    }
}
