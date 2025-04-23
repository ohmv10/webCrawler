package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"os"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func pageCloser(page *rod.Page){
	fmt.Println("page close triggered")
	page.Close()

}


func main() {
	browser := rod.New().MustConnect()
	page := browser.MustPage("https://www.instagram.com/accounts/login/")
	page.MustWaitLoad()
	defer pageCloser(page)
	time.Sleep(3 * time.Second)


	// Login
	page.MustElement("input[name='username']").MustInput(os.Getenv("INSTA_USERNAME"))
	page.MustElement("input[name='password']").MustInput(os.Getenv("INSTA_PASSWORD")).MustType(input.Enter)
	page.MustWaitNavigation()
	time.Sleep(5 * time.Second)

	// Navigate to profile
	page.MustNavigate("https://www.instagram.com/districtupdates/")
	page.MustWaitLoad()
	fmt.Println("Page loading start")
	time.Sleep(5 * time.Second)
	fmt.Println("Page loading end")
	
	fmt.Println("Page scroll start")
	page.Mouse.Scroll(0, 99999, 5)
	
	time.Sleep(10 * time.Second)
	fmt.Println("Page scroll end")

	// Scroll to bottom

	// Get all post links
	links := page.MustElements("a")
	postLinks := []string{}

	for _, link := range links {
		hrefPtr, _ := link.Attribute("href")
		if hrefPtr != nil && strings.Contains(*hrefPtr, "/districtupdates/p/") {
			postLinks = append(postLinks, "https://www.instagram.com"+*hrefPtr)
		}
	}

	fmt.Println("Pages scanned : ", len(postLinks) )

	// create page scanner
	pageScanner := getPageScanner(postLinks, browser)

	var wg sync.WaitGroup
	wg.Add(1)
	// scan pages
	pageScanner.scanPages(&wg)

	wg.Wait()
}
