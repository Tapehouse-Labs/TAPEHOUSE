// SPDX-License-Identifier: MIT OR Apache-2.0

package margin

import (
	"encoding/json"
	"fmt"
)

// Launch reads the engine's parameters file, stylus/contracts/margin/parameters.json: each asset's first volatility
// and weekend gap, each pair's first correlation and the market asset, in the order of names, and each asset's
// selling then buying depth in whole USD.
func Launch(data []byte, names []string) (Parameters, []uint32, error) {
	var file struct {
		Market      string
		Volatility  map[string]struct{ Initial uint32 }
		Correlation map[string]struct{ Initial uint16 }
		WeekendGap  map[string]struct{ Initial uint32 }
		Depth       map[string]struct{ Initial []uint32 }
	}
	if err := json.Unmarshal(data, &file); err != nil {
		return Parameters{}, nil, err
	}
	p := Parameters{Market: -1}
	var depths []uint32
	for i, n := range names {
		v, okV := file.Volatility[n]
		g, okG := file.WeekendGap[n]
		d, okD := file.Depth[n]
		if !okV || !okG || !okD || len(d.Initial) != 2 {
			return Parameters{}, nil, fmt.Errorf("parameters: no %s", n)
		}
		var symbol [32]byte
		copy(symbol[:], n)
		p.Symbols = append(p.Symbols, symbol)
		p.Volatilities = append(p.Volatilities, v.Initial)
		p.Gaps = append(p.Gaps, g.Initial)
		depths = append(depths, d.Initial...)
		for _, m := range names[i+1:] {
			c, ok := file.Correlation[n+"/"+m]
			if !ok {
				return Parameters{}, nil, fmt.Errorf("parameters: no %s/%s", n, m)
			}
			p.Correlations = append(p.Correlations, c.Initial)
		}
		if n == file.Market {
			p.Market = i
		}
	}
	return p, depths, nil
}

// Pools reads each asset's pool from the engine's parameters file: the name of its entry in the registry's
// .uniswapV3 group, `<ASSET>_<QUOTE>_<fee>`.
func Pools(data []byte) (map[string]string, error) {
	var file struct{ Pool map[string]string }
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	return file.Pool, nil
}
