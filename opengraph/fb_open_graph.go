package og

import (
	"strconv"
	"strings"
)

// Open Graph protocol, http://ogp.me
// by facebook
// https://developers.facebook.com/docs/sharing/webmasters

type Ogp struct {
	// the four required properties
	Title string      `json:"title"` // og:title
	Type  string      `json:"type"`  // og:type
	Image []*OgpImage `json:"image"` // og:image
	Url   string      `json:"url"`   // og:url

	// optional metadata
	Audio           []*OgpAudio `json:"audio"`          // og:audio
	Description     string      `json:"description"`    // og:description
	Determiner      string      `json:"determiner"`     // og:determiner
	Locale          string      `json:"locale"`         // og:locale
	LocaleAlternate []string    `json:locale_alternate` // og:locale:alternate
	SiteName        string      `json:"site_name"`      // og:site_name
	Video           []*OgpVideo `json:"video"`          // og:video
}

type OgpImage struct {
	Url       string `json:"url"`
	SecureUrl string `json:"secure_url"`
	Type      string `json:"type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type OgpVideo struct {
	Url       string `json:"url"`
	SecureUrl string `json:"secure_url"`
	Type      string `json:"type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type OgpAudio struct {
	Url       string `json:"url"`
	SecureUrl string `json:"secure_url"`
	Type      string `json:"type"`
}

const (
	ogPrefix    = "og:"
	imagePrefix = "image:"
	audioPrefix = "audio:"
	videoPrefix = "video:"
	title       = "title"
	typ         = "type"
	image       = "image"
	url         = "url"
	audio       = "audio"
	description = "description"
	determiner  = "determiner"
	locale      = "locale"
	altLocale   = "locale:alternate"
	siteName    = "site_name"
	video       = "video"
	secURL      = "secure_url"
	width       = "width"
	height      = "height"
)

func GetOgp(meta []map[string]string) *Ogp {
	var obj *Ogp
	var (
		img *OgpImage
		vid *OgpVideo
		aud *OgpAudio
	)
	for _, m := range meta {
		p, has := m["property"]
		if !has || !strings.HasPrefix(p, ogPrefix) {
			continue
		}
		content, has := m["content"]
		if !has {
			continue
		}
		if obj == nil {
			obj = new(Ogp)
		}
		name := p[len(ogPrefix):]
		if strings.HasPrefix(name, imagePrefix) {
			if img == nil {
				continue
			}
			name = name[len(imagePrefix):]
			switch name {
			case secURL:
				img.SecureUrl = content
			case typ:
				img.Type = content
			case height, width:
				size, err := strconv.Atoi(content)
				if err != nil {
					continue
				}
				if name == height {
					img.Height = size
				} else {
					img.Width = size
				}
			}
			continue
		}

		if strings.HasPrefix(name, videoPrefix) {
			if vid == nil {
				continue
			}
			name = name[len(videoPrefix):]
			switch name {
			case secURL:
				vid.SecureUrl = content
			case typ:
				vid.Type = content
			case height, width:
				size, err := strconv.Atoi(content)
				if err != nil {
					continue
				}
				if name == height {
					vid.Height = size
				} else {
					vid.Width = size
				}
			}
			continue
		}

		if strings.HasPrefix(name, audioPrefix) {
			if aud == nil {
				continue
			}
			switch name[len(audioPrefix):] {
			case secURL:
				aud.SecureUrl = content
			case typ:
				aud.Type = content
			}
			continue
		}

		switch name {
		case title:
			obj.Title = content
		case typ:
			obj.Type = obj.Type
		case image:
			if img != nil {
				obj.Image = append(obj.Image, img)
			}
			img = &OgpImage{}
			img.Url = content
		case url:
			obj.Url = content
		case audio:
			if aud != nil {
				obj.Audio = append(obj.Audio, aud)
			}
			aud = &OgpAudio{}
			aud.Url = content
		case description:
			obj.Description = content
		case determiner:
			obj.Determiner = content
		case locale:
			obj.Locale = content
		case altLocale:
			obj.LocaleAlternate = append(obj.LocaleAlternate, content)
		case siteName:
			obj.SiteName = content
		case video:
			if vid != nil {
				obj.Video = append(obj.Video, vid)
			}
			vid = &OgpVideo{}
			vid.Url = content
		}
	}
	if img != nil {
		obj.Image = append(obj.Image, img)
	}
	if vid != nil {
		obj.Video = append(obj.Video, vid)
	}
	if aud != nil {
		obj.Audio = append(obj.Audio, aud)
	}
	return obj
}
