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
	pageInfos		[]*PageInfo
}

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

func (p *PageScanner) scrollToEnd(divName string, infiniteScroll bool){
	var i int
	if infiniteScroll {
		i = -1000
	}else{
		i = 0
	}
	for i = 0 ; i<3;i++{
		numberOfElements := p.getPostArrayLength(divName)
		fmt.Println("Page scroll start")
		p.main_page.Mouse.Scroll(0, 99999, 5)
		time.Sleep(10 * time.Second)
		fmt.Println("Page scroll end")
		if numberOfElements == p.getPostArrayLength(divName){
			break
		}
	}
}

func (p *PageScanner) getPostArrayLength(divName string) int {
	elements := p.main_page.MustElements(fmt.Sprintf(`div.%s`,divName))
	return len(elements)
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
				Url : *href,
			}
			p.pageInfos = append(p.pageInfos, &pageInfo)
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

	waitConuter := 10

	for _, postURL := range p.pageInfos {

		waitConuter--
		if waitConuter <= 0 {
			time.Sleep(10 * time.Second)
			waitConuter = 10
		}

		p.main_page_wg.Add(1)
		go p.scanPage(postURL)
		// p.scanPage(postURL.url)
	}
	p.main_page_wg.Wait()
}

func (p *PageScanner) scanPage(pageInfo *PageInfo) {

	defer p.main_page_wg.Done()
	pageURL := pageInfo.Url

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
	p.getPostInfo(postPage, pageInfo)
	
	// fetch likes
	p.getLikes(postPage, pageInfo)

	// fetch hashtags(#) and profile_tags (@)
	p.getHashTagAndTag(postPage, pageInfo)

	// fmt.Println("--------------------------------------------------------")
	// fmt.Printf("Post: %s\nLikes: %s\nCaption: %s\n", pageURL, pageInfo.likes, pageInfo.caption)
	// fmt.Println("--------------------------------------------------------")
}

func (p *PageScanner) getLikes(postPage *rod.Page, pageInfo *PageInfo) {

	currentURL := postPage.MustInfo().URL
	likeSpan, err := postPage.Element(likesClassName)
	var likesText string
	

	if err == nil {
		likesText = likeSpan.MustText()
	} else {
		likesText = "[Error] : Unable to find likes on " + currentURL
	}
	pageInfo.Likes = likesText
}

func (p *PageScanner) getPostInfo(postPage *rod.Page, pageInfo *PageInfo) {

	currentURL := postPage.MustInfo().URL

	captionEl, err := postPage.Timeout(3 * time.Second).Element(postInfoClassName)

	var captionText string
	if err == nil {
		captionText = captionEl.MustText()
	} else {
		captionText = "[Error] : Caption not found on : "+ currentURL
	}
	pageInfo.Caption = captionText
}


func (p *PageScanner) getHashTagAndTag(postPage *rod.Page, pageInfo *PageInfo){
	captionElement, err := postPage.Timeout(3 * time.Second).Element(postInfoClassName)
	if err != nil {
		log.Fatal(fmt.Sprintf("PageScanner.getHashTagAndTags : %e", err))
		return
	}
	var hashtags, profileTags []string

	tags := captionElement.MustElements("a")
	for _, tag := range tags {
		tagText := tag.MustText()
		if tagText[0] == '#'{
			hashtags = append(hashtags, tagText[1:])
		}else if tagText[0] == '@'{
			profileTags = append(profileTags, tagText[1:])
		}else{
			log.Fatalf("PageScanner.getHashTagAndTags invalid tag text: %s", tagText )
		}
	}

	pageInfo.Hashtags = hashtags
	pageInfo.ProfileTags = profileTags
}