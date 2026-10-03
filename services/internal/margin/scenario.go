// SPDX-License-Identifier: MIT OR Apache-2.0

// Package margin evaluates the margin program's arithmetic bit for bit: the validity of a correlation matrix, the
// scenario set, the portfolio requirement with its liquidity add-on and pool reads, and the regime the current
// requirement follows through the band's session. Amounts are USD with 18 decimals and returns are in millionths.
package margin

import (
	"encoding/binary"
	"math/big"
)

// Constants of the margin program.
const (
	ppm         = 1_000_000
	wad         = 1_000_000_000_000_000_000
	factorScale = 1_000_000_000_000
	day         = 86_400
	pivotFloor  = factorScale / ppm
	// ZMarket is the standard normal's expected shortfall at 99%, in millionths.
	ZMarket = 2_665_214
	// One is a correlation of one, in basis points.
	One = 10_000
)

// Sizes are the lattice sizes the engine draws.
var Sizes = []int{32, 64, 128, 256}

var means = map[int][]int32{
	32: {-2252211, -1683275, -1420533, -1231294, -1078398, -947376, -830935, -724828, -626336, -533591, -445235,
		-360235, -277767, -197152, -117800, -39186},
	64: {-2510185, -1994237, -1764207, -1602344, -1474296, -1366770, -1273110, -1189477, -1113440, -1043356, -978060,
		-916692, -858599, -803270, -750300, -699356, -650167, -602505, -556174, -511008, -466862, -423608, -381134,
		-339336, -298123, -257411, -217121, -177182, -137522, -98078, -58787, -19586},
	128: {-2747787, -2272583, -2065885, -1922589, -1810660, -1717753, -1637708, -1566980, -1503333, -1445259,
		-1391692, -1341848, -1295133, -1251087, -1209345, -1169610, -1131642, -1095239, -1060233, -1026480, -993859,
		-962262, -931598, -901787, -872756, -844443, -816791, -789750, -763275, -737324, -711861, -686851, -662264,
		-638071, -614246, -590764, -567604, -544744, -522166, -499851, -477782, -455943, -434320, -412897, -391663,
		-370604, -349708, -328964, -308360, -287886, -267533, -247289, -227147, -207096, -187128, -167235, -147407,
		-127637, -107918, -88239, -68595, -48978, -29379, -9792},
	256: {-2969102, -2526473, -2337436, -2207729, -2107268, -2024501, -1953679, -1891500, -1835885, -1785435,
		-1739163, -1696343, -1656427, -1618989, -1583693, -1550268, -1518490, -1488176, -1459172, -1431346, -1404587,
		-1378797, -1353893, -1329802, -1306459, -1283807, -1261796, -1240379, -1219517, -1199172, -1179313, -1159907,
		-1140929, -1122354, -1104158, -1086320, -1068821, -1051644, -1034772, -1018189, -1001881, -985836, -970041,
		-954484, -939154, -924042, -909139, -894434, -879921, -865591, -851436, -837449, -823625, -809957, -796438,
		-783063, -769826, -756723, -743749, -730899, -718169, -705553, -693050, -680653, -668360, -656168, -644072,
		-632070, -620158, -608333, -596593, -584935, -573355, -561852, -550423, -539065, -527777, -516555, -505398,
		-494303, -483269, -472294, -461375, -450511, -439700, -428940, -418229, -407566, -396949, -386377, -375848,
		-365360, -354913, -344504, -334132, -323796, -313494, -303226, -292989, -282783, -272607, -262459, -252337,
		-242242, -232171, -222123, -212098, -202094, -192111, -182146, -172199, -162270, -152357, -142458, -132573,
		-122702, -112842, -102993, -93154, -83325, -73503, -63688, -53880, -44076, -34277, -24481, -14688, -4896},
}

// Parameters are what scenarios are drawn from, in the stored order of the assets. Market is the market asset's
// position, or -1 for the equal-weighted portfolio.
type Parameters struct {
	Symbols      [][32]byte
	Volatilities []uint32
	Correlations []uint16
	Gaps         []uint32
	Market       int
}

// Set is a scenario set: every row, as returns in millionths in the stored order of the assets.
type Set struct {
	Parameters Parameters
	Size       int
	order      []int
	rows       [][]int32
}

// Count is the number of scenarios in a set of lattice size size over n assets.
func Count(size, n int) int {
	return 3*size + 4 + 2*n
}

// Pair is the index of the pair (i, j), i < j, in the upper triangle of an n × n matrix read row by row.
func Pair(n, i, j int) int {
	return i*n - i*(i+1)/2 + (j - i - 1)
}

func generator(size int) int {
	switch size {
	case 32:
		return 5
	case 64:
		return 11
	}
	return 13
}

func mean(size, j int) int64 {
	if j < size/2 {
		return int64(means[size][j])
	}
	return -int64(means[size][size-1-j])
}

func shift(size, d int) int {
	frac := new(big.Int).Mul(big.NewInt(int64(d)), big.NewInt(618_033_988_749_894_848))
	frac.Mod(frac, big.NewInt(wad))
	return int(frac.Mul(frac, big.NewInt(int64(size))).Quo(frac, big.NewInt(wad)).Int64())
}

func draw(size, j, d int) int64 {
	g := 1
	for range d {
		g = g * generator(size) % size
	}
	return mean(size, (j*g+shift(size, d))%size)
}

