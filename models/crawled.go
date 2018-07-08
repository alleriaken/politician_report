package models

import (
	"fmt"
	"log"
	"encoding/json"
)

type Crawled struct {
	Id, CrawledAt, KeywordId, Negative, Positive  int
	Url, Host, Title, GoogleTitle, PreviewContent string
}

func AllCrawled(kw int) ([]*Crawled, error) {
	q := "SELECT * FROM crawled"
	if kw == 0 {
		q = fmt.Sprintf("SELECT * FROM crawled WHERE keyword_id = %d", kw)
	}

	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var crawleds = make([]*Crawled, 0)

	if err = rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	for rows.Next() {
		var crawled Crawled
		if err := rows.Scan(&crawled.Id, &crawled.Url, &crawled.Host,
			&crawled.Title, &crawled.CrawledAt, &crawled.KeywordId,
			&crawled.Negative, &crawled.Positive, &crawled.GoogleTitle,
			&crawled.PreviewContent); err != nil {
			log.Fatal(err)
		}
		log.Println(crawled.Title)
		out, _ := json.Marshal(crawled)
		log.Println(string(out))
		crawleds = append(crawleds, &crawled)
	}

	return crawleds, nil
}

func CheckCrawledUrl(kw_id int, url string) (bool) {
	q := fmt.Sprintf("SELECT * FROM crawled WHERE keyword_id = %d AND url = '%s'", kw_id, url)
	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for rows.Next()  {
		return true
	}

	return false
}

func SaveCrawled(kw_id int, url string, preview_content string, google_title string, title string) (int64) {
	query := "INSERT INTO crawled (keyword_id, url, preview_content, google_title, title) " +
		"VALUES (?, ?, ? ,?, ?)"
	stmt, _ := db.Prepare(query)
	res, _ := stmt.Exec(kw_id, url, preview_content, google_title, title)
	id, _ := res.LastInsertId()
	return id
}
