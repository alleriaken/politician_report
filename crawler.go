package main

import (
	"github.com/joho/godotenv"
	"./models"
	"github.com/gocolly/colly"
	"fmt"
	"net/url"
	"regexp"
	"time"
	"log"
	"net/smtp"
	"os"
)


func main() {

	godotenv.Load()
	models.InitDB()
	keywords, _ := models.AllKeyword()
	for _, kw := range keywords {
		crawlGoogle(kw)
		duration := time.Duration(10) * time.Second
		time.Sleep(duration)
	}
	//crawlGoogle()

	defer models.CloseDB()
}

func crawlGoogle(keyword *models.Keyword)  {
	total_visit_page := 5
	google_url := "https://www.google.com.vn/search?q=%s&tbs=qdr:%s&start=%d"
	for page:=0; page < total_visit_page; page++ {
		c := colly.NewCollector() // colly.Debugger(&debug.LogDebugger{})
		count := 1
		// Find and visit all links
		c.OnHTML("h3[class='r']", func(e *colly.HTMLElement) {
			//e.Request.Visit(e.Attr("href"))
			log.Println("Result ", count, ":")
			count++
			ele := e.DOM.Children().First()
			log.Println(ele.Text())
			url_pattern := regexp.MustCompile("/url\\?q=(https?)://(.*?)&sa=")
			href, _ := ele.Attr("href")
			matches := url_pattern.FindStringSubmatch(href)
			if len(matches) > 2 {
				protocol := matches[1]
				article_url := matches[2]
				full_url := fmt.Sprintf("%s://%s", protocol, article_url)
				log.Println(full_url)
				crawled := models.CheckCrawledUrl(keyword.Id, full_url)
				if crawled == false {
					models.SaveCrawled(keyword.Id, full_url, e.Text, ele.Text(), ele.Text())
					email_body := fmt.Sprintf("<html>Hello there, <br/>" +
						"We have found a new article about keyword : %s , on GOOGLE search engine <br/>" +
						"Please check it out: <br/>" +
						"Article : %s <br/></html>" +
						"Url: %s", keyword.Keyword, ele.Text(), full_url)
					emails, _ := models.EmailForKeyword(keyword.Id)
					for _, to := range emails {
						SendEmail("New article crawled for keyword: "+keyword.Keyword, email_body, to)
					}
				}
			} else {
				log.Println(matches)
			}

		})

		c.OnRequest(func(r *colly.Request) {
			log.Println("Visiting", r.URL)
		})

		url_path := fmt.Sprintf(google_url, url.PathEscape(keyword.Keyword), keyword.TimeOfArticle, page * 10)

		c.Visit(url_path)
		duration := time.Duration(10) * time.Second
		time.Sleep(duration)
	}
}

func SendEmail(title string, body string, to string) {
	from := os.Getenv("SENDER_EMAIL")
	pass := os.Getenv("SENDER_EMAIL_PASSWORD")

	msg := "From: %s\n" +
		"To: %s\n" +
		"Subject: %s\n" +
		"Mime-Version: 1.0;\n" +
		"Content-Type: text/html; charset=\"ISO-8859-1\";\n" +
		"Content-Transfer-Encoding: 7bit;\n\n" +
		"%s"

	msg_body := fmt.Sprintf(msg, from, to, title, body)

	log.Println(msg_body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg_body))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
}