// NewSet draws the set of lattice size size (one of Sizes) over a horizon of horizon seconds.
func NewSet(p Parameters, size int, horizon uint64) *Set {
	n := len(p.Symbols)
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	for i := 1; i < n; i++ {
		for k := i; k > 0 && string(p.Symbols[order[k-1]][:]) > string(p.Symbols[order[k]][:]); k-- {
			order[k-1], order[k] = order[k], order[k-1]
		}
	}
	rank := make([]int, n)
	for r, i := range order {
		rank[i] = r
	}
	rho := func(a, b int) int64 {
		if a == b {
			return One
		}
		return int64(p.Correlations[Pair(n, min(a, b), max(a, b))])
	}
	s := big.NewInt(factorScale)
	factor := make([]*big.Int, n*n)
	for i := range factor {
		factor[i] = new(big.Int)
	}
	dot := func(a, b, upto int) *big.Int {
		sum := new(big.Int)
		for k := range upto {
			sum.Add(sum, new(big.Int).Mul(factor[a*n+k], factor[b*n+k]))
		}
		return sum
	}
	for j := range n {
		acc := new(big.Int).Mul(s, s)
		acc.Sub(acc, dot(j, j, j))
		if acc.Sign() < 0 {
			acc.SetInt64(0)
		}
		if pivot := acc.Sqrt(acc); pivot.Cmp(big.NewInt(pivotFloor)) >= 0 {
			factor[j*n+j] = pivot
		}
		for i := j + 1; i < n; i++ {
			acc := big.NewInt(rho(order[i], order[j]) * (factorScale / One))
			acc.Mul(acc, s).Sub(acc, dot(i, j, j))
			if factor[j*n+j].Sign() > 0 {
				acc.Quo(acc, factor[j*n+j])
				if acc.Cmp(s) > 0 {
					acc.Set(s)
				} else if acc.Cmp(new(big.Int).Neg(s)) < 0 {
					acc.Neg(s)
				}
				factor[i*n+j] = acc
			}
		}
	}
	root := new(big.Int).SetUint64(horizon)
	root.Mul(root, s).Quo(root, big.NewInt(day)).Sqrt(root)
	weight := func(k int) int64 {
		if p.Market < 0 || p.Market == k {
			return 1
		}
		return 0
	}
	vol := func(k int) *big.Int { return big.NewInt(int64(p.Volatilities[k])) }
	variance := new(big.Int)
	for a := range n {
		for b := range n {
			term := new(big.Int).Mul(vol(a), vol(b))
			variance.Add(variance, term.Mul(term, big.NewInt(weight(a)*weight(b)*rho(a, b))))
		}
	}
	sigma := variance.Mul(variance, big.NewInt(One)).Sqrt(variance)
	moves := make([]*big.Int, n)
	for i := range n {
		covariance := new(big.Int)
		for k := range n {
			term := new(big.Int).Mul(vol(i), vol(k))
			covariance.Add(covariance, term.Mul(term, big.NewInt(weight(k)*rho(i, k))))
		}
		m := new(big.Int).Mul(big.NewInt(ZMarket), root)
		moves[i] = m.Mul(m, covariance).Quo(m, new(big.Int).Mul(s, sigma))
	}
	scaled := func(x *big.Int, i int) int32 {
		v := new(big.Int).Mul(x, vol(i))
		return clamp(v.Mul(v, root).Quo(v, big.NewInt(ppm*ppm)))
	}
	rows := make([][]int32, Count(size, n))
	for index := range rows {
		row := make([]int32, n)
		switch {
		case index < size:
			draws := make([]*big.Int, n)
			for d := range n {
				draws[d] = big.NewInt(draw(size, index, d))
			}
			for c, i := range order {
				x := new(big.Int)
				for k := 0; k <= c; k++ {
					x.Add(x, new(big.Int).Mul(factor[c*n+k], draws[k]))
				}
				row[i] = scaled(x.Quo(x, s), i)
			}
		case index < 2*size:
			x := big.NewInt(draw(size, index-size, 0))
			for i := range row {
				row[i] = scaled(x, i)
			}
		case index < 3*size:
			for i := range row {
				row[i] = scaled(big.NewInt(draw(size, index-2*size, rank[i])), i)
			}
		case index < 3*size+4:
			shock := index - 3*size
			for i := range row {
				magnitude := big.NewInt(int64(p.Gaps[i]))
				if shock < 2 {
					magnitude = new(big.Int).Set(moves[i])
				}
				if shock%2 == 0 {
					magnitude.Neg(magnitude)
				}
				row[i] = clamp(magnitude)
			}
		default:
			alone := index - 3*size - 4
			i := order[alone/2]
			gap := big.NewInt(int64(p.Gaps[i]))
			if alone%2 == 0 {
				gap.Neg(gap)
			}
			row[i] = clamp(gap)
		}
		rows[index] = row
	}
	return &Set{Parameters: p, Size: size, order: order, rows: rows}
}

// Row is scenario index of Count, as returns in millionths in the stored order of the assets.
func (s *Set) Row(index int) []int32 {
	return s.rows[index]
}

// Encoded is every scenario in the canonical order of the assets, ascending by symbol, as big-endian int32.
func (s *Set) Encoded() []byte {
	out := make([]byte, 0, len(s.rows)*len(s.order)*4)
	for _, row := range s.rows {
		for _, i := range s.order {
			out = binary.BigEndian.AppendUint32(out, uint32(row[i]))
		}
	}
	return out
}

func clamp(r *big.Int) int32 {
	if r.Cmp(big.NewInt(-ppm)) < 0 {
		return -ppm
	}
	if !r.IsInt64() || r.Int64() > 1<<31-1 {
		return 1<<31 - 1
	}
	return int32(r.Int64())
}
