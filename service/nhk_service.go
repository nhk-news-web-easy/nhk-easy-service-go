package service

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/db"
	"github.com/nhk-news-web-easy/nhk-easy-service-go/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

var (
	thirtyDaysInHours = float64(24 * 30)
)

func GetNews(startDate *timestamppb.Timestamp, endDate *timestamppb.Timestamp) ([]model.News, error) {
	if startDate == nil || endDate == nil {
		return nil, errors.New("startDate and endDate are required")
	}

	startTime, endTime := startDate.AsTime(), endDate.AsTime()

	if startTime.After(endTime) {
		return nil, errors.New("startDate should be less than endDate")
	}

	if endTime.Sub(startTime).Hours() > thirtyDaysInHours {
		return nil, errors.New("time range is too large")
	}

	rows, err := db.Query("select id, body, image_url, m3u8url, news_id, outline_with_ruby, "+
		"published_at_utc, title, title_with_ruby, url from news where published_at_utc >= ? and published_at_utc <= ?", startTime, endTime)

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

		newsList = append(newsList, news)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}
