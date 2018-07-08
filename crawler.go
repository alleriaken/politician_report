package main

import (
	"github.com/joho/godotenv"
	"./models"
	"github.com/gocolly/colly"
	"fmt"
	"net/url"
	"regexp"
)


func main() {

	godotenv.Load()
	models.InitDB()
	models.AllKeyword()
	//crawlGoogle()

	defer models.CloseDB()
}

func crawlGoogle()  {
	c := colly.NewCollector() // colly.Debugger(&debug.LogDebugger{})

	google_url := "https://www.google.com.vn/search?q=%s&tbs=qdr:%s"
	count := 0
	// Find and visit all links
	c.OnHTML("h3[class='r']", func(e *colly.HTMLElement) {
		//e.Request.Visit(e.Attr("href"))
		fmt.Println("Result ", count, ":")
		count++
		ele := e.DOM.Children().First()
		fmt.Println(ele.Text())
		url_pattern := regexp.MustCompile("/url\\?q=(https?)://(.*?)&sa=")
		href, _ := ele.Attr("href")
		matches := url_pattern.FindStringSubmatch(href)
		fmt.Println( matches)
		if len(matches) > 2 {
			protocol := matches[1]
			article_url := matches[2]
			fmt.Println(fmt.Sprintf("%s://%s", protocol, article_url))
		} else {
			fmt.Println(matches)
		}

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited")
	})

	url_path := fmt.Sprintf(google_url, url.PathEscape("Ho Duc Phoc"), "w" )

	c.Visit(url_path)
}
