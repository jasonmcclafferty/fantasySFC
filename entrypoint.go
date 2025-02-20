package main

import "fmt"

func main() {
	fmt.Println("Dia duit, a shaol!")

	res, err := scrape()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(res)
}
