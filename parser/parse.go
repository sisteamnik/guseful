package parser

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/sisteamnik/guseful/website"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
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

	for _, v := range pr {
		links := getLinks(getString(v.Index), v.Index)
		for _, val := range links {
			if website.IsVisitedPage(Db, val) {
				continue
			}
			content := getString(val)
			time.Sleep(2000 * time.Millisecond)
			code := 404
			if content != "" {
				code = 200
			}
			err := website.VisitPage(Db, val, []byte(content), code)
			if err != nil {
				continue
			}

			var rules = map[string]string{}
			err = json.Unmarshal([]byte(v.Rules), &rules)
			if err != nil {
				continue
			}
			doc, err := goquery.NewDocumentFromReader(strings.
				NewReader(content))
			if err != nil {
				continue
			}
			doc.Find("script,style").Each(func(i int, s *goquery.Selection) {
				removeNode(s)
			})
			rres := map[string][]string{}
			for key, value := range rules {
				s := doc.Find(value)
				if len(s.Nodes) > 1 {
					continue
				}
				con := strings.TrimSpace(s.Text())
				result := con

				ar := strings.Split(result, "\n")

				r, _ := regexp.Compile("[\n]{2,}")

				for i, v := range ar {
					stemp := strings.TrimSpace(v)
					if len(stemp) == 0 {
						ar[i] = ""
						continue
					}
				}

				if key == "content" {
					if len(result) < 500 {
						continue
					}
				}

				result = strings.Join(ar, "\n")
				result = r.ReplaceAllString(result, "\n\n")
				rres[key] = []string{result}

				rres["source"] = []string{val}
				if len(rres["title"]) == 1 && rres["title"][0] != "" &&
					len(rres["content"]) == 1 && len(rres["content"][0]) > 300 {
					res[val] = rres
				}

				if len(res) > 10 {
					return
				}
			}

		}
	}
	return
}

func (p *Prsr) GetAll(pagetype string) ([]Parser, error) {
	pr := []Parser{}
	_, err := p.Opts.Db.Select(&pr, "select * from Parser where LastVisit < ? and Type"+
		" = ?",
		time.Now().UnixNano()-p.Opts.Delay, pagetype)
	return pr, err
}

func (p *Prsr) CreateParser(prs *Parser) error {
	err := p.Opts.Db.Insert(prs)
	return err
}

func removeNode(selection *goquery.Selection) {
	if selection != nil {
		node := selection.Get(0)
		if node != nil && node.Parent != nil {
			node.Parent.RemoveChild(node)
		}
	}
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
	return res
}
