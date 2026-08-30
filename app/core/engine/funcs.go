package engine

import (
	"fmt"
	"hash/fnv"
)

// wraptmpl wraps a string in the template delimiters.
func wraptmpl(s string) string {
	return "{{ " + s + " }}"
}

// Defaults for portBlock. The window stops below 32768 because that is where
// Linux starts handing out ephemeral source ports (net.ipv4.ip_local_port_range);
// macOS starts at 49152. Staying under the lower of the two keeps a reserved port
// from being claimed by an unrelated outbound connection, which shows up as an
// intermittent bind failure that is painful to trace.
const (
	defaultPortRangeStart = 20000
	defaultPortRangeEnd   = 32767
	defaultPortBlockSize  = 16

	minUnprivilegedPort = 1024
	maxPort             = 65535
)

// portBlock maps a seed string to a deterministic, block-aligned base port.
//
// It takes the seed and optionally start, end, and size:
//
//	{{ portBlock .Ctx.OutputDirBase }}
//	{{ portBlock .Ctx.OutputDirBase 20000 32767 16 }}
//
// The result is the first port of a block of `size` consecutive ports, so a
// caller adds an offset per service: base+1, base+2, and so on up to size-1.
//
// Aligning to block boundaries is the point. Hashing straight to a port and then
// incrementing lets two seeds land a few ports apart and partially overlap, so
// some services bind and others do not. With alignment, two seeds either share
// the whole block or share nothing.
//
// This is stable, not unique. The same seed always yields the same block, with no
// state to store or clean up, but distinct seeds can collide. Collision odds
// follow the birthday bound: the defaults give 798 blocks, so roughly 5.5% at 10
// concurrent seeds and 21% at 20. Treat the result as a good default that a human
// can override, not as a reservation.
func portBlock(seed string, opts ...int) (int, error) {
	start, end, size := defaultPortRangeStart, defaultPortRangeEnd, defaultPortBlockSize

	switch len(opts) {
	case 0:
	case 3:
		start, end, size = opts[0], opts[1], opts[2]
	default:
		return 0, fmt.Errorf("portBlock: expected 0 or 3 options (start, end, size), got %d", len(opts))
	}

	if seed == "" {
		return 0, fmt.Errorf("portBlock: seed must not be empty")
	}

	if size <= 0 {
		return 0, fmt.Errorf("portBlock: size must be positive, got %d", size)
	}

	if start < minUnprivilegedPort {
		return 0, fmt.Errorf("portBlock: start must be >= %d, got %d (ports below that need root)", minUnprivilegedPort, start)
	}

	if end > maxPort {
		return 0, fmt.Errorf("portBlock: end must be <= %d, got %d", maxPort, end)
	}

	if end <= start {
		return 0, fmt.Errorf("portBlock: end (%d) must be greater than start (%d)", end, start)
	}

	blocks := (end - start + 1) / size
	if blocks < 1 {
		return 0, fmt.Errorf("portBlock: range %d-%d is smaller than one block of %d", start, end, size)
	}

	// fnv32a rather than one of sprig's checksum functions: those return hex for
	// everything except adler32, and feeding hex to atoi silently yields 0, which
	// hands every seed the same block.
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))

	return start + int(h.Sum32()%uint32(blocks))*size, nil
}
