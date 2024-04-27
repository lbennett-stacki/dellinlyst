package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parsePageProducts(body io.ReadCloser) ([]Product, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println("error parsing body:", err)
		return nil, err
	}

	title := strings.TrimSpace(doc.Find("#ProductGridContainer .title").Text())
	if strings.HasPrefix(title, "No products found") {
		return nil, ErrNoProductsFound
	}

	var products []Product

	doc.Find("#ProductGridContainer quick-buy").Each(func(i int, selection *goquery.Selection) {
		productData, exists := selection.Attr("data-quick-buy")
		if exists {
			var product struct {
				Product Product `json:"product"`
			}
			if err := json.Unmarshal([]byte(productData), &product); err == nil {
				products = append(products, product.Product)
			} else {
				fmt.Println("errorChannel unmarshalling product data:", err)
			}
		}
	})

	return products, nil
}
