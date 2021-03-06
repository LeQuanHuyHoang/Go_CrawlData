package main

import (
	"encoding/json"
	"fmt"
	"go_crawler/pkg/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/lithammer/shortuuid/v3"

	"github.com/gocolly/colly"
)

type Variant struct {
	URL string
	IMG string
}

type ProductDetail struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	ProductType  string
	Stock        string `json:"stock"`
	PriceCurrent string `json:"price_current"`
	Price        string `json:"price"`
	SKU          string `json:"sku"`
	IMG          string
	MODEL        string
	VENDOR       string
	SEASON       string
	CATEGORY     string
	DESCRIPTION  string
	VARIANT      []Variant
}

type Category struct {
	ID          string
	Name        string
	URL         string
	HasChildren bool
	ParentId    string
	ParentName  string
}

type CategoryLv1 struct {
	Lv1CategoryID   string `bson.D:"lv1_category_id" json:"lv1_category_id"`
	Lv1CategoryURL  string `bson.D:"lv1_category_url" json:"lv1_category_url"`
	Lv1CategoryName string `bson.D:"lv1_category_url" json:"Lv1_category_name"`
}

type CategoryLv2 struct {
	Lv2CategoryURL  string `json:"lv2_category_url"`
	Lv2CategoryName string `json:"lv2_category_name"`
	Lv1CategoryName string `json:"lv1_category_name"`
}

type CategoryLv3 struct {
	Lv3CategoryURL  string `json:"lv3_category_url"`
	Lv3CategoryName string `json:"lv3_category_name"`
	Lv2CategoryName string `json:"lv2_category_name"`
}

func NewCategory() *Category {
	return &Category{}
}

func NewCategoryLv1() *CategoryLv1 {
	return &CategoryLv1{}
}

func NewCategoryLv2() *CategoryLv2 {
	return &CategoryLv2{}
}

func NewCategoryLv3() *CategoryLv3 {
	return &CategoryLv3{}
}

func NewProductDetail() *ProductDetail {
	return &ProductDetail{}
}

func main() {
	category := NewCategory()
	_, err := category.GetCategory("https://hoang-phuc.com/")
	utils.CheckError(err)

	/* 	category_lv1 := NewCategoryLv1()
	   	lcate_lv1, err := category_lv1.GetCategoryLv1("https://hoang-phuc.com/")
	   	utils.CheckError(err)

	   	category_lv2 := NewCategoryLv2()
	   	lcate_lv2, err := category_lv2.GetCategoryLv2("https://hoang-phuc.com/", lcate_lv1)
	   	utils.CheckError(err)

	   	category_lv3 := NewCategoryLv3()
	   	err = category_lv3.GetCategoryLv3("https://hoang-phuc.com/", lcate_lv2)
	   	utils.CheckError(err) */

	/* 	product_detail := NewProductDetail()
	   	var list_category []CategoryLv3
	   	list_category, err := ReadCategoryJson()
	   	utils.CheckError(err)
	   	err = product_detail.GetProductDetail(list_category)
	   	utils.CheckError(err) */
	log.Printf("Crawling Complete\n")
}

