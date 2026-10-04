// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

// PositionState is all an alert remembers of one position: whether its threshold alert is spent until the ratio
// recovers, the closure its weekend alert was sent for, whether it is short and was told so, the auction it was told
// of and whether it was told it cannot be judged.
type PositionState struct {
	Alerted  bool   `json:"alerted,omitempty"`
	Closure  uint64 `json:"closure,omitempty"`
	Short    bool   `json:"short,omitempty"`
	Auction  uint64 `json:"auction,omitempty"`
	Unjudged bool   `json:"unjudged,omitempty"`
}
