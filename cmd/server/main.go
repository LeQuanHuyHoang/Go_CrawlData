package main

import (
	_ "encoding/csv"
	"encoding/json"
	"fmt"
	"go_crawler/pkg/utils"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type ProductDetail struct {
	URL          string `json:"url"`
	Name         string `json:"name"`
	Stock        string `json:"stock"`
	PriceCurrent string `json:"price_current"`
	Price        string `json:"price"`
	SKU          string `json:"sku"`
}

type CategoryLv1 struct {
	Lv1CategoryURL  string `json:"lv1_category_url"`
	Lv1CategoryName string `json:"Lv1_category_name"`
}

type CategoryLv2 struct {
	Lv2CategoryURL  string `json:"lv2_category_url"`
	Lv2CategoryName string `json:"lv2_category_name"`
}

type CategoryLv3 struct {
	Lv3CategoryURL  string `json:"lv3_category_url"`
	Lv3CategoryName string `json:"lv3_category_name"`
}

func NewCategoryLv1() *CategoryLv1 {
	return &CategoryLv1{}
}

func NewCategoryLv2() *CategoryLv2 {
	return &CategoryLv2{}
}

func NewCategory() *CategoryLv3 {
	return &CategoryLv3{}
}

func NewProductDetail() *ProductDetail {
	return &ProductDetail{}
}

/*
type Products struct {
	List []ProductDetail `json:"product_detail"`
}
*/
/* func NewProducts() *Products {
	return &Products{}
} */

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
	/* 	product := NewProducts()
	   	err := product.CrawlData()
	   	utils.CheckError(err)
	   	productJson, err := json.Marshal(product)
	   	checkError(err)
	   	err = ioutil.WriteFile("data.json", productJson, 0644)
	   	checkError(err)
	   	log.Printf("Crawling Complete\n") */
	category := NewCategory()
	product_detail := NewProductDetail()
	err := category.GetCategoryLv3("https://hoang-phuc.com/")
	utils.CheckError(err)
	var list_category []CategoryLv3
	list_category, err = category.ReadCategoryJson()
	utils.CheckError(err)
	err = product_detail.GetProductDetail(list_category)
	utils.CheckError(err)
	log.Println()
	log.Printf("Crawling Complete\n")
}

func (cate *CategoryLv3) GetCategoryLv3(url string) error {
	list_category := make([]CategoryLv3, 0)
	category := &CategoryLv3{}
	c := colly.NewCollector(
		colly.AllowedDomains("hoang-phuc.com"),
	)
	c.OnHTML("div.submenu-container.grid > ul > li", func(e *colly.HTMLElement) {
		category.Lv3CategoryURL = e.ChildAttr("a", "href")
		category.Lv3CategoryName = e.ChildText("span > span")
		list_category = append(list_category, *category)
	})
	c.Visit(url)

	categoryJson, err := json.Marshal(list_category)
	utils.CheckError(err)
	err = ioutil.WriteFile("category.json", categoryJson, 0644)
	utils.CheckError(err)
	log.Println(c)
	return nil
}

func (p *ProductDetail) GetProductDetail(cate []CategoryLv3) error {
	listproducts := make([]ProductDetail, 0)
	product_detail := &ProductDetail{}
	for _, v := range cate {
		url := v.Lv3CategoryURL
		fmt.Println(url)
		c := colly.NewCollector(
			colly.AllowedDomains("hoang-phuc.com"),
		)

		productInfo := c.Clone()

		c.OnHTML("li.item.product.product-item", func(e *colly.HTMLElement) {
			e.ForEach("a.product.photo.product-item-photo", func(_ int, el *colly.HTMLElement) {
				product_detail.URL = el.Attr("href")
				productInfo.Visit(e.Request.AbsoluteURL(product_detail.URL))
			})
		})
		productInfo.OnHTML("div.product-info-main", func(e *colly.HTMLElement) {
			product_detail.Name = e.ChildText("h1.page-title")
			product_detail.Stock = e.ChildAttr("div.stock.available", "title")
			product_detail.PriceCurrent = e.ChildText("span.normal-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price")
			product_detail.Price = e.ChildText("span.old-price.sly-old-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price")
			product_detail.SKU = e.ChildText("div.product.attribute.sku > div.value")

			listproducts = append(listproducts, *product_detail)
		})

		c.OnHTML("li.item.pages-item-next > a.action.next", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			fmt.Println(link)
			c.Visit(e.Request.AbsoluteURL(link))
		})
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting: ", r.URL.String())
		})
		c.Visit(url)
		productJson, err := json.Marshal(listproducts)
		utils.CheckError(err)
		err = ioutil.WriteFile("data.json", productJson, 0644)
		utils.CheckError(err)
		log.Println(c)
	}
	return nil
}

func (cate *CategoryLv3) ReadCategoryJson() ([]CategoryLv3, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("category.json")
	// if we os.Open returns an error then handle it
	utils.CheckError(err)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var rs []CategoryLv3
	json.Unmarshal([]byte(byteValue), &rs)

	return rs, nil
}

/* func (products *Products) CrawlData() error {
	c := colly.NewCollector(
		colly.AllowedDomains("hoang-phuc.com"),
	)

	productInfo := c.Clone()
	pro := &Product{}

	c.OnHTML("li.item.product.product-item", func(e *colly.HTMLElement) {
		e.ForEach("a.product.photo.product-item-photo", func(_ int, el *colly.HTMLElement) {
			pro.URL = el.Attr("href")
			productInfo.Visit(e.Request.AbsoluteURL(pro.URL))
		})
	})
	productInfo.OnHTML("div.product-info-main", func(e *colly.HTMLElement) {
		pro.Name = e.ChildText("h1.page-title")
		pro.Stock = e.ChildAttr("div.stock.available", "title")
		pro.PriceCurrent = e.ChildText("span.normal-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price")
		pro.Price = e.ChildText("span.old-price.sly-old-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price")
		pro.SKU = e.ChildText("div.product.attribute.sku > div.value")

		products.List = append(products.List, *pro)
	})

	c.OnHTML("li.item.pages-item-next > a.action.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.Visit("https://hoang-phuc.com/nam/giay-nam/giay-the-thao.html")
	log.Println(c)


	return nil
}
*/
/* 	for i := 1; i <= TotalPages; i++ {
	fmt.Printf("Crawler Page: %d\n", i)

	c.Visit("https://hoang-phuc.com/nam/giay-nam/giay-the-thao.html?p=" + strconv.Itoa(i))

} */
