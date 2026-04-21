package helper

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"youtube-rss.api.aleatoreo.com/internal/domain"
	"youtube-rss.api.aleatoreo.com/package/youtuberss"
)

const CHANNEL_INDEX = "channel/"

func getYoutubeChannelUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Failed to get channel url")
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyString := string(bodyBytes)
	url = bodyString[strings.Index(bodyString, "https://www.youtube.com/channel/"):]
	url = url[:strings.Index(url, "\">")]

	return url, nil
}

func getYouTubeRSS(url string) (string, error) {
	if !strings.Contains(url, CHANNEL_INDEX) {
		var err error
		url, err = getYoutubeChannelUrl(url)
		if err != nil {
			return "", err
		}
	}
	channelId := url[strings.Index(url, CHANNEL_INDEX)+len(CHANNEL_INDEX):]

	return "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelId, nil
}

func ValidateChannel(channel domain.Channel) (domain.Channel, error) {
	if channel.RssUrl == "" && channel.Url == "" {
		return channel, errors.New("Channel must have a rss url or a url")
	}

	if channel.Url != "" && channel.RssUrl == "" {
		rssUrl, err := getYouTubeRSS(channel.Url)
		if err != nil {
			return channel, err
		}
		channel.RssUrl = rssUrl
	}

	yrss := youtuberss.New()
	feed, err := yrss.GetFeed(channel.RssUrl)
	if err != nil {
		return channel, err
	}
	if feed.ID == "" {
		return channel, errors.New("Cannot communicate with YouTube")
	}
	channel.Name = feed.Title

	return channel, nil
}
