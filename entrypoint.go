package main

import "fmt"

func main() {
	fmt.Println("Dia duit, a shaol!")

	fixture, res, err := scrape()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(fixture)
	fmt.Println(res)
}
