package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/db"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/model"
	"log"
)

func GetNews() ([]model.News, error) {
	rows, err := db.Query("select id, body, image_url, m3u8url, news_id, outline_with_ruby, published_at_utc, title, title_with_ruby, url from news")

	if err != nil {
		log.Printf("failed to query db %v", err)

		return nil, err
	}

	defer rows.Close()

	var newsList []model.News

	for rows.Next() {
		var news model.News

		err := rows.Scan(&news.Id, &news.Body, &news.ImageUrl, &news.M3u8Url, &news.NewsId, &news.OutlineWithRuby, &news.PublishedAtUtc, &news.Title, &news.TitleWithRuby, &news.Url)

		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}
