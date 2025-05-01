package main

import (
	"fmt"
	"instagramVideoDownloader"
	"os"
	// "path"
	// "strings"
	"sync"

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
var folder string = "cocomelon"
var pageURL string = "https://www.instagram.com/cocomelon/"

var folders []string = []string{
	"toocool",
	"yesminister",
	"pattyshukla",
	"minigolf",
	"souljams",
	"worldofwhisky",
	"pindezaika",
	"rageroom",
	"zindagikekhayal",
	"caravaggio",
	"bestcomedylineup",
	"bestcomedylineup",
	"shapingtextures",
	"taarukraina",
	"ankurtewari",
	"ankurtewari",
	"potteryastherapy",
	"arrahman",
}

var urls []string= []string{
	"https://www.instagram.com/bandraborn/",
	"https://www.instagram.com/yesministerbyessex/",	
	"https://www.instagram.com/club.loka/",	
	"https://www.instagram.com/minigolfmadness/",	
	"https://www.instagram.com/bira91taproom/",	
	"https://www.instagram.com/whiskysamba/",	
	"https://www.instagram.com/shangrilanewdelhi/",	
	"https://www.instagram.com/delhirageroom/",	
	"https://www.instagram.com/amandeep.khayal/",	
	"https://www.instagram.com/knmaindia/",	
	"https://www.instagram.com/madhurvirli/",	
	"https://www.instagram.com/pranavsharm_a/",
	"https://www.instagram.com/p/DGiQSxuNYuJ/",	
	"https://www.instagram.com/taarukraina/",	
	"https://www.instagram.com/ankurtewatia_/",	
	"https://www.instagram.com/ghalatfamily/",
	"https://www.instagram.com/naveenchhaya.16/",		
	"https://www.instagram.com/arrahman/",	
}
func main() {
	pageScanner := createPageScanner("https://www.instagram.com/accounts/login/")
	defer pageScanner.closeMainPage()
	// Login
	pageScanner.loginInstagram(os.Getenv("INSTA_USERNAME"),os.Getenv("INSTA_PASSWORD"))
	if len(folders) != len(urls) {
		fmt.Println("Error: folders and urls length mismatch")
		return
	}
	for i:=0; i<len(urls); i++ {
		folder = folders[i]
		pageURL = urls[i]


	// Navigate to profile
	// pageScanner.navigateToProfileWithURL("https://www.instagram.com/districtupdates/")
	pageScanner.navigateToProfileWithURL(pageURL)

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
	mongoInstance.connectCollection(folder)
	mongoInstance.insertData(pageScanner.pageInfos)

	var download_wg sync.WaitGroup

	for _, pageInfo := range pageScanner.pageInfos{
		download_wg.Add(1)
		func (url string)  {
			defer download_wg.Done()
			_,err := instagramvideodownloader.DownloadLinkGenerator(pageInfo.Url, folder)
			if err != nil {
				fmt.Println("Error generating download link:", err)
				return
			}
		}(pageInfo.Url)
	}
	download_wg.Wait()
	}
}

