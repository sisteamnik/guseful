package parser

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	Db := p.Opts.Db
	pr := []Parser{}
	res = map[string]map[string][]string{}
	Db.Select(&pr, "select * from Parser where LastVisit < ? and Type"+
		" = ?",
		time.Now().UnixNano()-p.Opts.Delay, pagetype)

	fmt.Println(pr)

	for _, v := range pr {
		links := getLinks(getString(v.Index), v.Index)
		fmt.Println(links)
		for _, val := range links {
			fmt.Println("Start visit ", val)
			if website.IsVisitedPage(Db, val) {
				fmt.Println(val, " - visited")
				continue
			}
			content := getString(val)
			time.Sleep(500 * time.Millisecond)
			code := 404
			if content != "" {
				code = 200
			}
			err := website.VisitPage(Db, val, []byte(content), code)
			if err != nil {
				fmt.Println("err", err)
				continue
			}

			var rules = map[string]string{}
			err = json.Unmarshal([]byte(v.Rules), &rules)
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			doc, err := goquery.NewDocumentFromReader(strings.
				NewReader(content))
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			rres := map[string][]string{}
			for key, value := range rules {
				s := doc.Find(value)
				rres[key] = []string{s.Text()}
			}
			rres["source"] = []string{val}
			fmt.Println(rres)
			res[val] = rres
		}
	}
	fmt.Println("All visited")
	return
}

func getString(u string) string {
	response, err := http.Get(u)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
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

	mu, err := url.Parse(contextUrl)
	if err != nil {
		return res
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")

		u, err := mu.Parse(val)
		if err != nil {
			return
		}
		if u.Host != mu.Host && u.Host != "www."+mu.Host &&
			mu.Host != "www."+u.Host {
			return
		}
		val = website.Normalize(u.String())
		for _, v := range res {
			if v == val {
				return
			}
		}
		res = append(res, val)
	})
	fmt.Println(res)
	return res
}
