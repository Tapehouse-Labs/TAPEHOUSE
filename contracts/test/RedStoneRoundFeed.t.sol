// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IMultiplier, RedStoneRoundFeed} from "../src/RedStoneRoundFeed.sol";

/// @notice The feed against a real payload of RedStone's `redstone-primary-prod` data service: ten signed packages,
/// five signers' of TSLA and of AMZN, captured from RedStone's gateway.
contract MultiplierDouble is IMultiplier {
    uint256 public override uiMultiplier = 1e18;

    function set(uint256 value) external {
        uiMultiplier = value;
    }
}

contract RedStoneRoundFeedTest is Test {
    uint256 internal constant PACKAGE_SIZE = 142;
    bytes9 internal constant MARKER = 0x000002ed57011e0000;

    bytes internal payload;
    uint256 internal timestamp;
    RedStoneRoundFeed internal tsla;
    RedStoneRoundFeed internal amzn;

    function setUp() public {
        payload = vm.parseBytes(vm.trim(vm.readFile("test/data/redstone-prices.hex")));
        uint256 sigStart = payload.length - 9 - 3 - 2 - 65;
        for (uint256 i = sigStart - 13; i < sigStart - 7; ++i) {
            timestamp = (timestamp << 8) | uint8(payload[i]);
        }
        timestamp /= 1000;
        tsla = new RedStoneRoundFeed(bytes32("TSLA"), "TSLA / USD, RedStone", IMultiplier(address(0)));
        amzn = new RedStoneRoundFeed(bytes32("AMZN"), "AMZN / USD, RedStone", IMultiplier(address(0)));
        vm.warp(timestamp + 30);
    }

    function _update(RedStoneRoundFeed feed, bytes memory data) internal returns (bool ok, bytes memory result) {
        return address(feed).call(bytes.concat(feed.update.selector, data));
    }

    function _repack(uint256[] memory indices) internal view returns (bytes memory out) {
        for (uint256 i; i < indices.length; ++i) {
            out = bytes.concat(out, _slice(payload, indices[i] * PACKAGE_SIZE, PACKAGE_SIZE));
        }
        out = bytes.concat(out, bytes2(uint16(indices.length)), bytes3(0), MARKER);
    }

    function _slice(bytes memory data, uint256 start, uint256 size) internal pure returns (bytes memory out) {
        out = new bytes(size);
        for (uint256 i; i < size; ++i) {
            out[i] = data[start + i];
        }
    }

    function test_RecordsTheMedianOfTheSignersAsTheFirstRound() public {
        (bool ok,) = _update(tsla, payload);
        assertTrue(ok);
        (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredIn) =
            tsla.latestRoundData();
        assertEq(roundId, (uint80(1) << 64) | 1);
        assertEq(answeredIn, roundId);
        assertEq(startedAt, timestamp);
        assertEq(updatedAt, timestamp);
        assertGt(answer, 10_000_000_000, "above $100");
        assertLt(answer, 200_000_000_000, "below $2,000");
        assertEq(tsla.decimals(), 8);
        assertEq(tsla.description(), "TSLA / USD, RedStone");
        (uint80 id, int256 again,,,) = tsla.getRoundData(roundId);
        assertEq(id, roundId);
        assertEq(again, answer);
    }

    function test_EachFeedTakesItsOwnValue() public {
        (bool a,) = _update(tsla, payload);
        (bool b,) = _update(amzn, payload);
        assertTrue(a && b);
        (, int256 t,,,) = tsla.latestRoundData();
        (, int256 z,,,) = amzn.latestRoundData();
        assertTrue(t != z);
    }

    function test_AnEmptyFeedAndAnUnknownRoundHaveNoData() public {
        vm.expectRevert("No data present");
        tsla.latestRoundData();
        _update(tsla, payload);
        vm.expectRevert("No data present");
        tsla.getRoundData(2);
        vm.expectRevert("No data present");
        tsla.getRoundData((uint80(2) << 64) | 1);
        vm.expectRevert("No data present");
        tsla.getRoundData(1);
    }

    function test_APayloadIsRecordedOnceAndOnlyNewerOnesFollow() public {
        _update(tsla, payload);
        (bool ok, bytes memory reason) = _update(tsla, payload);
        assertFalse(ok);
        assertEq(bytes4(reason), RedStoneRoundFeed.NotNewer.selector);
    }

    function test_AStaleOrFuturePayloadIsRefused() public {
        vm.warp(timestamp + 181);
        (bool ok, bytes memory reason) = _update(tsla, payload);
        assertFalse(ok);
        assertEq(bytes4(reason), RedStoneRoundFeed.TimestampIsTooOld.selector);
        vm.warp(timestamp - 61);
        (ok, reason) = _update(tsla, payload);
        assertFalse(ok);
        assertEq(bytes4(reason), RedStoneRoundFeed.TimestampIsFromTheFuture.selector);
        vm.warp(timestamp + 180);
        (ok,) = _update(tsla, payload);
        assertTrue(ok);
    }

    function test_ATamperedValueIsRefused() public {
        bytes memory forged = payload;
        forged[40] = bytes1(uint8(forged[40]) ^ 0x01);
        (bool ok, bytes memory reason) = _update(tsla, forged);
        assertFalse(ok);
        bytes4 selector = bytes4(reason);
        assertTrue(
            selector == RedStoneRoundFeed.SignerNotAuthorised.selector
                || selector == RedStoneRoundFeed.InvalidSignature.selector
                || selector == RedStoneRoundFeed.InsufficientNumberOfUniqueSigners.selector,
            "a forged package is refused"
        );
    }

    function test_FewerThanThreeUniqueSignersAreRefused() public {
        uint256[] memory two = new uint256[](2);
        two[0] = 0;
        two[1] = 1;
        (bool ok, bytes memory reason) = _update(tsla, _repack(two));
        assertFalse(ok);
        assertEq(bytes4(reason), RedStoneRoundFeed.InsufficientNumberOfUniqueSigners.selector);

        uint256[] memory repeated = new uint256[](3);
        (ok, reason) = _update(tsla, _repack(repeated));
        assertFalse(ok);
        assertEq(bytes4(reason), RedStoneRoundFeed.InsufficientNumberOfUniqueSigners.selector);
    }

    function test_AFeedTheDataServiceDoesNotCarryIsRefused() public {
        RedStoneRoundFeed other = new RedStoneRoundFeed(bytes32("NOPE"), "none", IMultiplier(address(0)));
        (bool ok, bytes memory reason) = _update(other, payload);
        assertFalse(ok);
        assertEq(bytes4(reason), RedStoneRoundFeed.InsufficientNumberOfUniqueSigners.selector);
    }

    function test_APayloadWithoutTheMarkerIsRefused() public {
        (bool ok, bytes memory reason) = _update(tsla, hex"1234");
        assertFalse(ok);
        assertEq(bytes4(reason), RedStoneRoundFeed.InvalidPayload.selector);
    }

    function test_ATokensMultiplierScalesTheAnswerToTheTokensPrice() public {
        MultiplierDouble multiplier = new MultiplierDouble();
        RedStoneRoundFeed plain = new RedStoneRoundFeed(bytes32("TSLA"), "share", IMultiplier(address(0)));
        RedStoneRoundFeed token = new RedStoneRoundFeed(bytes32("TSLA"), "token", multiplier);
        multiplier.set(2.5e18);
        _update(plain, payload);
        _update(token, payload);
        (, int256 share,,,) = plain.latestRoundData();
        (, int256 price,,,) = token.latestRoundData();
        assertEq(price, share * 25 / 10, "the token's price is the share's times its multiplier");
        assertEq(address(token.stockToken()), address(multiplier));
    }
}