func (cate *Category) GetCategory(url string) ([]Category, error) {
	list_category_lv1 := make([]Category, 0)
	list_category_lv2 := make([]Category, 0)
	list_category_lv3 := make([]Category, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("hoang-phuc.com"),
	)
	c.OnHTML("div.magemenu-menu.horizontal-menu > ul > li", func(e *colly.HTMLElement) {
		category := Category{}
		category.URL = e.ChildAttr("a", "href")
		fmt.Println(category.URL)
		category.ID = e.ChildAttr("a", "class")
		category.Name = e.ChildText("li.menu > a > span > span")
		category.HasChildren = false
		category.ParentId = ""
		category.ParentName = ""
		if category.Name != "B???o h??nh" {
			list_category_lv1 = append(list_category_lv1, category)
		}
	})

	c.OnHTML("div.grid-child > ul > li", func(e *colly.HTMLElement) {
		category := Category{}
		category.ID = shortuuid.New()
		category.URL = e.ChildAttr("a", "href")
		category.Name = e.ChildText("div.grid-child > ul > li> a > span>span")
		for i, v := range list_category_lv1 {
			link := strings.TrimSuffix(v.URL, ".html")
			if strings.Contains(category.URL, link) {
				category.ParentId = v.ID
				category.ParentName = v.Name
				list_category_lv1[i].HasChildren = true
			}
		}
		list_category_lv2 = append(list_category_lv2, category)
	})
	category := Category{}
	c.OnHTML("div.submenu-container > ul > li", func(e *colly.HTMLElement) {
		category.ID = shortuuid.New()
		category.URL = e.ChildAttr("a", "href")
		category.Name = e.ChildText("span > span")
		for i, v := range list_category_lv2 {
			link := strings.TrimSuffix(v.URL, ".html")
			if strings.Contains(category.URL, link) {
				category.ParentId = v.ID
				category.ParentName = v.Name
				list_category_lv2[i].HasChildren = true
			}
		}
		list_category_lv3 = append(list_category_lv3, category)
	})

	c.Visit(url)

	categoryJson1, err := json.MarshalIndent(list_category_lv1, "", "   ")
	utils.CheckError(err)
	err = ioutil.WriteFile("category_lv1.json", categoryJson1, 0644)
	utils.CheckError(err)

	categoryJson2, err := json.MarshalIndent(list_category_lv2, "", "   ")
	utils.CheckError(err)
	err = ioutil.WriteFile("category_lv2.json", categoryJson2, 0644)
	utils.CheckError(err)

	categoryJson3, err := json.MarshalIndent(list_category_lv3, "", "   ")
	utils.CheckError(err)
	err = ioutil.WriteFile("category_lv3.json", categoryJson3, 0644)
	utils.CheckError(err)

	fmt.Println(url)
	return nil, nil
}

func (cate *CategoryLv1) GetCategoryLv1(url string) ([]CategoryLv1, error) {
	list_category := make([]CategoryLv1, 0)
	category := &CategoryLv1{}
	c := colly.NewCollector(
		colly.AllowedDomains("hoang-phuc.com"),
	)
	c.OnHTML("div.magemenu-menu.horizontal-menu > ul > li", func(e *colly.HTMLElement) {
		category.Lv1CategoryURL = e.ChildAttr("a", "href")
		category.Lv1CategoryName = e.ChildText("li.menu > a > span > span")
		category.Lv1CategoryID = shortuuid.New()
		if category.Lv1CategoryName != "B???o h??nh" {
			list_category = append(list_category, *category)
		}
	})
	c.Visit(url)

	categoryJson, err := json.Marshal(list_category)
	utils.CheckError(err)
	err = ioutil.WriteFile("category_lv1.json", categoryJson, 0644)
	utils.CheckError(err)
	log.Println(c)

	/* 	categoryBson, err := bson.Marshal(bson.M{"data": list_category})
	   	fmt.Println(categoryBson)
	   	utils.CheckError(err)
	   	test := []interface{}{categoryBson}
	   	fmt.Println(test...)
	   	client, err := database.GetMongoClient()
	   	utils.CheckError(err)
	   	//Create a handle to the respective collection in the database.
	   	collection := client.Database(database.DB).Collection("CategoryLv1")
	   	//Perform InsertOne operation & validate against the error.
	   	_, err = collection.InsertMany(context.TODO(), test)
	   	utils.CheckError(err) */

	return list_category, nil
}

