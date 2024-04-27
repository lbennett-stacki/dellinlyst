package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var ErrNoProductsFound = fmt.Errorf("no products found")

func crawl(url string) ([]Product, error) {
	page := 1
	products := []Product{}

	for {
		productsChannel := make(chan []Product)
		errorChannel := make(chan error)
		go fetchPageBody(url, page, productsChannel, errorChannel)

		fmt.Println("fetching products on page", page)

		select {
		case newProducts := <-productsChannel:
			products = append(products, newProducts...)
			page++
		case err := <-errorChannel:
			if errors.Is(err, ErrNoProductsFound) {
				return products, nil
			}
			fmt.Println("Error fetching page:", err)
			return nil, err
		}
	}
}

func writeOutput(products []Product) {
	fmt.Println("Writing products:", len(products))

	file, err := os.Create(fmt.Sprintf("outputs/%d.json", time.Now().Unix()))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(products); err != nil {
		fmt.Println("Error encoding products:", err)
		return
	}

	fmt.Println("Products output saved")
}

func CrawlDelli() {
	url := "https://delli.market/collections/shop-all"

	products, err := crawl(url)
	if err != nil {
		fmt.Println("Error crawling Delli:", err)
		return
	}

	writeOutput(products)
}
