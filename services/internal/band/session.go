// SPDX-License-Identifier: MIT OR Apache-2.0

package band

import "math/big"

// Feed IDs of the signed New York market status, written together or not at all.
const (
	CurrentStatus  = "NY_MARKET_CURRENT_STATUS"
	NextStatus     = "NY_MARKET_NEXT_STATUS"
	NextChangeTime = "NY_MARKET_NEXT_CHANGE_TIME"
)

// Timing of Chainlink's 24/5 session, in milliseconds.
const (
	statusScale            = 100_000_000
	PostMarketMs           = 14_400_000
	ReopenBeforeOpenMs     = 48_600_000
	ReopenBeforeMidnightMs = 14_400_000
	StatusMaxAgeMs         = 3_600_000
)

// Nyse is NYSE's state, as RedStone signs it.
type Nyse uint8

// NYSE's states; zero is none.
const (
	Regular     Nyse = 1
	ClosedShort Nyse = 2
	ClosedLong  Nyse = 3
)

// Status is the signed market status with the package timestamp of all three values.
type Status struct {
	Current   Nyse
	Next      Nyse
	ChangeMs  uint64
	PackageMs uint64
}

func decodeNyse(value *big.Int) Nyse {
	if !value.IsUint64() {
		return 0
	}
	switch value.Uint64() {
	case 100_000_000:
		return Regular
	case 101_000_000:
		return ClosedShort
	case 102_000_000:
		return ClosedLong
	}
	return 0
}

// DecodeStatus decodes the stored current status, next status and change time, and reports whether all three come
// from one package and hold known values.
func DecodeStatus(values [3]*big.Int, packageMs [3]uint64) (Status, bool) {
	if packageMs[1] != packageMs[0] || packageMs[2] != packageMs[0] {
		return Status{}, false
	}
	change := new(big.Int).Quo(values[2], big.NewInt(statusScale))
	s := Status{Current: decodeNyse(values[0]), Next: decodeNyse(values[1]), PackageMs: packageMs[0]}
	if !change.IsUint64() || change.Sign() == 0 || s.Current == 0 || s.Next == 0 {
		return Status{}, false
	}
	s.ChangeMs = change.Uint64()
	return s, true
}

// Session is Chainlink's 24/5 session at one instant. NyseNext is zero once NYSE's change has passed, and
// BoundaryMs zero while the next boundary is not known.
type Session struct {
	Open       bool
	Nyse       Nyse
	NyseNext   Nyse
	ChangeMs   uint64
	BoundaryMs uint64
}

// Derive is the session at nowMs from the stored status and closeMs, the last regular close the band saw, and
// whether it is known: a status at most StatusMaxAgeMs old.
func Derive(status Status, known bool, closeMs, nowMs uint64) (Session, bool) {
	if !known || (nowMs > status.PackageMs && nowMs-status.PackageMs > StatusMaxAgeMs) {
		return Session{}, false
	}
	var s Session
	var reopensMs uint64
	if nowMs < status.ChangeMs {
		s.Nyse, s.NyseNext, s.ChangeMs = status.Current, status.Next, status.ChangeMs
		switch status.Next {
		case Regular:
			reopensMs = saturatingSub(status.ChangeMs, ReopenBeforeOpenMs)
		case ClosedShort:
			reopensMs = saturatingSub(status.ChangeMs, ReopenBeforeMidnightMs)
		}
	} else {
		if status.Current == Regular {
			closeMs = status.ChangeMs
		}
		s.Nyse = status.Next
	}
	postMarketEnds := saturatingAdd(closeMs, PostMarketMs)
	inPostMarket := nowMs < postMarketEnds
	switch {
	case s.Nyse == Regular:
		s.Open = true
	case s.Nyse == ClosedShort && s.NyseNext == ClosedLong:
		s.Open = inPostMarket
	case s.Nyse == ClosedShort:
		s.Open = true
	default:
		s.Open = inPostMarket || (reopensMs != 0 && nowMs >= reopensMs)
	}
	switch {
	case !s.Open && s.Nyse == ClosedLong:
		s.BoundaryMs = reopensMs
	case !s.Open:
	case inPostMarket && (s.Nyse == ClosedLong || (s.Nyse == ClosedShort && s.NyseNext == ClosedLong)):
		s.BoundaryMs = postMarketEnds
	case s.Nyse == Regular && s.NyseNext == ClosedLong:
		s.BoundaryMs = saturatingAdd(s.ChangeMs, PostMarketMs)
	}
	return s, true
}

// NextClose is the last regular close once status is written: its change time when signed in regular hours.
func NextClose(status Status, known bool, closeMs uint64) uint64 {
	if known && status.Current == Regular {
		return status.ChangeMs
	}
	return closeMs
}

// InRegularHours reports whether the signed status says NYSE is in regular hours in a known session.
func InRegularHours(s Session, known bool) bool {
	return known && s.Nyse == Regular
}

func saturatingSub(a, b uint64) uint64 {
	if b > a {
		return 0
	}
	return a - b
}
