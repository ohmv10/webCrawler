package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-rod/rod/lib/input"

	"github.com/go-rod/rod"
)

type PageScanner struct {
	pageURLs 		[]string
	browser  		*rod.Browser
	main_page_wg    *sync.WaitGroup
	main_page 		*rod.Page
	pageInfos		[]PageInfo
}

// func getPageScanner(pageURLs []string) PageScanner {
// 	var pg_wg sync.WaitGroup
// 	return PageScanner{
// 		pageURLs: 		pageURLs,
// 		main_page_wg:   &pg_wg,
// 	}
// }

func createPageScanner(url string) PageScanner {
	var page_wg sync.WaitGroup
	browser := rod.New().MustConnect()
	page := browser.MustPage(url)
	page.MustWaitLoad()
	time.Sleep(3*time.Second)

	return PageScanner{
		main_page_wg: &page_wg, 
		browser: browser,
		main_page: page,
	}
}

func (p *PageScanner) closeMainPage(){
	pageCloser(p.main_page)
}

func (p *PageScanner) loginInstagram(username, password string){
	// login instagram
	p.main_page.MustElement("input[name='username']").MustInput(username)
	p.main_page.MustElement("input[name='password']").MustInput(password).MustType(input.Enter)
	p.main_page.MustWaitNavigation()
	
	// wait to avoid suspissions
	time.Sleep(5 * time.Second)
}

func (p *PageScanner) navigateToProfileWithURL(url string){
	// navigate
	p.main_page.MustNavigate(url)
	p.main_page.MustWaitLoad()
	
	// wait to avoid suspission
	fmt.Println("Page loading start")
	time.Sleep(5 * time.Second)
	fmt.Println("Page loading end")
}


func (p *PageScanner) scrollToEnd(){

	fmt.Println("Page scroll start")
	p.main_page.Mouse.Scroll(0, 99999, 5)
	time.Sleep(10 * time.Second)
	fmt.Println("Page scroll end")

}

func (p *PageScanner) updatePostSlice(divName string){
	
	// get div
	fmt.Println("Element finder start")
	divs := p.main_page.MustElements(fmt.Sprintf(`div.%s`, divName))
	// div, err := p.main_page.Element(fmt.Sprintf(`div.%s`, divName))
	fmt.Println("Element finder end")

	for _, div := range divs{
		links := div.MustElements("a")
		
		for _, link := range links {
			// Get the href attribute
			href, err := link.Attribute("href")
			
			if err != nil || href == nil {
				log.Fatal("no link found")
				continue
			}
			// Append to the pageURLs slice
			pageInfo := PageInfo{
				url : *href,
			}
			p.pageInfos = append(p.pageInfos, pageInfo)
		}
	}
	// log all the urls
	// p.printURLS()
}

func (p *PageScanner) printURLS(){
	fmt.Println("URL printer : ")
	for _,pageInfo := range p.pageInfos{
		fmt.Println(pageInfo)
	}
}


func (p *PageScanner) scanPages(main_wg *sync.WaitGroup) {
	defer main_wg.Done()
	fmt.Println("Page scanner started")
	for _, postURL := range p.pageInfos {
		p.main_page_wg.Add(1)
		go p.scanPage(postURL.url)
		// p.scanPage(postURL.url)
	}
	p.main_page_wg.Wait()
}

func (p *PageScanner) scanPage(pageURL string) {

	defer p.main_page_wg.Done()

	// fmt.Print("Page Started : ", pageURL," : ")
	// fmt.Println(pageURL == "https://www.instagram.com/districtupdates/p/DIjYNt1TVYI/")
	
	// Load page and wait 2 sec for suspission
	postPage := p.browser.MustPage(fmt.Sprintf("https://www.instagram.com%s",pageURL)).MustWaitLoad()
	// fmt.Println(postPage.MustInfo().URL)
	// postPage := p.browser.MustPage(pageURL)
	time.Sleep(10 * time.Second)

	// close page at end
	defer postPage.Close()

	// fetch captions
	postInfo := p.getPostInfo(postPage)

	// fetch likes
	likes := p.getLikes(postPage)
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("Post: %s\nLikes: %s\nCaption: %s\n", pageURL, likes, postInfo)
	fmt.Println("--------------------------------------------------------")
}

func (p *PageScanner) getLikes(postPage *rod.Page) string {

	currentURL := postPage.MustInfo().URL
	likeSpan, err := postPage.Element("span.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x1hl2dhg.x16tdsg8.x1vvkbs")
	var likesText string
	

	if err == nil {
		likesText = likeSpan.MustText()
	} else {
		likesText = "[Error] : Unable to find likes on " + currentURL
	}


	return likesText
}

func (p *PageScanner) getPostInfo(postPage *rod.Page) string {


	currentURL := postPage.MustInfo().URL

	captionEl, err := postPage.Timeout(3 * time.Second).Element("h1._ap3a._aaco._aacu._aacx._aad7._aade")

	var captionText string
	if err == nil {
		captionText = captionEl.MustText()
	} else {
		captionText = "[Error] : Caption not found on : "+ currentURL
	}
	return captionText
}