func (cate *CategoryLv2) GetCategoryLv2(url string, lv1 []CategoryLv1) ([]CategoryLv2, error) {
	list_category := make([]CategoryLv2, 0)
	category := &CategoryLv2{}
	c := colly.NewCollector(
		colly.AllowedDomains("hoang-phuc.com"),
	)
	c.OnHTML("div.grid-child > ul > li", func(e *colly.HTMLElement) {
		category.Lv2CategoryURL = e.ChildAttr(" a", "href")
		category.Lv2CategoryName = e.ChildText("div.grid-child > ul > li> a > span>span")
		for _, v := range lv1 {
			link := strings.TrimSuffix(v.Lv1CategoryURL, ".html")
			if strings.Contains(category.Lv2CategoryURL, link) {
				category.Lv1CategoryName = v.Lv1CategoryName
			}
		}
		list_category = append(list_category, *category)
	})
	c.Visit(url)

	categoryJson, err := json.Marshal(list_category)
	utils.CheckError(err)
	err = ioutil.WriteFile("category_lv2.json", categoryJson, 0644)
	utils.CheckError(err)
	log.Println(c)
	return list_category, nil
}

func (cate *CategoryLv3) GetCategoryLv3(url string, lv2 []CategoryLv2) error {
	list_category := make([]CategoryLv3, 0)
	category := &CategoryLv3{}
	c := colly.NewCollector(
		colly.AllowedDomains("hoang-phuc.com"),
	)
	c.OnHTML("div.submenu-container > ul > li", func(e *colly.HTMLElement) {
		category.Lv3CategoryURL = e.ChildAttr("a", "href")
		category.Lv3CategoryName = e.ChildText("span > span")
		for _, v := range lv2 {
			link := strings.TrimSuffix(v.Lv2CategoryURL, ".html")
			if strings.Contains(category.Lv3CategoryURL, link) {
				category.Lv2CategoryName = v.Lv2CategoryName
			}
		}
		list_category = append(list_category, *category)
	})
	c.Visit(url)

	categoryJson, err := json.Marshal(list_category)
	utils.CheckError(err)
	err = ioutil.WriteFile("category_lv3.json", categoryJson, 0644)
	utils.CheckError(err)
	log.Println(c)
	return nil
}

func (p *ProductDetail) GetProductDetail(cate []CategoryLv3) error {
	listproducts := make([]ProductDetail, 0)
	product_detail := &ProductDetail{}
	related_product := Variant{}

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
		productInfo.OnHTML("div.gallery-placeholder > img", func(e *colly.HTMLElement) {
			product_detail.IMG = e.Attr("data-original")
		})

		productInfo.OnHTML("div.product.info.detailed", func(e *colly.HTMLElement) {
			product_detail.MODEL = e.ChildText("td.col.data[data-th='Model']")
			product_detail.CATEGORY = e.ChildText("td.col.data[data-th='Gi???i t??nh']")
			product_detail.SEASON = e.ChildText("td.col.data[data-th='Season']")
			product_detail.VENDOR = e.ChildText("td.col.data[data-th='Th????ng hi???u']")
			product_detail.DESCRIPTION = e.ChildText("div.value.description-item > p")
		})

		productInfo.OnHTML("div.product-info-main", func(e *colly.HTMLElement) {
			product_detail.ID = e.ChildAttr("div.price-box.price-final_price", "data-product-id")
			product_detail.Name = e.ChildText("h1.page-title")
			product_detail.Stock = e.ChildAttr("div.stock.available", "title")
			product_detail.PriceCurrent = e.ChildText("span.normal-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price")
			product_detail.Price = e.ChildText("span.old-price.sly-old-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price")
			product_detail.SKU = e.ChildText("div.product.attribute.sku > div.value")

			e.ForEach("ul.related-product-list > li", func(_ int, h *colly.HTMLElement) {
				related_product.URL = h.ChildAttr("a", "href")
				related_product.IMG = h.ChildAttr("img", "data-original")
				product_detail.VARIANT = append(product_detail.VARIANT, related_product)
			})

			listproducts = append(listproducts, *product_detail)
			product_detail.VARIANT = nil
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
		productJson, err := json.MarshalIndent(listproducts, "", "   ")
		utils.CheckError(err)
		err = ioutil.WriteFile("data.json", productJson, 0644)
		utils.CheckError(err)
		log.Println(c)
	}
	return nil
}

func ReadCategoryJson() ([]CategoryLv3, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("category_lv3.json")
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
