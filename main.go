package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
)

type Backend struct {
	name        string
	connections int
}

type LoadBalancer struct {
	pool        map[string]*Backend
	sorted_keys []string
	pool_len    int
	mutex       sync.Mutex
}

// TODO (round-robin): implement per the lesson description.
func RR_pick(load_balancer *LoadBalancer) string {
	if load_balancer.pool_len == 0 {
		return "EMPTY"
	}
	min, min_key := math.MaxInt, ""
	load_balancer.mutex.Lock()
	for i := 0; i < load_balancer.pool_len; i++ {
		key := load_balancer.sorted_keys[i]
		curr_connections := load_balancer.pool[key].connections
		if curr_connections < min {
			min = curr_connections
			min_key = key
		}
	}
	load_balancer.pool[min_key].connections += 1
	load_balancer.mutex.Unlock()
	return load_balancer.pool[min_key].name
}

func RR_done(load_balancer *LoadBalancer, done_strings []string) {
	done_len := len(done_strings)
	if done_len != 1 {
		fmt.Println("wrong format")
	}
	load_balancer.mutex.Lock()
	backend, ok := load_balancer.pool[done_strings[0]]
	cur_value := -1
	if ok {
		cur_value = backend.connections
	} else {
		fmt.Println("wrong backend name")
		load_balancer.mutex.Unlock()
		return
	}
	if cur_value > 0 {
		backend.connections = cur_value - 1
	}
	fmt.Println("OK")
	load_balancer.mutex.Unlock()
}

func RR_pool(load_balancer *LoadBalancer, pool_strings []string) {
	pool_len := len(pool_strings)
	if pool_len < 1 {
		fmt.Println("pool can not be empty")
	}
	load_balancer.mutex.Lock()
	load_balancer.pool = make(map[string]*Backend)
	load_balancer.sorted_keys = make([]string, pool_len)
	for i := 0; i < pool_len; i++ {
		name := pool_strings[i]
		load_balancer.sorted_keys[i] = name
		load_balancer.pool[name] = &Backend{name: name, connections: 0}
	}
	load_balancer.pool_len = pool_len
	fmt.Println("OK")
	load_balancer.mutex.Unlock()
}

func RR_status(load_balancer *LoadBalancer) {
	load_balancer.mutex.Lock()
	for i := 0; i < load_balancer.pool_len; i++ {
		key := load_balancer.sorted_keys[i]
		fmt.Printf("%s:%d\n", key, load_balancer.pool[key].connections)
	}
	load_balancer.mutex.Unlock()
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	lb := LoadBalancer{nil, nil, 0, sync.Mutex{}}
	for sc.Scan() {
		words := strings.Fields(sc.Text())
		if len(words) == 0 {
			fmt.Println("Incorrect line")
			continue
		}
		cmd, args := words[0], words[1:]
		switch cmd {
		case "PICK":
			{
				pick := RR_pick(&lb)
				fmt.Println(pick)

			}
		case "DONE":
			{
				RR_done(&lb, args)

			}
		case "POOL":
			{
				RR_pool(&lb, args)
			}
		case "STATUS":
			{
				RR_status(&lb)
			}
		}
	}
}
