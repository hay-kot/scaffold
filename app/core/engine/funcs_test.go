package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_portBlock_Deterministic(t *testing.T) {
	first, err := portBlock("hive-desktop-0lyqf7")
	require.NoError(t, err)

	for range 10 {
		again, err := portBlock("hive-desktop-0lyqf7")
		require.NoError(t, err)
		assert.Equal(t, first, again, "same seed must always yield the same block")
	}
}

func Test_portBlock_BlockAligned(t *testing.T) {
	// Alignment is the property that stops two seeds from partially overlapping.
	seeds := []string{"recipinned", "hive-desktop", "mealie-nla6iz", "infra-7dn6nd", "a", "zzzzzzzzzz"}

	for _, seed := range seeds {
		t.Run(seed, func(t *testing.T) {
			got, err := portBlock(seed)
			require.NoError(t, err)

			offset := got - defaultPortRangeStart
			assert.Zero(t, offset%defaultPortBlockSize, "block base must sit on a block boundary")
			assert.GreaterOrEqual(t, got, defaultPortRangeStart)
			assert.LessOrEqual(t, got+defaultPortBlockSize-1, defaultPortRangeEnd, "whole block must fit inside the range")
		})
	}
}

func Test_portBlock_StaysBelowEphemeralRange(t *testing.T) {
	// Linux hands out ephemeral ports from 32768 up. A default block that reached
	// into that range would bind fine and then fail intermittently later.
	for i := range 500 {
		got, err := portBlock(fmt.Sprintf("seed-%d", i))
		require.NoError(t, err)
		assert.Less(t, got+defaultPortBlockSize-1, 32768)
	}
}

func Test_portBlock_CustomRange(t *testing.T) {
	got, err := portBlock("recipinned", 40000, 49000, 10)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, got, 40000)
	assert.LessOrEqual(t, got+9, 49000)
	assert.Zero(t, (got-40000)%10)
}

func Test_portBlock_DistinctSeedsSpread(t *testing.T) {
	// Not a uniqueness guarantee, just a check that the hash is not degenerate.
	seen := map[int]int{}
	const n = 200

	for i := range n {
		got, err := portBlock(fmt.Sprintf("project-%d", i))
		require.NoError(t, err)
		seen[got]++
	}

	assert.Greater(t, len(seen), n*3/4, "fnv32a should spread %d seeds across most blocks", n)
}

func Test_portBlock_Errors(t *testing.T) {
	tests := []struct {
		name string
		seed string
		opts []int
		msg  string
	}{
		{name: "empty seed", seed: "", msg: "seed must not be empty"},
		{name: "wrong option count", seed: "x", opts: []int{1000, 2000}, msg: "expected 0 or 3 options"},
		{name: "zero size", seed: "x", opts: []int{20000, 30000, 0}, msg: "size must be positive"},
		{name: "negative size", seed: "x", opts: []int{20000, 30000, -4}, msg: "size must be positive"},
		{name: "privileged start", seed: "x", opts: []int{80, 30000, 16}, msg: "start must be >="},
		{name: "end above max", seed: "x", opts: []int{20000, 70000, 16}, msg: "end must be <="},
		{name: "inverted range", seed: "x", opts: []int{30000, 20000, 16}, msg: "must be greater than start"},
		{name: "range smaller than block", seed: "x", opts: []int{20000, 20004, 16}, msg: "smaller than one block"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := portBlock(tt.seed, tt.opts...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.msg)
		})
	}
}

func Test_portBlock_InTemplate(t *testing.T) {
	e := New()

	out, err := e.TmplString(`{{ portBlock .Seed }}`, Vars{"Seed": "recipinned"})
	require.NoError(t, err)

	direct, err := portBlock("recipinned")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprint(direct), out)

	// The offset pattern a port scaffold actually uses.
	out, err = e.TmplString(`{{ $b := portBlock .Seed }}{{ add $b 1 }},{{ add $b 2 }}`, Vars{"Seed": "recipinned"})
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d,%d", direct+1, direct+2), out)
}

func Test_portBlock_TemplateErrorSurfaces(t *testing.T) {
	e := New()

	_, err := e.TmplString(`{{ portBlock "" }}`, Vars{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "seed must not be empty")
}
