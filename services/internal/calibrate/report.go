// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"fmt"

	"github.com/tapehouse/tapehouse/services/internal/backtest"
)

// divergence is how far the filter's volatility may stand from the method's before the report flags it.
const divergence = 1.5

func summary(label string, r backtest.Report) string {
	l, c := r.Leverage, r.Capacity
	return fmt.Sprintf("%s: %d long and %d mixed portfolios lose more than their margin at the deployed %d× cap; capacity %.1f%%, bad debt %.4f%%; %d of %d long portfolios covered at their worst closure",
		label, l.AtDeployedCap["long"].Portfolios, l.AtDeployedCap["mixed"].Portfolios, l.DeployedCap, 100*c.Model, 100*c.BadDebtModel,
		r.Coverage["long"].CoveredAtWorst, r.Coverage["long"].Portfolios)
}

// Lines describes the proposal: each value that changes with its target, the filter's volatility beside the method's,
// the backtest of both sets, the distance of the band from the reference, and the tail.
func (p Proposal) Lines() []string {
	var out []string
	change := func(what string, was, is, target uint64) {
		if was != is {
			out = append(out, fmt.Sprintf("%s %d → %d (target %d)", what, was, is, target))
		}
	}
	for i, a := range p.Assets {
		change(a+" volatility", uint64(p.Base.Volatilities[i]), uint64(p.Proposed.Volatilities[i]), uint64(p.Targets.Volatilities[i]))
		change(a+" gap", uint64(p.Base.Gaps[i]), uint64(p.Proposed.Gaps[i]), uint64(p.Targets.Gaps[i]))
		change(a+" selling depth", uint64(p.Base.Depths[2*i]), uint64(p.Proposed.Depths[2*i]), uint64(p.Proposed.Depths[2*i]))
		change(a+" buying depth", uint64(p.Base.Depths[2*i+1]), uint64(p.Proposed.Depths[2*i+1]), uint64(p.Proposed.Depths[2*i+1]))
	}
	for k, pair := range Pairs(p.Assets) {
		change(pair[0]+"/"+pair[1]+" correlation", uint64(p.Base.Correlations[k]), uint64(p.Proposed.Correlations[k]), uint64(p.Targets.Correlations[k]))
	}
	if len(out) == 0 {
		out = append(out, "no value changes")
	}
	if p.Proposed.Shrunk > 0 {
		out = append(out, fmt.Sprintf("the correlations were blended 2^-%d toward the current matrix to stay positive definite", p.Proposed.Shrunk))
	}
	for i, a := range p.Assets {
		if i >= len(p.ReferenceVolatilities) {
			break
		}
		line := fmt.Sprintf("reference volatility %s %d against the method's %d", a, p.ReferenceVolatilities[i], p.Targets.Volatilities[i])
		if ratio := float64(p.ReferenceVolatilities[i]) / float64(p.Targets.Volatilities[i]); ratio > divergence || ratio < 1/divergence {
			line += ": diverges beyond ×1.5"
		}
		out = append(out, line)
	}
	out = append(out, summary("backtest of the current parameters", p.Current), summary("backtest of the proposed parameters", p.Next))
	for _, d := range p.Reference {
		line := fmt.Sprintf("band %s over %d values: centre %.2f bps median, %.2f bps at the 95th percentile from the reference; half-width %.2f× its; holds the next value %.2f%% against the reference's %.2f%%",
			d.Asset, d.Points, d.MedianMidBps, d.P95MidBps, d.MedianHalfRatio, 100*d.BandCoverage, 100*d.ReferenceCoverage)
		if d.Redeploy {
			line += ": redeploy the band"
		}
		out = append(out, line)
	}
	if t := p.Tail; t != nil {
		line := fmt.Sprintf("tail: threshold %.6f gaps, scale %.6f, probability %.6f over %d excesses; reopen ratio %.2f to %.2f; constants %v against the deployed %v",
			t.Threshold, t.Scale, t.Probability, t.Excesses, t.ReopenLow, t.ReopenHigh, t.Scaled, t.Deployed)
		if t.Redeploy {
			line += ": deploy a new gap cover"
		}
		out = append(out, line)
	}
	out = append(out, fmt.Sprintf("weekend leverage cap: the rule gives %d× against the deployed %d×; moving it is a redeployment", p.Current.Leverage.Cap, p.Current.Leverage.DeployedCap))
	out = append(out, "simulation: "+p.Simulation)
	return out
}
