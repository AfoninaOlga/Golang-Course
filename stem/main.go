package main

import "fmt"

func main() {
	input := ParseInput()

	if input == "" {
		fmt.Println("Use \"-s\" flag to pass a string for stemming")
	}

	res, err := StemInput(input)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
