package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to hash bruteforce decoder")

	for {
		fmt.Println("Print our hash (print 0 to stop):")
		input, _ := reader.ReadString('\n')
		input = input[0:len(input)-2] + ""

		if input == "0" {
			fmt.Println("Closing program")
			os.Exit(0)
		}

		r, _ := regexp.Compile("[a-z0-9]{64}")
		if r.MatchString(input) {
			fmt.Println("Choose number of threads")
			num, _ := reader.ReadString('\n')
			num = num[0:len(num)-2] + ""
			n, err := strconv.Atoi(num)
			if err != nil || n < 0 {
				fmt.Println("wrong thread number")
			}

			var wg sync.WaitGroup
			for _, v := range divide(validCombinations(5), n) {
				wg.Add(1)
				go func(pass []string, input string) {
					for _, fetch := range pass {
						bruteForce(input, fetch)
					}
					wg.Done()
				}(v, input)
				//go decode(v, input)
			}
			wg.Wait()
		} else {
			fmt.Println("hash is not in SHA-256 standard")
			continue
		}
	}
}

func bruteForce(hash string, fetch string) {
	h := sha256.New()
	h.Write([]byte(fetch))

	if fmt.Sprintf("%x", h.Sum(nil)) == hash {
		fmt.Println("fetched password: " + fetch)
	}
}

func validCombinations(maxChar int) []string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	var result []string = []string{""}
	var last, next int
	for i := 1; i <= maxChar; i++ {
		last = len(result)
		for _, str := range result[next:] {
			for _, char := range alphabet {
				result = append(result, str+string(char))
			}
		}
		next = last
	}
	return result
}

func divide(logs []string, numCPU int) [][]string {
	var divided [][]string

	chunkSize := (len(logs) + numCPU - 1) / numCPU

	for i := 0; i < len(logs); i += chunkSize {
		end := i + chunkSize

		if end > len(logs) {
			end = len(logs)
		}

		divided = append(divided, logs[i:end])
	}

	return divided
}
