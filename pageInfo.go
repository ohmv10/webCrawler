package main

// type PageInfo struct {
// 	url string
// 	likes, caption string
// 	hashtags,profileTags []string
// }

type PageInfo struct {
	Url         string   `bson:"url"`
	Likes       string   `bson:"likes"` // We'll keep it string for now due to commas
	Caption     string   `bson:"caption"` // Not present in sample, so can be optional or omitted
	Hashtags    []string `bson:"hashtags"`
	ProfileTags []string `bson:"profileTags"`
}
