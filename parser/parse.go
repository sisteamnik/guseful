package parser

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/coopernurse/gorp"
	"github.com/kennygrant/sanitize"
	"github.com/sisteamnik/guseful/gz"
	"github.com/sisteamnik/guseful/website"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewParser(opts Opts) *Prsr {
	p := new(Prsr)
	p.Opts = opts

	return p
}

func (p *Prsr) ParseAll(pagetype string) (res map[string]map[string][]string) {
	pr := []Parser{}
	res = map[string]map[string][]string{}
	p.Opts.db.Select(&pr, "select * from Parser where LastVisit < ? and Type"+
		" = ?",
		time.Now().UnixNano()-p.Opts.Delay, pagetype)

	for _, v := range pr {
		links := getLinks(getString(v.Index))
		for _, val := range links {
			content := getString(val)
			code := 404
			if content != "" {
				code = 200
			}
			err := website.VisitPage(Db, val, gz.Gz(content), code)
			if err != nil {
				continue
			}

			var rules = map[string]string{}
			err = json.Unmarshal(v.Rules, &rules)
			if err != nil {
				continue
			}
			doc, err := goquery.NewDocumentFromReader(strings.
				NewReader(content))
			if err != nil {
				continue
			}
			rres := map[string][]string{}
			for key, value := range rules {
				s := doc.Find(value)
				rres[key] = []string{s.Text()}
			}
			res[val] = rres
		}
	}
	return
}

func getString(u string) string {
	response, err := http.Get(u)
	if err != nil {
		return ""
	} else {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return ""
		}
		response.Body.Close()
		return string(contents)
	}
	return ""
}

func getLinks(p string, contextUrl string) []string {
	res := []string{}
	r := strings.NewReader(p)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return res
	}

	u, err := url.Parse(contextUrl)
	if err != nil {
		return res
	}
	u.Fragment = ""

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")

		u, err := u.Parse(val)
		if err != nil {
			return
		}
		u.Fragment = ""
		for _, v := range res {
			if v == u.String() {
				return
			}
		}
		res = append(res, u.String())
	})
	return res
}
