package website

import (
	"github.com/PuerkitoBio/purell"
	"github.com/coopernurse/gorp"
)

func VisitPage(Db *gorp.DbMap, fullurl string, body []byte, code int) error {
	res := website.SitePage{}
	var err error

	//fmt.Printf("New pages %d, to up %d\n", len(ToIns), len(ToUp))

	fullurl = n(fullurl)
	u, _ := url.Parse(fullurl)

	if !IsExistSite(fullurl) {
		_, err = AddSite(Db, fullurl)
		if err != nil {
			return err
		}
	}

	if !IsExistPage(fullurl) {
		res.Url = u.RequestURI()
		res.Error = int64(code)
		res.Body = gz.Gz(string(body))
		res.Visited = time.Now().UnixNano()

		Db.Insert(&res)

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

	Db.Update(&res)

	return nil
}

func AddSite(fullurl string) (website.Site, error) {
	fullurl = n(fullurl)
	u, _ := url.Parse(fullurl)

	if u.Host == "" || u.Host == ":" {
		return website.Site{}, errors.New("empty host")
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

func AddPage(Db, fullurl string) (website.SitePage, error) {
	fullurl = n(fullurl)
	u, _ := url.Parse(fullurl)
	res := website.SitePage{}
	if IsExistPage(fullurl) {
		return res, errors.New("exist")
	}

	if !IsExistSite(fullurl) {
		_, err := AddSite(fullurl)
		if err != nil {
			return res, err
		}
	}

	site, err := GetSite(fullurl)
	if err != nil {
		return res, errors.New("site not exist")
	}

	res.SiteId = site.Id
	res.Url = u.RequestURI()

	ToInsLock.Lock()
	ToIns = append(ToIns, res)
	ToInsLock.Unlock()

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
	u, _ := url.Parse(fullurl)
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

	vis, _ := Db.SelectInt("select Visited from SitePage"+
		" where SiteId in (select Id from Site where Domain = ?"+
		" limit 1) and Url = ?", u.Host, u.RequestURI())
	if vis > 0 {
		visited = true
	}

	return visited
}

func GetPage(Db *gorp.DbMap, fullurl string) (website.SitePage, error) {
	fullurl = n(fullurl)
	w := website.SitePage{}
	u, _ := url.Parse(fullurl)

	err := Db.SelectOne(&w, "select * from SitePage where SiteId"+
		" in (select Id from Site where Domain = ? limit 1) and Url = ?",
		u.Host, u.RequestURI())
	if err != nil || w.Url == "" {
		return w, errors.New("not exist")
	}
	w.Url = fullurl

	return w, nil
}

func GetSite(Db *gorp.DbMap, fullurl string) (website.Site, error) {
	fullurl = n(fullurl)
	w := website.Site{}
	u, _ := url.Parse(fullurl)
	err := Db.SelectOne(&w, "select * from Site where Domain="+
		" ?", u.Host)
	if err != nil {
		return w, err
	}
	if w.Id == 0 {
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
	fullurl = u.String()
	return fullurl
}

func n(fullurl string) string {
	return Normalize(fullurl)
}
