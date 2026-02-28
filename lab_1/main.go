package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {

	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Incorrect amount of arguments expected 1 got", len(args))
		os.Exit(1)
	}

	num, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Printf("parsing nummber error: %s", args[0])
		os.Exit(1)
	}

	var isPositive, isNatural, isEven bool
	isPositive = num > 0.0
	isNatural = math.Floor(num) == num
	if isNatural {
		isEven = int(num)%2 == 0
	}

	fmt.Printf("Number %s is:\n", args[0])
	if isPositive {
		fmt.Println("  Positive")
	} else {
		fmt.Println("  Negative")
	}

	if isNatural {
		fmt.Println("  Natural")
		if isEven {
			fmt.Println("  Even")
		} else {
			fmt.Println("  Odd")
		}
	} else {
		fmt.Println("  Real")
	}
}
