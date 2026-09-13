package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Backend struct {
	name string
}

type LoadBalancer struct {
	pool     []Backend
	pool_len int
	counter  int
	mutex    sync.Mutex
}

// TODO (round-robin): implement per the lesson description.
func RR_pick(load_balancer *LoadBalancer) string {
	if load_balancer.pool_len == 0 {
		return "EMPTY"
	}
	load_balancer.mutex.Lock()
	index := load_balancer.counter % load_balancer.pool_len
	load_balancer.counter++
	load_balancer.mutex.Unlock()
	return load_balancer.pool[index].name
}

func RR_pool(load_balancer *LoadBalancer, pool_strings []string) {
	pool_len := len(pool_strings)
	if pool_len < 1 {
		fmt.Println("pool can't be empty")
	}
	load_balancer.mutex.Lock()
	load_balancer.pool = make([]Backend, pool_len)
	for i := 0; i < pool_len; i++ {
		load_balancer.pool[i].name = pool_strings[i]
	}
	load_balancer.pool_len = pool_len
	load_balancer.counter = 0
	fmt.Println("OK")
	load_balancer.mutex.Unlock()
}

func RR_reset(load_balancer *LoadBalancer) {
	load_balancer.mutex.Lock()
	load_balancer.counter = 0
	fmt.Println("OK")
	load_balancer.mutex.Unlock()
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	lb := LoadBalancer{[]Backend{}, 0, 0, sync.Mutex{}}
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
		case "POOL":
			{
				RR_pool(&lb, args)
			}
		case "RESET":
			{
				RR_reset(&lb)
			}
		}
	}
}
