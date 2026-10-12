// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {StockTokenWrapper, IWrappedStock} from "../src/StockTokenWrapper.sol";

contract StockDouble is ERC20("Tesla", "TSLA") {
    uint256 public uiMultiplier = 1e18;
    uint256 public newUIMultiplier = 1e18;
    uint256 public effectiveAt;
    bool public tokenPaused;
    bool public paused;
    address public ACCESS_CONTROLLED_REGISTRY = address(0xBEEF);

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }

    function schedule(uint256 multiplier, uint256 when) external {
        newUIMultiplier = multiplier;
        effectiveAt = when;
    }

    function setTokenPaused(bool value) external {
        tokenPaused = value;
    }

    function setPaused(bool value) external {
        paused = value;
    }
}

contract StockTokenWrapperTest is Test {
    StockDouble internal stock;
    StockTokenWrapper internal wrapper;
    address internal alice = makeAddr("alice");

    function setUp() public {
        stock = new StockDouble();
        wrapper = new StockTokenWrapper(IWrappedStock(address(stock)));
        stock.mint(alice, 10e18);
    }

    function test_NamesTheWrappedToken() public view {
        assertEq(wrapper.name(), "Wrapped Tesla");
        assertEq(wrapper.symbol(), "wTSLA");
        assertEq(wrapper.decimals(), 18);
        assertEq(address(wrapper.underlying()), address(stock));
    }

    function test_WrapsAndUnwrapsOneToOne() public {
        vm.startPrank(alice);
        stock.approve(address(wrapper), 4e18);
        wrapper.depositFor(alice, 4e18);
        assertEq(wrapper.balanceOf(alice), 4e18);
        assertEq(stock.balanceOf(address(wrapper)), 4e18);
        wrapper.withdrawTo(alice, 1.5e18);
        vm.stopPrank();
        assertEq(wrapper.balanceOf(alice), 2.5e18);
        assertEq(stock.balanceOf(alice), 7.5e18);
        assertEq(wrapper.totalSupply(), stock.balanceOf(address(wrapper)));
    }

    function test_ForwardsTheStockTokensViews() public {
        assertEq(wrapper.uiMultiplier(), 1e18);
        assertEq(wrapper.newUIMultiplier(), 1e18);
        assertEq(wrapper.effectiveAt(), 0);
        assertEq(wrapper.ACCESS_CONTROLLED_REGISTRY(), address(0xBEEF));
        assertFalse(wrapper.oraclePaused());
        assertFalse(wrapper.paused());

        stock.schedule(2e18, 1_790_000_000);
        stock.setTokenPaused(true);
        stock.setPaused(true);
        assertEq(wrapper.newUIMultiplier(), 2e18);
        assertEq(wrapper.effectiveAt(), 1_790_000_000);
        assertTrue(wrapper.oraclePaused());
        assertTrue(wrapper.paused());
    }

    function test_PermitSpendsWithoutAnApproval() public {
        (address owner, uint256 key) = makeAddrAndKey("owner");
        stock.mint(owner, 3e18);
        vm.startPrank(owner);
        stock.approve(address(wrapper), 3e18);
        wrapper.depositFor(owner, 3e18);
        vm.stopPrank();

        uint256 deadline = block.timestamp + 1 hours;
        bytes32 digest = keccak256(
            abi.encodePacked(
                "\x19\x01",
                wrapper.DOMAIN_SEPARATOR(),
                keccak256(
                    abi.encode(
                        keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"),
                        owner,
                        alice,
                        2e18,
                        wrapper.nonces(owner),
                        deadline
                    )
                )
            )
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, digest);
        wrapper.permit(owner, alice, 2e18, deadline, v, r, s);
        vm.prank(alice);
        IERC20(address(wrapper)).transferFrom(owner, alice, 2e18);
        assertEq(wrapper.balanceOf(alice), 2e18);
    }

    function test_CannotWithdrawMoreThanWrapped() public {
        vm.startPrank(alice);
        stock.approve(address(wrapper), 1e18);
        wrapper.depositFor(alice, 1e18);
        vm.expectRevert();
        wrapper.withdrawTo(alice, 2e18);
        vm.stopPrank();
    }
}
