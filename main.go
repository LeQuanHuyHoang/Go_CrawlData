package main

import (
	_ "encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	_ "os"

	"github.com/gocolly/colly"
)

type Product struct {
	URL          string `json:"URL"`
	Name         string `json:"Name"`
	Stock        string `json:"Stock"`
	PriceCurrent string `json:"PriceCurrent"`
	Price        string `json:"Price"`
	SKU          string `json:"SKU"`
}

type Products struct {
	List []Product `json:"Products"`
}

func NewProducts() *Products {
	return &Products{}
}

func main() {
	/* 	fName := "data.csv"
	   	file, err := os.Create(fName)
	   	if err != nil {
	   		log.Fatalf("Could not create file, err :%kq", err)
	   		return
	   	}
	   	defer file.Close()

	   	writer := csv.Writer()
	   	defer writer.Flush() */
	product := NewProducts()
	err := product.CrawlData()
	checkError(err)
	productJson, err := json.Marshal(product)
	checkError(err)
	err = ioutil.WriteFile("data.json", productJson, 0644)
	checkError(err)
	log.Printf("Crawling Complete\n")
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (products *Products) CrawlData() error {
	c := colly.NewCollector(
		colly.AllowedDomains("hoang-phuc.com"),
	)

	productInfo := c.Clone()

	c.OnHTML("li.item.product.product-item", func(e *colly.HTMLElement) {
		e.ForEach("a.product.photo.product-item-photo", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("href")
			productInfo.Visit(e.Request.AbsoluteURL(link))
		})
	})
	productInfo.OnHTML("div.product-info-main", func(e *colly.HTMLElement) {
		pro := Product{
			URL:   e.Attr("li.product-item > a.href"),
			Name:  e.ChildText("h1.page-title"),
			Stock: e.ChildText("div.stock available"),
			Price: e.ChildText("span.price"),
			SKU:   e.ChildText("div.value"),
		}
		products.List = append(products.List, pro)
	})

	c.OnHTML("li.item.pages-item-next > a.action.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit("https://hoang-phuc.com/nam/giay-nam/giay-the-thao.html")
	log.Println(c)
	/* 	for i := 1; i <= TotalPages; i++ {
		fmt.Printf("Crawler Page: %d\n", i)

		c.Visit("https://hoang-phuc.com/nam/giay-nam/giay-the-thao.html?p=" + strconv.Itoa(i))

	} */

	return nil
}
