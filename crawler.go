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
	"unicode"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
				"strings"
	"net/http"
	"io/ioutil"
	"github.com/grokify/html-strip-tags-go"
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
	total_visit_page := 10
	total_results := 0

	article_urls := ""

	google_url := "https://www.google.com.vn/search?q=%s&tbs=qdr:%s&start=%d"
	for page:=0; page < total_visit_page; page++ {
		c := colly.NewCollector() // colly.Debugger(&debug.LogDebugger{})
		count := 1
		// Find and visit all links
		c.OnHTML("h3[class='r']", func(e *colly.HTMLElement) {
			//e.Request.Visit(e.Attr("href"))
			log.Println("Result ", count, ":")
			count++
			log.Println(e.Text)
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
					resp, err := http.Get(full_url)
					if err != nil {
						return
					}
					defer resp.Body.Close()
					log.Println(resp.Body)
					body, _ := ioutil.ReadAll(resp.Body)
					bs := string(body)
					stripped := normalize_string(bs)
					log.Println(stripped)
					kw := normalize_string(keyword.Keyword)
					log.Println(kw)
					ct := strings.Contains(stripped, kw)
					log.Println(ct)
					if ct {
						total_results += 1
						models.SaveCrawled(keyword.Id, full_url, e.Text, ele.Text(), ele.Text())
						article_urls += fmt.Sprintf("Article : %s <br/>"+
							"Url: %s<br/>----------------------------------------<br/><br/>", ele.Text(), full_url)
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


	if total_results > 0 {
		email_body := fmt.Sprintf("<html>Hello there, <br/>"+
			"We have found %d new articles about keyword : %s , on GOOGLE search engine <br/>"+
			"List of articles: <br/><br/>"+
			"%s<br/></html>", total_results, keyword.Keyword, article_urls)
		emails, _ := models.EmailForKeyword(keyword.Id)
		for _, to := range emails {
			SendEmail(string(total_results)+" new articles crawled for keyword: "+keyword.Keyword, email_body, to)
		}
		log.Println(total_results)
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

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func remove_unicode_accent(txt string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	b := make([]byte, len(txt))
	t.Transform(b, []byte(txt), true)

	rs := strings.Replace(string(b), "đ", "d", -1)
	rs =  strings.Replace(rs, "Đ", "D", -1)
	return rs
}

func remove_spaces(txt string) string {
	pattern, _ := regexp.Compile(`([\s\r\n])+`)
	final := pattern.ReplaceAllString(txt, "")
	final = strings.Replace(final, "&nbsp;", "", -1)
	return final
}

func normalize_string(txt string) string {
	txt = strings.ToLower(remove_spaces(remove_unicode_accent(strip.StripTags(txt))))
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	txt = reg.ReplaceAllString(txt, "")
	return txt
}