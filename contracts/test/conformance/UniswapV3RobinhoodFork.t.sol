// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IUniswapV3Factory, IUniswapV3Pool} from "./Interfaces.sol";

contract UniswapV3RobinhoodForkTest is Test {
    uint256 internal constant WEEKDAY_BLOCK = 69_922_505;
    uint32 internal constant TWAP_WINDOW = 30 minutes;

    string internal robinhood = vm.readFile("../deployments/4663.json");

    function test_EveryPoolIsTheFactoryPoolForItsKey() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        IUniswapV3Factory factory = IUniswapV3Factory(vm.parseJsonAddress(robinhood, ".uniswapV3.Factory"));
        string[] memory keys = _poolKeys();
        for (uint256 i; i < keys.length; ++i) {
            string[] memory parts = vm.split(keys[i], "_");
            assertTrue(_isQuote(parts[1]) && !_isQuote(parts[0]), keys[i]);
            address base = vm.parseJsonAddress(robinhood, string.concat(".tokens.", parts[0]));
            address quote = vm.parseJsonAddress(robinhood, string.concat(".tokens.", parts[1]));
            uint24 fee = uint24(vm.parseUint(parts[2]));
            IUniswapV3Pool pool = _pool(keys[i]);

            assertEq(factory.getPool(base, quote, fee), address(pool), keys[i]);
            (address token0, address token1) = base < quote ? (base, quote) : (quote, base);
            assertEq(pool.token0(), token0, keys[i]);
            assertEq(pool.token1(), token1, keys[i]);
            assertEq(pool.fee(), fee, keys[i]);
            assertGt(pool.liquidity(), 0, keys[i]);
        }
    }

    function test_EveryPoolServesAThirtyMinuteTwap() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        string[] memory keys = _poolKeys();
        uint32[] memory secondsAgos = new uint32[](2);
        secondsAgos[0] = TWAP_WINDOW;
        for (uint256 i; i < keys.length; ++i) {
            IUniswapV3Pool pool = _pool(keys[i]);
            assertGe(_longestWindow(pool), TWAP_WINDOW, keys[i]);
            pool.observe(secondsAgos);
        }
    }

    function test_ObserveRevertsPastTheOldestObservation() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        IUniswapV3Pool pool = _pool("NVDA_USDG_500");
        uint32 window = _longestWindow(pool);
        uint32[] memory secondsAgos = new uint32[](2);
        secondsAgos[0] = window;
        pool.observe(secondsAgos);
        secondsAgos[0] = window + 1;
        vm.expectRevert(bytes("OLD"));
        pool.observe(secondsAgos);
    }

    function test_LongestWindowsAtThePinnedBlock() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        assertEq(_longestWindow(_pool("NVDA_USDG_500")), 343_965);
        assertEq(_longestWindow(_pool("TSLA_USDG_3000")), 460_012);
        assertEq(_longestWindow(_pool("AAPL_USDG_500")), 109_642);
        assertEq(_longestWindow(_pool("MSFT_USDG_3000")), 642_633);
        assertEq(_longestWindow(_pool("GOOGL_USDG_500")), 83_711);
        assertEq(_longestWindow(_pool("SPY_WETH_500")), 102_248);
        assertEq(_longestWindow(_pool("SPY_USDG_500")), 428_107);
    }

    function _longestWindow(IUniswapV3Pool pool) internal view returns (uint32) {
        (,, uint16 index, uint16 cardinality,,,) = pool.slot0();
        (uint32 oldest,,, bool initialized) = pool.observations((index + 1) % cardinality);
        if (!initialized) (oldest,,,) = pool.observations(0);
        return uint32(block.timestamp) - oldest;
    }

    function _isQuote(string memory symbol) internal pure returns (bool) {
        return keccak256(bytes(symbol)) == keccak256("USDG") || keccak256(bytes(symbol)) == keccak256("WETH");
    }

    function _poolKeys() internal view returns (string[] memory keys) {
        string[] memory names = vm.parseJsonKeys(robinhood, ".uniswapV3");
        keys = new string[](names.length - 1);
        uint256 n;
        for (uint256 i; i < names.length; ++i) {
            if (keccak256(bytes(names[i])) != keccak256("Factory")) keys[n++] = names[i];
        }
    }

    function _pool(string memory key) internal view returns (IUniswapV3Pool) {
        return IUniswapV3Pool(vm.parseJsonAddress(robinhood, string.concat(".uniswapV3.", key)));
    }
}
