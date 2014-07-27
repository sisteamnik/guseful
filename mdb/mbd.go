package mdb

import (
	"fmt"
	"github.com/PuerkitoBio/purell"
	"github.com/coopernurse/gorp"
	"github.com/sisteamnik/guseful/website"
	"net/url"
	"sync"
	"time"
)

type Mdb struct {
	sites     []website.Site
	sitesLock sync.RWMutex

	pages     []website.SitePage
	pagesLock sync.RWMutex

	newpages     []website.SitePage
	newpagesLock sync.RWMutex
}

func (m *Mdb) store(db *gorp.DbMap) bool {
	tx, _ := db.Begin()
	defer tx.Commit()

	m.newpagesLock.Lock()
	if len(m.newpages) < 10 {
		return false
	}
	defer m.newpagesLock.Unlock()

	for _, v := range m.newpages {
		tx.Insert(&v)
	}
	m.newpages = []website.SitePage{}
	return true
}

func (m *Mdb) AddSite(s website.Site) bool {
	m.sitesLock.Lock()
	defer m.sitesLock.Unlock()
	m.sites = append(m.sites, s)
	return true
}

func (m *Mdb) AddPage(p website.SitePage) bool {
	_, ok := m.GetPage(p.Url)
	if ok {
		return false
	}
	_, ok = m.GetSite(0, p.Url)
	if !ok {
		return false
	}

	u := m.NormilizeUrl(p.Url)
	pu, _ := url.Parse(u)

	p.Url = pu.RequestURI()

	m.pagesLock.Lock()
	defer m.pagesLock.Unlock()
	m.pages = append(m.pages, p)
	m.newpagesLock.Lock()
	m.newpages = append(m.newpages, p)
	m.newpagesLock.Unlock()
	return true
}

func (m *Mdb) Start(s []website.Site, p []website.SitePage, db *gorp.DbMap) {
	m.sitesLock.Lock()
	m.pagesLock.Lock()
	defer m.sitesLock.Unlock()
	defer m.pagesLock.Unlock()
	m.sites = s
	m.pages = p
	go func() {
		for {
			time.Sleep(5 * time.Second)
			m.store(db)
		}
	}()
	fmt.Println("started with sites", len(m.sites))
	fmt.Println("started with pages", len(m.pages))
	fmt.Println("not visited", len(m.GetNotVisited()))

}

func (m *Mdb) GetPage(rawurl string) (website.SitePage, bool) {
	site, ok := m.GetSite(0, rawurl)
	siteId := int64(0)
	if ok {
		siteId = site.Id
	}

	u := m.NormilizeUrl(rawurl)
	pu, _ := url.Parse(u)

	m.pagesLock.RLock()
	defer m.pagesLock.RUnlock()
	for _, v := range m.pages {
		if v.SiteId == siteId && pu.RequestURI() == v.Url {
			v.Url = pu.String()
			return v, true
		}
	}
	return website.SitePage{}, false
}

func (m *Mdb) GetPageById(id int64) (website.SitePage, bool) {

	m.pagesLock.RLock()
	defer m.pagesLock.RUnlock()
	for _, v := range m.pages {
		if v.Id == id {
			v.Url = m.FullUrl(v.SiteId, v.Url)
			return v, true
		}
	}
	return website.SitePage{}, false
}

func (m *Mdb) FullUrl(siteId int64, urlpath string) string {
	u := m.NormilizeUrl(urlpath)
	pu, _ := url.Parse(u)

	if pu.Scheme != "" {
		return pu.String()
	}

	m.pagesLock.RLock()
	defer m.pagesLock.RUnlock()

	s, ok := m.GetSite(siteId, urlpath)
	if !ok {
		return ""
	}
	pu.Scheme = "http://"
	pu.Host = s.Domain
	return pu.String()
}

func (m *Mdb) GetNotVisited() (ss []string) {
	//m.pagesLock.RLock()
	//defer m.pagesLock.RUnlock()

	for _, v := range m.pages {
		if v.Visited == false && m.Alowed(m.FullUrl(v.SiteId, v.Url)) {
			ss = append(ss, m.FullUrl(v.SiteId, v.Url))
		}
	}
	return
}

func (m *Mdb) Alowed(rawurl string) bool {
	fmt.Println("ty")

	p, ok := m.GetSite(0, rawurl)
	if ok {
		return p.Allow
	}
	return false
}

func (m *Mdb) GetSite(id int64, rawurl string) (website.Site, bool) {
	m.sitesLock.RLock()
	defer m.sitesLock.RUnlock()

	u := m.NormilizeUrl(rawurl)
	pu, _ := url.Parse(u)

	if id == 0 && pu.Host == "" {
		return website.Site{}, false
	}

	for _, v := range m.sites {
		if v.Id == id || v.Domain == pu.Host {
			return v, true
		}
	}
	return website.Site{}, false
}

func (m *Mdb) SetVisited(rawurl string) bool {
	_, ok := m.GetPage(rawurl)

	u := m.NormilizeUrl(rawurl)
	pu, _ := url.Parse(u)

	if ok {
		m.pagesLock.Lock()
		for i, v := range m.pages {
			if v.Url == pu.RequestURI() {
				m.pages[i].Visited = true
			}
		}
		m.pagesLock.Unlock()
	}
	return false
}

func (m *Mdb) IsVisited(rawurl string) bool {
	return false
}

func (m *Mdb) ExistPage(rawurl string) bool {
	_, ok := m.GetPage(rawurl)
	if ok {
		return true
	}
	return false
}

func (m *Mdb) ExistSite(rawurl string) bool {
	m.sitesLock.RLock()
	defer m.sitesLock.RUnlock()
	u := m.NormilizeUrl(rawurl)
	pu, _ := url.Parse(u)
	for _, v := range m.sites {
		if v.Domain == pu.Host {
			return true
		}
	}
	return false
}

func (m *Mdb) NormilizeUrl(rawurl string) string {
	return purell.MustNormalizeURLString(rawurl, purell.FlagsUsuallySafeGreedy)
}
