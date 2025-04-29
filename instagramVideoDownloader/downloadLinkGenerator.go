package instagramvideodownloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)


type InstagramGraphQLResponse struct {
	Data struct {
		XDTShortcodeMedia struct {
			IsVideo  bool   `json:"is_video"`
			VideoURL string `json:"video_url"`
		} `json:"xdt_shortcode_media"`
	} `json:"data"`
}

func DownloadLinkGenerator(reelShortcode string)(string, error){

	if reelShortcode == ""  {
		return "", errors.New("reelShortcode is empty")
	}
	resp, err := getInstagramPostGraphQL(reelShortcode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw body:", string(bodyBytes))
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var igResp InstagramGraphQLResponse
		err := json.NewDecoder(resp.Body).Decode(&igResp)
		fmt.Println("igResp: ", igResp)

		

		if  err != nil {
			return  "", errors.New("failed to decode response")
		}
		if igResp.Data.XDTShortcodeMedia.VideoURL == "" {
			return "", errors.New("video URL not found")
		}

		if !igResp.Data.XDTShortcodeMedia.IsVideo {
			return "", errors.New("post is not a video")
		}
		igResp.Data.XDTShortcodeMedia.VideoURL = url.QueryEscape(igResp.Data.XDTShortcodeMedia.VideoURL)
		return igResp.Data.XDTShortcodeMedia.VideoURL, nil
		
	case http.StatusNotFound:
		return "", errors.New("post not found")
	case http.StatusTooManyRequests, http.StatusUnauthorized:
		return "", errors.New("too many requests, try again later")
	default:
		return "", errors.New("unexpected status code")
	}
}

func getInstagramPostGraphQL(shortcode string) (*http.Response, error) {
	requestURL := "https://www.instagram.com/graphql/query"
	body := generateRequestBody(shortcode)

	req, err := http.NewRequest("POST", requestURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-FB-Friendly-Name", "PolarisPostActionLoadPostQueryQuery")
	req.Header.Set("X-BLOKS-VERSION-ID", "0d99de0d13662a50e0958bcb112dd651f70dea02e1859073ab25f8f2a477de96")
	req.Header.Set("X-CSRFToken", "uy8OpI1kndx4oUHjlHaUfu")
	req.Header.Set("X-IG-App-ID", "1217981644879628")
	req.Header.Set("X-FB-LSD", "AVrqPT0gJDo")
	req.Header.Set("X-ASBD-ID", "359341")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", fmt.Sprintf("https://www.instagram.com/p/%s/", shortcode))	

	client := &http.Client{}
	return client.Do(req)
}


func generateRequestBody(shortcode string) string {
	variables := map[string]interface{}{
		"shortcode":               shortcode,
		"fetch_tagged_user_count": nil,
		"hoisted_comment_id":      nil,
		"hoisted_reply_id":        nil,
	}
	variablesJSON, _ := json.Marshal(variables)

	params := url.Values{
		"av":                       {"0"},
		"__d":                      {"www"},
		"__user":                   {"0"},
		"__a":                      {"1"},
		"__req":                    {"b"},
		"__hs":                     {"20183.HYP:instagram_web_pkg.2.1...0"},
		"dpr":                      {"3"},
		"__ccg":                    {"GOOD"},
		"__rev":                    {"1021613311"},
		"__s":                      {"hm5eih:ztapmw:x0losd"},
		"__hsi":                    {"7489787314313612244"},
		"__dyn":                    {"7xeUjG1mxu1syUbFp41twpUnwgU7SbzEdF8aUco2qwJw5ux609vCwjE1EE2Cw8G11wBz81s8hwGxu786a3a1YwBgao6C0Mo2swtUd8-U2zxe2GewGw9a361qw8Xxm16wa-0oa2-azo7u3C2u2J0bS1LwTwKG1pg2fwxyo6O1FwlA3a3zhA6bwIxe6V8aUuwm8jwhU3cyVrDyo"},
		"__csr":                    {"goMJ6MT9Z48KVkIBBvRfqKOkinBtG-FfLaRgG-lZ9Qji9XGexh7VozjHRKq5J6KVqjQdGl2pAFmvK5GWGXyk8h9GA-m6V5yF4UWagnJzazAbZ5osXuFkVeGCHG8GF4l5yp9oOezpo88PAlZ1Pxa5bxGQ7o9VrFbg-8wwxp1G2acxacGVQ00jyoE0ijonyXwfwEnwWwkA2m0dLw3tE1I80hCg8UeU4Ohox0clAhAtsM0iCA9wap4DwhS1fxW0fLhpRB51m13xC3e0h2t2H801HQw1bu02j-"},
		"__comet_req":              {"7"},
		"lsd":                      {"AVrqPT0gJDo"},
		"jazoest":                  {"2946"},
		"__spin_r":                 {"1021613311"},
		"__spin_b":                 {"trunk"},
		"__spin_t":                 {"1743852001"},
		"__crn":                    {"comet.igweb.PolarisPostRoute"},
		"fb_api_caller_class":      {"RelayModern"},
		"fb_api_req_friendly_name": {"PolarisPostActionLoadPostQueryQuery"},
		"variables":                {string(variablesJSON)},
		"server_timestamps":        {"true"},
		"doc_id":                   {"8845758582119845"},
	}

	return params.Encode()
}
