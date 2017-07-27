package twitter

import (
	"strconv"
	"strings"
)

// Twitter cards,
// https://dev.twitter.com/cards/markup

type TwitterCard struct {
	Card        string         `json:"card"`
	Site        *TwitterSite   `json:"site"`
	Creator     *TwitterSite   `json:"creator"`
	Description string         `json:"description"`
	Title       string         `json:"title"`
	Image       *TwitterImage  `json:"image"`
	Player      *TwitterPlayer `json:"player"`
}

type TwitterSite struct {
	User      string `json:"user"`
	TwitterID string `json:"twitter_id"`
}

type TwitterImage struct {
	Url string `json:"url"`
	Alt string `json:"alt"`
}

type TwitterPlayer struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Stream string `json:"stream"`
}

// TODO App

const (
	twitterPrefix = "twitter:"
	card          = "card"
	site          = "site"
	siteID        = "site:id"
	creator       = "creator"
	creatorID     = "creator:id"
	title         = "title"
	description   = "description"
	image         = "image:src"
	imageAlt      = "image:alt"
	player        = "player"
	playerWidth   = "player:width"
	playerHeight  = "player:height"
	playerStream  = "player:stream"
)

func GetTwitterCard(meta []map[string]string) *TwitterCard {
	var obj *TwitterCard
	for _, m := range meta {
		name, has := m["name"]
		if !has || !strings.HasPrefix(name, twitterPrefix) {
			continue
		}
		content, has := m["content"]
		if !has {
			continue
		}
		if obj == nil {
			obj = new(TwitterCard)
		}
		name = name[len(twitterPrefix):]
		switch name {
		case card:
			obj.Card = content
		case site:
			if obj.Site == nil {
				obj.Site = &TwitterSite{}
			}
			obj.Site.User = content
		case siteID:
			if obj.Site == nil {
				obj.Site = &TwitterSite{}
			}
			obj.Site.TwitterID = content
		case creator:
			if obj.Creator == nil {
				obj.Creator = &TwitterSite{}
			}
			obj.Creator.User = content
		case creatorID:
			if obj.Creator == nil {
				obj.Creator = &TwitterSite{}
			}
			obj.Creator.TwitterID = content
		case description:
			obj.Description = content
		case title:
			obj.Title = content
		case image:
			if obj.Image == nil {
				obj.Image = &TwitterImage{}
			}
			obj.Image.Url = content
		case imageAlt:
			if obj.Image == nil {
				obj.Image = &TwitterImage{}
			}
			obj.Image.Alt = content
		case player:
			if obj.Player == nil {
				obj.Player = &TwitterPlayer{}
			}
			obj.Player.Url = content
		case playerStream:
			if obj.Player == nil {
				obj.Player = &TwitterPlayer{}
			}
			obj.Player.Stream = content
		case playerHeight, playerWidth:
			size, err := strconv.Atoi(content)
			if err != nil {
				continue
			}
			if obj.Player == nil {
				obj.Player = &TwitterPlayer{}
			}
			if name == playerHeight {
				obj.Player.Height = size
			} else {
				obj.Player.Width = size
			}
		}
	}
	return obj
}
