package main

import "fmt"

func primeNumDecom(n int) {
	k := 2
	for n > 1 {
		if n%k == 0 {
			fmt.Println(k)
			n = n / k
		} else {
			k = k + 1
		}
	}
}

func main() {
	var n int
	fmt.Scan(&n)

	primeNumDecom(n)

}
