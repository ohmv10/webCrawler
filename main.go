package main

import (
	"fmt"
	"instagramVideoDownloader"
	"path"
	"strings"
	// "sync"
	// "os"

	"github.com/go-rod/rod"
)

func pageCloser(page *rod.Page){
	fmt.Println("page close triggered")
	page.Close()

}

// x1lliihq x1n2onr6 xh8yej3 x4gyw5p x11i5rnm x1ntc13c x9i3mqj x2pgyrj
var mediaClassName string = "x1lliihq.x1n2onr6.xh8yej3.x4gyw5p.x11i5rnm.x1ntc13c.x9i3mqj.x2pgyrj"
// xdj266r x11i5rnm xat24cr x1mh8g0r xexx8yu x4uap5 x18d9i69 xkhd6sd x1hl2dhg x16tdsg8 x1vvkbs
var likesClassName string = "span.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x1hl2dhg.x16tdsg8.x1vvkbs"
// var postInfoClassName string = "h1._ap3a._aaco._aacu._aacx._aad7._aade"
var postInfoClassName string = "h1._ap3a._aaco._aacu._aacx._aad7._aade"

func main() {
	// pageScanner := createPageScanner("https://www.instagram.com/accounts/login/")
	// defer pageScanner.closeMainPage()

	// // Login
	// pageScanner.loginInstagram(os.Getenv("INSTA_USERNAME"),os.Getenv("INSTA_PASSWORD"))

	// // Navigate to profile
	// // pageScanner.navigateToProfileWithURL("https://www.instagram.com/districtupdates/")
	// pageScanner.navigateToProfileWithURL("https://www.instagram.com/districtupdates/")

	// //scroll to end
	// pageScanner.scrollToEnd(mediaClassName,false)

	// // Get all post links
	// pageScanner.updatePostSlice(mediaClassName)

	// // scan pages
	// var wg sync.WaitGroup
	// wg.Add(1)
	// pageScanner.scanPages(&wg)
	// wg.Wait()

	// for _, postInfo := range pageScanner.pageInfos {
	// 	fmt.Printf("url : %s \n",postInfo.Url)
	// 	fmt.Printf("likes : %s \n",postInfo.Likes)
	// 	fmt.Printf("hashtags : %s \n",postInfo.Hashtags)
	// 	fmt.Printf("profileTags : %s \n",postInfo.ProfileTags)
	// 	fmt.Println("--------------------")
	// }

	// uri := "mongodb://localhost:27017"
	// mongoInstance := createInstance(uri)
	// mongoInstance.connectDB("insta")
	// mongoInstance.connectCollection("pravah")
	// mongoInstance.insertData(pageScanner.pageInfos)
	
	// download videos test
	
	uri := "mongodb://localhost:27017"
	mongoInstance := createInstance(uri)
	mongoInstance.connectDB("insta")
	mongoInstance.connectCollection("pravah")
	var pageScanner PageScanner
	mongoInstance.loadPageInfosFromDB(&pageScanner)

	for _, pageInfo := range pageScanner.pageInfos{

		reelShortCode := path.Base(strings.TrimSuffix(pageInfo.Url, "/"))
		fmt.Println("Reel Shortcode: ", reelShortCode)
		downloadURL,err := instagramvideodownloader.DownloadLinkGenerator(reelShortCode)
		if err != nil {
			fmt.Println("Error generating download link:", err)
			continue
		}
		err = instagramvideodownloader.DownloadProxy(downloadURL, fmt.Sprintf("%s.mp4",pageInfo.Url))

		if err != nil {
			fmt.Println("Error downloading video:", err)
			continue
		}
	}
}
