// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"math/big"
	"slices"
)

// Kind is a kind of alert.
type Kind uint8

// The kinds of alert.
const (
	Threshold Kind = iota + 1
	Weekend
	Liquidation
	Unjudged
)

// Observation is what the rules know of one position at one tick, on the chain's clock.
type Observation struct {
	// Judged is whether the liquidator can judge the position; Held, whether it holds a Stock Token.
	Judged, Held bool
	// Equity and Requirement are the liquidator's, in USD with 18 decimals.
	Equity, Requirement *big.Int
	// Short is whether the liquidator would start the position's auction now, and Auction when the auction running now
	// started, in seconds, zero for none.
	Short   bool
	Auction uint64
	// Closing is whether the requirement is rising to a weekend or holiday close, and Closure that close's boundary.
	Closing bool
	Closure uint64
}

// rearm is how far above the threshold, in percent of it, the ratio must rise before its alert is armed again.
const rearm = 105

// Decide returns the alerts due for a position with the subscriber's threshold, in thousandths, and the state it is
// in once they are sent. A threshold alert is sent once and armed again only when the ratio of equity to requirement
// has risen 5% above the threshold; a weekend alert once for each closure; a liquidation alert once for each time the
// position falls short and for each auction that starts anew while it still is; an unjudged notice once for each
// time the position cannot be judged. A threshold alert due with a liquidation alert is spent without being sent:
// the liquidation alert says more. While it cannot be judged nothing else is said.
func Decide(threshold int, state PositionState, o Observation) ([]Kind, PositionState) {
	var due []Kind
	if !o.Judged {
		if !state.Unjudged {
			due = append(due, Unjudged)
		}
		state.Unjudged = true
		return due, state
	}
	state.Unjudged = false
	if o.Requirement.Sign() > 0 {
		scaled := new(big.Int).Mul(o.Equity, big.NewInt(1000))
		limit := new(big.Int).Mul(o.Requirement, big.NewInt(int64(threshold)))
		if scaled.Cmp(limit) < 0 && !state.Alerted {
			due, state.Alerted = append(due, Threshold), true
		}
		if state.Alerted && scaled.Mul(scaled, big.NewInt(100)).Cmp(limit.Mul(limit, big.NewInt(rearm))) >= 0 {
			state.Alerted = false
		}
	} else {
		state.Alerted = false
	}
	if o.Closing && o.Held && state.Closure != o.Closure {
		due, state.Closure = append(due, Weekend), o.Closure
	}
	switch {
	case !o.Short:
		state.Short, state.Auction = false, 0
	case !state.Short || (o.Auction != 0 && state.Auction != 0 && o.Auction != state.Auction):
		due, state.Short, state.Auction = append(due, Liquidation), true, o.Auction
	case state.Auction == 0:
		state.Auction = o.Auction
	}
	if slices.Contains(due, Liquidation) {
		due = slices.DeleteFunc(due, func(k Kind) bool { return k == Threshold })
	}
	return due, state
}

// revert puts back what old held of the state an alert of kind sets, for an alert that was not sent.
func revert(next *PositionState, old PositionState, kind Kind) {
	switch kind {
	case Threshold:
		next.Alerted = old.Alerted
	case Weekend:
		next.Closure = old.Closure
	case Liquidation:
		next.Short, next.Auction = old.Short, old.Auction
	case Unjudged:
		next.Unjudged = old.Unjudged
	}
}
