// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test, stdError} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {MorphoBandOracle} from "../src/MorphoBandOracle.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {TokenDouble} from "./doubles/PoolDouble.sol";

contract MorphoBandOracleTest is Test {
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant NVDA_24_7 = "NVDA---24_7";
    bytes32 internal constant TSLA = "TSLA";
    uint64 internal constant LOW = 17_901_000_000;
    uint256 internal constant LOW_IN_MORPHO = 179_010_000_000_000_000_000_000_000;

    BandDouble internal band;
    TokenDouble internal nvda;
    TokenDouble internal usdg;
    MorphoBandOracle internal oracle;
    address internal owner = makeAddr("owner");

    function setUp() public {
        vm.warp(1_790_390_000);
        band = new BandDouble();
        nvda = new TokenDouble(18);
        usdg = new TokenDouble(6);
        band.setAsset(NVDA, BandDouble.Asset(address(0xC1), NVDA_24_7, bytes32(0), address(nvda)));
        band.setQuote(NVDA, BandDouble.Quote(3, 2, 18_000_000_000, 55, LOW, 18_099_000_000));
        oracle = _oracle(band);
    }

    function _oracle(BandDouble band_) internal returns (MorphoBandOracle) {
        return new MorphoBandOracle(IBand(address(band_)), NVDA, IERC20Metadata(address(usdg)), owner);
    }

    function _band(address token) internal returns (BandDouble other) {
        other = new BandDouble();
        other.setAsset(NVDA, BandDouble.Asset(address(0xC1), NVDA_24_7, bytes32(0), token));
        other.setQuote(NVDA, BandDouble.Quote(2, 1, 18_100_000_000, 80, 17_955_200_000, 18_244_800_000));
    }

    function test_ConstructorNeedsABandThatNamesTheStockToken() public {
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoStockToken.selector, address(band), bytes32("AAPL")));
        new MorphoBandOracle(IBand(address(band)), "AAPL", IERC20Metadata(address(usdg)), owner);
        band.setAsset(TSLA, BandDouble.Asset(address(0xC2), "TSLA---24_7", bytes32(0), address(0)));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoStockToken.selector, address(band), TSLA));
        new MorphoBandOracle(IBand(address(band)), TSLA, IERC20Metadata(address(usdg)), owner);
    }

    function test_TheOracleNamesItsMarketsTokensAndItsFirstBand() public {
        vm.expectEmit();
        emit MorphoBandOracle.BandSet(IBand(address(0)), IBand(address(band)));
        oracle = _oracle(band);
        assertEq(address(oracle.band()), address(band));
        assertEq(oracle.symbol(), NVDA);
        assertEq(oracle.collateralToken(), address(nvda));
        assertEq(oracle.loanToken(), address(usdg));
        assertEq(oracle.BAND_DECIMALS(), 8);
        assertEq(oracle.scaleFactor(), 1e16);
        assertEq(oracle.owner(), owner);
    }

    function test_EveryStateWithAPriceAnswersTheLowEdgeInMorphosScale() public {
        for (uint8 state = 1; state <= 3; ++state) {
            band.setQuote(NVDA, BandDouble.Quote(state, 1, 18_000_000_000 + state, 55, LOW, 18_099_000_000));
            assertEq(oracle.price(), LOW_IN_MORPHO);
        }
        uint256 collateral = 10e18;
        uint256 maxBorrow = collateral * oracle.price() / 1e36 * 0.625e18 / 1e18;
        assertEq(maxBorrow, 1_118_812_500);
    }

    function test_ADegradedQuoteAnswersTheLowEdgeBesideALapsedHalt() public {
        band.setQuote(NVDA, BandDouble.Quote(1, 2, 18_000_000_000, 55, LOW, 18_099_000_000));
        band.setHalt(NVDA, BandDouble.Halt(false, 1_790_389_000, 1_790_388_100, false));
        assertEq(oracle.price(), LOW_IN_MORPHO);
        (bool signedHalt, uint64 until,,) = oracle.halt();
        assertFalse(signedHalt);
        assertLt(until, vm.getBlockTimestamp());
    }

    function test_AZeroedQuoteHasNoAnswerWhateverHaltTheBandReports() public {
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        band.setHalt(NVDA, BandDouble.Halt(true, 1_790_390_600, 1_790_389_990, false));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoAnswer.selector, NVDA));
        oracle.price();
        band.setHalt(NVDA, BandDouble.Halt(false, 0, 0, true));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoAnswer.selector, NVDA));
        oracle.price();
        band.setHalt(NVDA, BandDouble.Halt(false, 0, 0, false));
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(2, 1_790_389_000, 1e18, 4e18));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoAnswer.selector, NVDA));
        oracle.price();
    }

    function test_AZeroLowEdgeHasNoAnswer() public {
        band.setQuote(NVDA, BandDouble.Quote(3, 2, 1, 1500, 0, 1));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoAnswer.selector, NVDA));
        oracle.price();
    }

    function test_AHaltedStateHasNoAnswerWhateverBoundsComeWithIt() public {
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 18_000_000_000, 55, LOW, 18_099_000_000));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoAnswer.selector, NVDA));
        oracle.price();
    }

    function test_AnUnsettledSequencerHasNoAnswer() public {
        band.setSequencer(address(0x5E9), false);
        oracle = _oracle(band);
        vm.expectRevert(MorphoBandOracle.SequencerNotSettled.selector);
        oracle.price();
        band.setSequencer(address(0x5E9), true);
        assertEq(oracle.price(), LOW_IN_MORPHO);
    }

    function test_AHaltedBandIsNoAnswerBeforeTheSequencerIsRead() public {
        band.setSequencer(address(0x5E9), false);
        oracle = _oracle(band);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.NoAnswer.selector, NVDA));
        oracle.price();
    }

    function test_WithoutASequencerFeedTheSequencerIsNeverRead() public {
        band.setSequencer(address(0), false);
        assertEq(oracle.price(), LOW_IN_MORPHO);
    }

    function test_TheScaleFollowsBothTokensDecimals() public {
        uint8[2][4] memory decimals = [[uint8(18), 18], [uint8(6), 6], [uint8(0), 28], [uint8(18), 6]];
        uint256[4] memory scales = [uint256(1e28), 1e28, 1, 1e40];
        for (uint256 i; i < 4; ++i) {
            BandDouble other = _band(address(new TokenDouble(decimals[i][1])));
            MorphoBandOracle scaled = new MorphoBandOracle(
                IBand(address(other)), NVDA, IERC20Metadata(address(new TokenDouble(decimals[i][0]))), owner
            );
            assertEq(scaled.scaleFactor(), scales[i]);
            assertEq(scaled.price(), 17_955_200_000 * scales[i]);
        }
        BandDouble finer = _band(address(new TokenDouble(29)));
        TokenDouble coarse = new TokenDouble(0);
        vm.expectRevert(stdError.arithmeticError);
        this.deploy(finer, address(coarse));
    }

    function deploy(BandDouble band_, address loanToken) external returns (MorphoBandOracle) {
        return new MorphoBandOracle(IBand(address(band_)), NVDA, IERC20Metadata(loanToken), owner);
    }

    function test_TheConstructorRefusesAScaleAPriceCouldOverflow() public {
        BandDouble coarse = _band(address(new TokenDouble(0)));
        coarse.setQuote(NVDA, BandDouble.Quote(2, 1, type(uint64).max, 0, type(uint64).max, type(uint64).max));
        MorphoBandOracle widest = this.deploy(coarse, address(new TokenDouble(29)));
        assertEq(widest.scaleFactor(), 1e57);
        assertEq(widest.price(), uint256(type(uint64).max) * 1e57);
        TokenDouble finest = new TokenDouble(30);
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.ScaleTooLarge.selector, 1e58));
        this.deploy(coarse, address(finest));
    }

    function testFuzz_APriceIsTheLowEdgeInMorphosScale(
        uint8 state,
        uint64 low,
        uint8 loanDecimals,
        uint8 collateralDecimals
    ) public {
        state = uint8(bound(state, 1, 3));
        low = uint64(bound(low, 1, type(uint64).max));
        loanDecimals = uint8(bound(loanDecimals, 0, 36));
        collateralDecimals =
            uint8(bound(collateralDecimals, loanDecimals > 29 ? loanDecimals - 29 : 0, uint256(loanDecimals) + 28));
        BandDouble other = _band(address(new TokenDouble(collateralDecimals)));
        other.setQuote(NVDA, BandDouble.Quote(state, 1, low, 0, low, low));
        MorphoBandOracle scaled = this.deploy(other, address(new TokenDouble(loanDecimals)));
        uint256 price = scaled.price();
        assertEq(price, uint256(low) * 10 ** (uint256(loanDecimals) + 28 - collateralDecimals));
        assertGt(price, 0);
    }

    function test_OnlyTheOwnerRepointsTheOracle() public {
        BandDouble other = _band(address(nvda));
        address stranger = makeAddr("stranger");
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, stranger));
        vm.prank(stranger);
        oracle.setBand(IBand(address(other)));
        assertEq(address(oracle.band()), address(band));
    }

    function test_ARepointNeedsABandThatNamesTheSameStockToken() public {
        TokenDouble token = new TokenDouble(18);
        BandDouble other = _band(address(token));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.AssetMismatch.selector, address(other), address(token)));
        vm.prank(owner);
        oracle.setBand(IBand(address(other)));
        BandDouble share = _band(address(0));
        vm.expectRevert(abi.encodeWithSelector(MorphoBandOracle.AssetMismatch.selector, address(share), address(0)));
        vm.prank(owner);
        oracle.setBand(IBand(address(share)));
        assertEq(address(oracle.band()), address(band));
    }

    function test_ARepointReadsTheNewBandAndItsSequencer() public {
        BandDouble other = _band(address(nvda));
        other.setSequencer(address(0x5E9), false);
        vm.expectEmit(address(oracle));
        emit MorphoBandOracle.BandSet(IBand(address(band)), IBand(address(other)));
        vm.prank(owner);
        oracle.setBand(IBand(address(other)));
        assertEq(address(oracle.band()), address(other));
        vm.expectRevert(MorphoBandOracle.SequencerNotSettled.selector);
        oracle.price();
        other.setSequencer(address(0x5E9), true);
        assertEq(oracle.price(), 17_955_200_000 * 1e16);
        band.setSequencer(address(0), false);
        vm.prank(owner);
        oracle.setBand(IBand(address(band)));
        assertEq(oracle.price(), LOW_IN_MORPHO);
    }

    function test_OwnershipMovesInTwoStepsAndCannotBeRenounced() public {
        address timelock = makeAddr("timelock");
        vm.prank(owner);
        oracle.transferOwnership(timelock);
        assertEq(oracle.owner(), owner);
        assertEq(oracle.pendingOwner(), timelock);
        vm.prank(timelock);
        oracle.acceptOwnership();
        assertEq(oracle.owner(), timelock);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, owner));
        vm.prank(owner);
        oracle.setBand(IBand(address(band)));
        vm.expectRevert(MorphoBandOracle.OwnershipCannotBeRenounced.selector);
        vm.prank(timelock);
        oracle.renounceOwnership();
    }

    function test_HaltIsTheBandsForTheAsset() public {
        band.setHalt(NVDA, BandDouble.Halt(true, 1_790_390_600, 1_790_389_990, true));
        band.setHalt(TSLA, BandDouble.Halt(false, 7, 8, false));
        (bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused) = oracle.halt();
        assertTrue(signedHalt);
        assertEq(until, 1_790_390_600);
        assertEq(issuedAt, 1_790_389_990);
        assertTrue(oraclePaused);
    }

    function test_GasOfEachCall() public {
        oracle.price();
        vm.snapshotGasLastCall("price");
        BandDouble other = _band(address(nvda));
        vm.prank(owner);
        oracle.setBand(IBand(address(other)));
        vm.snapshotGasLastCall("setBand");
    }
}
