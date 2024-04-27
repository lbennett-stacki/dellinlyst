package crawler

import (
	"fmt"
	"net/http"
	"strconv"
)

func fetchPageBody(url string, page int, productsChannel chan<- []Product, errorChannel chan<- error) {
	response, err := http.Get(fmt.Sprintf("%s?page=%s", url, strconv.Itoa(page)))
	if err != nil {
		errorChannel <- err
		return
	}
	defer response.Body.Close()

	products, err := parsePageProducts(response.Body)
	if err != nil {
		errorChannel <- err
		return
	}

	productsChannel <- products
}
