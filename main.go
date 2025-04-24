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


var mediaClassName string = "x1lliihq.x1n2onr6.xh8yej3.x4gyw5p.x11i5rnm.x1ntc13c.x9i3mqj.x2pgyrj"
var likesClassName string = "span.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x1hl2dhg.x16tdsg8.x1vvkbs"
var postInfoClassName string = "h1._ap3a._aaco._aacu._aacx._aad7._aade"
// var postInfoClassName string = "span.x193iq5w.xeuugli.x1fj9vlw.x13faqbe.x1vvkbs.xt0psk2.x1i0vuye.xvs91rp.xo1l8bm.x5n08af.x10wh9bi.x1wdrske.x8viiok.x18hxmgj"

func main() {
	pageScanner := createPageScanner("https://www.instagram.com/accounts/login/")
	defer pageScanner.closeMainPage()

	// Login
	pageScanner.loginInstagram(os.Getenv("INSTA_USERNAME"),os.Getenv("INSTA_PASSWORD"))

	// Navigate to profile
	// pageScanner.navigateToProfileWithURL("https://www.instagram.com/districtupdates/")
	pageScanner.navigateToProfileWithURL("https://www.instagram.com/iitbbs.pravaah/")

	//scroll to end
	pageScanner.scrollToEnd(mediaClassName,false)

	// Get all post links
	pageScanner.updatePostSlice(mediaClassName)

	// scan pages
	var wg sync.WaitGroup
	wg.Add(1)
	pageScanner.scanPages(&wg)
	wg.Wait()

	for _, postInfo := range pageScanner.pageInfos {
		fmt.Printf("url : %s \n",postInfo.Url)
		fmt.Printf("likes : %s \n",postInfo.Likes)
		fmt.Printf("hashtags : %s \n",postInfo.Hashtags)
		fmt.Printf("profileTags : %s \n",postInfo.ProfileTags)
		fmt.Println("--------------------")
	}

	uri := "mongodb://localhost:27017"
	mongoInstance := createInstance(uri)
	mongoInstance.connectDB("insta")
	mongoInstance.connectCollection("pravah")
	mongoInstance.insertData(pageScanner.pageInfos)
}
