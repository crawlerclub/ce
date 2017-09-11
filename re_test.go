package ce

import (
	"regexp"
	"testing"
)

func TestReImg(t *testing.T) {
	cases := [][]string{
		{`    <img id="i-cd494f35d267741a" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7" data-src="http://i.dailymail.co.uk/i/pix/2017/08/20/23/4365CF9100000578-0-image-m-32_1503267632977.jpg" height="504" width="634" alt="Sickening: An execution scene from Channel 4 drama The State" class="blkBorder img-share"/>`, "http://i.dailymail.co.uk/i/pix/2017/08/20/23/4365CF9100000578-0-image-m-32_1503267632977.jpg"},
		{`<img class="lazy opacity_0 " id="img_0" data-original="http://zkres2.myzaker.com/201709/59b5d288a07aecc83802e3c5_640.jpg" data-height='318' data-width='600' />`,
			`http://zkres2.myzaker.com/201709/59b5d288a07aecc83802e3c5_640.jpg`},
	}
	//ReImgSrc = regexp.MustCompile(`(?ims)<img.+?src=\s*?"(.+?)"|'(.+?)'.*?>`)
	//ReImgSrc = regexp.MustCompile(`(?ims).+?src=\s*?"(.+?)"|'(.+?)'`)
	ReImgSrc = regexp.MustCompile(`(?ims)(?:.+?src|data-original)=\s*?"(.+?)"|'(.+?)'`)
	for _, c := range cases {
		ret := ReImg.FindAllStringSubmatch(c[0], -1)
		t.Log(ret)
		for _, r := range ret {
			src := ReImgSrc.FindAllStringSubmatch(r[0], -1)
			t.Log(src)
			for _, s := range src {
				t.Log(s[1])
			}
		}
		t.Log("\n\n")
	}
}
