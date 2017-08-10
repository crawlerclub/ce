package ce

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
)

func MD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func MakeAbsoluteUrl(href, baseurl string) (string, error) {
	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	base, err := url.Parse(baseurl)
	if err != nil {
		return "", err
	}
	u = base.ResolveReference(u)
	return u.String(), nil
}
