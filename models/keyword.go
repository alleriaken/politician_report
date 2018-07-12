package models

import (
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
)

type Keyword struct {
	Id, LastCrawledAt, CrawlDelayTime int
	TimeOfArticle, Keyword            string
}

func AllKeyword() ([]*Keyword, error) {
	rows, err := db.Query("SELECT * FROM keywords")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var words = make([]*Keyword, 0)

	if err = rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	for rows.Next() {
		var kw Keyword
		if err := rows.Scan(&kw.Id, &kw.Keyword, &kw.TimeOfArticle, &kw.LastCrawledAt, &kw.CrawlDelayTime); err != nil {
			log.Fatal(err)
		}
		log.Println(kw.Keyword)
		out, _ := json.Marshal(kw)
		log.Println(string(out))
		words = append(words, &kw)
	}

	return words, nil
}

func EmailForKeyword(keyword_id int) ([]string, error) {
	q := fmt.Sprintf("select * from email_delivering where keyword_id is null or keyword_id =0 or keyword_id = %d", keyword_id)
	rows, err := db.Query(q)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	var emails []string

	log.Println(rows.Columns())

	for rows.Next() {
		var id int
		var email string
		var keyword_id sql.NullInt64
		if err := rows.Scan(&id, &keyword_id, &email); err != nil {
			log.Fatal(err)
		}
		log.Println(email)
		emails = append(emails, email)
	}

	return emails, nil
}