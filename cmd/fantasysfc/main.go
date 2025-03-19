package main

import (
	"fmt"

	"github.com/jasonmcclafferty/fantasySFC/internal/scraper"
)

func main() {
	fmt.Println("Dia duit, a shaol!")

	fixture, res, err := scraper.Scrape()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(fixture)
	fmt.Println(res)
}
