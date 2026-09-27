package main

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
)

var (
	errEmptyPool      = errors.New("pool is empty")
	errUnknownBackend = errors.New("unknown backend")
)

// LoadBalancer picks backends using power-of-two-choices: of the two
// candidates, the one with fewer active connections wins (ties go to the
// lexicographically smaller name).
type LoadBalancer struct {
	mu    sync.Mutex
	conns map[string]int
	names []string // sorted, for deterministic STATUS output
}

func (lb *LoadBalancer) SetPool(names []string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.conns = make(map[string]int, len(names))
	for _, n := range names {
		lb.conns[n] = 0
	}
	lb.names = slices.Sorted(maps.Keys(lb.conns))
}

func (lb *LoadBalancer) Pick(a, b string) (string, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.conns) == 0 {
		return "", errEmptyPool
	}
	ca, okA := lb.conns[a]
	cb, okB := lb.conns[b]
	if !okA || !okB {
		return "", errUnknownBackend
	}
	pick := a
	if cb < ca || (cb == ca && b < a) {
		pick = b
	}
	lb.conns[pick]++
	return pick, nil
}

func (lb *LoadBalancer) Done(name string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	c, ok := lb.conns[name]
	if !ok {
		return errUnknownBackend
	}
	if c > 0 {
		lb.conns[name] = c - 1
	}
	return nil
}

func (lb *LoadBalancer) Status() []string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lines := make([]string, 0, len(lb.names))
	for _, n := range lb.names {
		lines = append(lines, fmt.Sprintf("%s:%d", n, lb.conns[n]))
	}
	return lines
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var lb LoadBalancer
	for sc.Scan() {
		words := strings.Fields(sc.Text())
		if len(words) == 0 {
			continue
		}
		cmd, args := words[0], words[1:]
		switch cmd {
		case "POOL":
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "POOL: need at least one backend")
				continue
			}
			lb.SetPool(args)
			fmt.Fprintln(out, "OK")
		case "PICK":
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "PICK: need exactly two backends")
				continue
			}
			pick, err := lb.Pick(args[0], args[1])
			switch {
			case errors.Is(err, errEmptyPool):
				fmt.Fprintln(out, "EMPTY")
			case err != nil:
				fmt.Fprintln(os.Stderr, "PICK:", err)
			default:
				fmt.Fprintln(out, pick)
			}
		case "DONE":
			if len(args) != 1 {
				fmt.Fprintln(os.Stderr, "DONE: need exactly one backend")
				continue
			}
			if err := lb.Done(args[0]); err != nil {
				fmt.Fprintln(os.Stderr, "DONE:", err)
				continue
			}
			fmt.Fprintln(out, "OK")
		case "STATUS":
			for _, line := range lb.Status() {
				fmt.Fprintln(out, line)
			}
		default:
			fmt.Fprintln(os.Stderr, "unknown command:", cmd)
		}
	}
}
