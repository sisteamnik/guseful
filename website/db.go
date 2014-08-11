package website

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/purell"
	"github.com/coopernurse/gorp"
	"github.com/sisteamnik/guseful/gz"
	"net/url"
	"time"
)

func VisitPage(Db *gorp.DbMap, fullurl string, body []byte, code int) error {
	res := SitePage{}
	var err error

	fullurl = n(fullurl)
	u, _ := url.Parse(fullurl)

	if !IsExistSite(Db, fullurl) {
		_, err = AddSite(Db, fullurl)
		if err != nil {
			return err
		}
	}

	if !IsExistPage(Db, fullurl) {
		res.Url = u.RequestURI()
		res.Error = int64(code)
		res.Body = gz.Gz(string(body))
		res.Visited = time.Now().UnixNano()

		site, err := GetSite(Db, fullurl)
		if err != nil {
			return err
		}
		res.SiteId = site.Id

		err = Db.Insert(&res)
		if err != nil {
			return err
		}

		return nil

	} else {
		res, err = GetPage(Db, fullurl)
		if err != nil {
			return err
		}
	}

	res.Error = int64(code)
	res.Body = gz.Gz(string(body))
	res.Visited = time.Now().UnixNano()

	_, err = Db.Update(&res)

	return err
}

func AddSite(Db *gorp.DbMap, fullurl string) (Site, error) {
	fullurl = n(fullurl)
	u, _ := url.Parse(fullurl)

	if u.Host == "" || u.Host == ":" {
		return Site{}, errors.New("empty host")
	}

	s, _ := GetSite(Db, fullurl)
	if s.Id != 0 {
		return s, errors.New("exist")
	} else {
		s.Domain = u.Host
		err := Db.Insert(&s)
		go AddPage(Db, fullurl)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

func AddPage(Db *gorp.DbMap, fullurl string) (SitePage, error) {
	fullurl = n(fullurl)
	u, _ := url.Parse(fullurl)
	res := SitePage{}
	if IsExistPage(Db, fullurl) {
		return res, errors.New("exist")
	}

	if !IsExistSite(Db, fullurl) {
		_, err := AddSite(Db, fullurl)
		if err != nil {
			return res, err
		}
	}

	site, err := GetSite(Db, fullurl)
	if err != nil {
		return res, errors.New("site not exist")
	}

	res.SiteId = site.Id
	res.Url = u.RequestURI()

	err = Db.Insert(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func IsExistPage(Db *gorp.DbMap, fullurl string) bool {
	fullurl = n(fullurl)
	exist := false

	p, err := GetPage(Db, fullurl)
	if err == nil && p.Id != 0 {
		exist = true
	}

	return exist
}

func IsExistSite(Db *gorp.DbMap, fullurl string) bool {
	fullurl = n(fullurl)
	exist := false
	p, err := GetSite(Db, fullurl)
	if err == nil && p.Id != 0 {
		exist = true
	}

	return exist
}

func IsVisitedPage(Db *gorp.DbMap, fullurl string) bool {
	fullurl = n(fullurl)
	u, _ := url.Parse(fullurl)
	visited := false

	vis, err := Db.SelectInt("select Visited from SitePage"+
		" where SiteId in (select Id from Site where Domain = ?"+
		") and Url = ?", u.Host, u.RequestURI())
	if err != nil {
		return false
	}
	if vis > 0 {
		visited = true
	}

	return visited
}

func GetPage(Db *gorp.DbMap, fullurl string) (SitePage, error) {
	fullurl = n(fullurl)
	w := SitePage{}
	u, _ := url.Parse(fullurl)

	err := Db.SelectOne(&w, "select * from SitePage where SiteId"+
		" in (select Id from Site where Domain = ?) and Url = ?",
		u.Host, u.RequestURI())
	if err != nil || w.Url == "" {
		fmt.Println(err)
		return w, errors.New("not exist")
	}
	w.Url = fullurl

	return w, nil
}

func GetSite(Db *gorp.DbMap, fullurl string) (Site, error) {
	fullurl = n(fullurl)
	w := Site{}
	u, _ := url.Parse(fullurl)
	err := Db.SelectOne(&w, "select * from Site where Domain="+
		" ?", u.Host)
	if err != nil {
		fmt.Println(err)
		return w, err
	}
	if w.Id == 0 {
		fmt.Println("Site not exist")
		return w, errors.New("not exist")
	}
	return w, nil
}

func Normalize(fullurl string) string {
	fullurl = purell.MustNormalizeURLString(fullurl, purell.FlagsSafe|
		purell.FlagSortQuery|purell.FlagRemoveFragment|
		purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveDotSegments|
		purell.FlagRemoveWWW)
	u, _ := url.Parse(fullurl)
	if u.Path == "" {
		u.Path = "/"
	}
	if u.Scheme == "" {
		u.Scheme = "http://"
	}
	q := u.Query()
	q.Del("track")
	u.RawQuery = q.Encode()
	fullurl = u.String()
	return fullurl
}

func n(fullurl string) string {
	return Normalize(fullurl)
}
