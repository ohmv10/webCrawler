package main

import (
	"fmt"
	"sync"
	// "time"
	"os"
	"github.com/go-rod/rod"
)

func pageCloser(page *rod.Page){
	fmt.Println("page close triggered")
	page.Close()

}


func main() {

	pageScanner := createPageScanner("https://www.instagram.com/accounts/login/")
	defer pageScanner.closeMainPage()

	// Login
	pageScanner.loginInstagram(os.Getenv("INSTA_USERNAME"),os.Getenv("INSTA_PASSWORD"))

	// Navigate to profile
	pageScanner.navigateToProfileWithURL("https://www.instagram.com/districtupdates/")

	//scroll to end
	pageScanner.scrollToEnd()

	// Get all post links
	pageScanner.updatePostSlice("x1lliihq.x1n2onr6.xh8yej3.x4gyw5p.x11i5rnm.x1ntc13c.x9i3mqj.x2pgyrj")

	// scan pages
	var wg sync.WaitGroup
	wg.Add(1)
	pageScanner.scanPages(&wg)
	wg.Wait()

	for _, postInfo := range pageScanner.pageInfos {
		fmt.Printf("url : %s \n",postInfo.url)
		fmt.Printf("likes : %s \n",postInfo.likes)
		fmt.Printf("hashtags : %s \n",postInfo.hashtags)
		fmt.Printf("profileTags : %s \n",postInfo.profileTags)
		fmt.Println("--------------------")
	}
}