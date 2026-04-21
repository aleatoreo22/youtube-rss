package youtuberss

import (
	"bufio"
	"encoding/xml"
	"net/http"
	"strings"
)

type YoutubeClient struct {
}

func New() *YoutubeClient {
	return &YoutubeClient{}
}

func (yc *YoutubeClient) GetFeed(url string) (*Feed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var xmlstring strings.Builder
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		xmlstring.WriteString(scanner.Text())
	}

	rss := &Feed{}
	err = xml.Unmarshal([]byte(xmlstring.String()), rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}
