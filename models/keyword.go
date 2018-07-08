package models

import (
	"fmt"
	"log"
	"encoding/json"
)

type Keyword struct {
	Id, LastCrawledAt, CrawlDelayTime int
	TimeOfArticle, Keyword            string
}

func AllKeyword() ([]*Keyword, error) {
	rows, err := db.Query("SELECT * FROM keywords")
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(kw.Keyword)
		out, _ := json.Marshal(kw)
		fmt.Printf(string(out))
		words = append(words, &kw)
	}

	return words, nil
}