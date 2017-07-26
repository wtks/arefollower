package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type (
	Ranking struct {
		PubDate time.Time
		Videos  []*Video
	}
	Video struct {
		Id           string
		Title        string
		UploadDate   time.Time
		ThumbUrl     string
		Length       string
		View         int
		Comment      int
		Mylist       int
		Url          string
		Tags         []string
	}
)

const (
	RANKING_URL    = "http://www.nicovideo.jp/ranking/fav/daily/are?rss=2.0&lang=ja-jp"
	VIDEO_INFO_URL = "http://ext.nicovideo.jp/api/getthumbinfo/"
	VIDEO_URL_BASE = "http://www.nicovideo.jp/watch/"
	TIME_FORMAT    = "Mon, 2 Jan 2006 15:04:05 -0700"
)

func FetchRanking() (*Ranking, error) {
	type Rss struct {
		Channel struct {
			PubDate string `xml:"pubDate"`
			Items   []struct {
				Link string `xml:"link"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	data, err := httpGet(RANKING_URL)
	if err != nil {
		return nil, err
	}

	rss := new(Rss)
	if err := xml.Unmarshal([]byte(data), &rss); err != nil {
		return nil, err
	}

	date, _ := time.Parse(TIME_FORMAT, rss.Channel.PubDate)
	result := &Ranking{
		PubDate: date,
		Videos:  make([]*Video, len(rss.Channel.Items)),
	}
	for i, item := range rss.Channel.Items {
		video, _ := FetchVideoInfo(strings.TrimPrefix(item.Link, VIDEO_URL_BASE))
		result.Videos[i] = video
	}

	return result, nil
}

func FetchVideoInfo(id string) (*Video, error) {
	type Res struct {
		Status string `xml:"status,attr"`
		Thumb  struct {
			VideoId       string   `xml:"video_id"`
			Title         string   `xml:"title"`
			ThumbnailUrl  string   `xml:"thumbnail_url"`
			FirstRetrieve string   `xml:"first_retrieve"`
			Length        string   `xml:"length"`
			ViewCounter   int      `xml:"view_counter"`
			CommentNum    int      `xml:"comment_num"`
			MylistCounter int      `xml:"mylist_counter"`
			WatchUrl      string   `xml:"watch_url"`
			Tags          []string `xml:"tags"`
		} `xml:"thumb"`
	}

	data, err := httpGet(VIDEO_INFO_URL + id)
	if err != nil {
		return nil, err
	}
	res := new(Res)
	if err := xml.Unmarshal([]byte(data), &res); err != nil || res.Status != "ok" {
		return nil, err
	}

	date, _ := time.Parse(TIME_FORMAT, res.Thumb.FirstRetrieve)

	return &Video{
		Id:           res.Thumb.VideoId,
		Title:        res.Thumb.Title,
		UploadDate:   date,
		ThumbUrl:     res.Thumb.ThumbnailUrl,
		Length:       res.Thumb.Length,
		View:         res.Thumb.ViewCounter,
		Comment:      res.Thumb.CommentNum,
		Mylist:       res.Thumb.MylistCounter,
		Url:          res.Thumb.WatchUrl,
		Tags:         res.Thumb.Tags,
	}, nil
}

func httpGet(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
