package main

import "fmt"

func main() {
	input := ParseInput()

	if input == "" {
		fmt.Println("Use \"-s\" flag to pass a string for stemming")
	}

	res := StemInput(input)

	fmt.Println(res)
}
