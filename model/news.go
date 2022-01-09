package model

import "time"

type News struct {
	Id              int
	Body            string
	ImageUrl        string
	M3u8Url         string
	NewsId          string
	OutlineWithRuby string
	PublishedAtUtc  time.Time
	Title           string
	TitleWithRuby   string
	Url             string
}
