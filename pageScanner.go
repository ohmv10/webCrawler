package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type PageScanner struct {
	pageURLs []string
	browser  *rod.Browser
	wg       *sync.WaitGroup
}

func getPageScanner(pageURLs []string, browser *rod.Browser) PageScanner {
	var pg_wg sync.WaitGroup
	return PageScanner{
		pageURLs: pageURLs,
		browser:  browser,
		wg:       &pg_wg,
	}
}

func (p *PageScanner) getLikes(postPage *rod.Page) string {

	currentURL := postPage.MustInfo().URL
	if(currentURL == "https://www.instagram.com/districtupdates/p/DIjYNt1TVYI/"){
		fmt.Println("Problemmatic Page")
	}
	likeSpan, err := postPage.Element("span.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x1hl2dhg.x16tdsg8.x1vvkbs")
	var likesText string
	

	if err == nil {
		likesText = likeSpan.MustText()
	} else {
		likesText = "[Error] : Unable to find likes on " + currentURL
	}

	if(currentURL == "https://www.instagram.com/districtupdates/p/DIjYNt1TVYI/"){
		fmt.Println("error : ", err)
		fmt.Println("Likes on error page : ", likesText)
	}

	return likesText
}

func (p *PageScanner) getPostInfo(postPage *rod.Page) string {

	currentURL := postPage.MustInfo().URL

	if(currentURL == "https://www.instagram.com/districtupdates/p/DIjYNt1TVYI/"){
		fmt.Println("Problemmatic Page CAPTION")
	}

	captionEl, err := postPage.Timeout(3 * time.Second).Element("h1._ap3a._aaco._aacu._aacx._aad7._aade")

	var captionText string
	if err == nil {
		captionText = captionEl.MustText()
	} else {
		captionText = "[Error] : Caption not found"
	}
	if(currentURL == "https://www.instagram.com/districtupdates/p/DIjYNt1TVYI/"){
		fmt.Println("error caption : ", err)
		fmt.Println("Likes on error page : ", captionText)
	}
	return captionText
}

func (p *PageScanner) scanPage(pageURL string) {

	defer p.wg.Done()
	// Load page and wait 2 sec for suspission
	postPage := p.browser.MustPage(pageURL).MustWaitLoad()
	// postPage := p.browser.MustPage(pageURL)
	fmt.Print("Page Started : ", pageURL)
	fmt.Println(pageURL == "https://www.instagram.com/districtupdates/p/DIjYNt1TVYI/")
	fmt.Println(postPage.MustInfo().URL)
	time.Sleep(10 * time.Second)

	// close page at end
	// defer postPage.Close()

	// fetch captions
	postInfo := p.getPostInfo(postPage)

	// fetch likes
	likes := p.getLikes(postPage)
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("Post: %s\nLikes: %s\nCaption: %s\n", pageURL, likes, postInfo)
	fmt.Println("--------------------------------------------------------")
}

func (p *PageScanner) scanPages(main_wg *sync.WaitGroup) {
	defer main_wg.Done()
	for _, postURL := range p.pageURLs {
		p.wg.Add(1)
		go p.scanPage(postURL)
	}
	p.wg.Wait()
}


