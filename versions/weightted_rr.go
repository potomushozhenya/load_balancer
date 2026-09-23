package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Backend struct {
	name        string
	init_weight int
	weight      int
}

type LoadBalancer struct {
	pool     []Backend
	pool_len int
	mutex    sync.Mutex
}

// TODO (round-robin): implement per the lesson description.
func RR_pick(load_balancer *LoadBalancer) string {
	if load_balancer.pool_len == 0 {
		return "EMPTY"
	}
	max, max_index, sum := 0, 0, 0
	load_balancer.mutex.Lock()
	for i := 0; i < load_balancer.pool_len; i++ {
		new_value := load_balancer.pool[i].weight + load_balancer.pool[i].init_weight
		sum += new_value
		load_balancer.pool[i].weight = new_value
		if new_value > max {
			max = new_value
			max_index = i
		}
	}
	load_balancer.pool[max_index].weight -= sum
	load_balancer.mutex.Unlock()
	return load_balancer.pool[max_index].name
}

func RR_pool(load_balancer *LoadBalancer, pool_strings []string) {
	pool_len := len(pool_strings)
	if pool_len < 1 {
		fmt.Println("pool can't be empty")
	}
	load_balancer.mutex.Lock()
	load_balancer.pool = make([]Backend, pool_len)
	for i := 0; i < pool_len; i++ {
		pool_strings_splitted := strings.Split(pool_strings[i], ":")
		parsed_weight, err := strconv.Atoi(pool_strings_splitted[1])
		if err != nil {
			fmt.Println(err)
			continue
		}
		load_balancer.pool[i].name = pool_strings_splitted[0]
		load_balancer.pool[i].init_weight = parsed_weight
	}
	load_balancer.pool_len = pool_len

	fmt.Println("OK")
	load_balancer.mutex.Unlock()
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	lb := LoadBalancer{[]Backend{}, 0, sync.Mutex{}}
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
		case "PICKN":
			{
				n, err := strconv.Atoi(args[0])
				if err != nil {
					fmt.Println(err)
					continue
				}
				picks := make([]string, n, n)
				for i := 0; i < n; i++ {
					picks[i] = RR_pick(&lb)
				}
				fmt.Println(strings.Join(picks, ","))

			}
		case "POOL":
			{
				RR_pool(&lb, args)
			}
		}
	}
}
