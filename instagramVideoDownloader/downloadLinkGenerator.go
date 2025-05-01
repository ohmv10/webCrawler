package instagramvideodownloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"

	// "io"
	"net/http"
	// "net/http/httptrace"
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

type Response struct {
	Data struct {
		XdtShortcodeMedia struct {
			IsVideo  bool   `json:"is_video"`
			VideoURL string `json:"video_url"`
			EdgeMediaToCaption struct {
				Edges []struct {
					Node struct {
						Text string `json:"text"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_caption"`
			EdgeMediaPreviewLike struct {
				Count int `json:"count"`
			} `json:"edge_media_preview_like"`
			EdgeMediaToTaggedUser struct {
				Edges []struct {
					Node struct {
						User struct {
							Username string `json:"username"`
							FullName string `json:"full_name"`
						} `json:"user"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_tagged_user"`
			Owner struct {
				Username string `json:"username"`
			} `json:"owner"`
		} `json:"xdt_shortcode_media"`
	} `json:"data"`
	Hashtags []string `json:"hashtags"`
	Caption  string   `json:"caption"`
}


// extractHashtags extracts hashtags from text
func extractHashtags(text string) []string {
	var hashtags []string
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			hashtags = append(hashtags, word)
		}
	}
	return hashtags
}

func FillData(igResp Response){
	var captions []string
	for _, edge := range igResp.Data.XdtShortcodeMedia.EdgeMediaToCaption.Edges {
		captions = append(captions, edge.Node.Text)
	}
	igResp.Caption = strings.Join(captions, "\n\n") // Join captions with double newline

	// Extract hashtags from captions
	for _, edge := range igResp.Data.XdtShortcodeMedia.EdgeMediaToCaption.Edges {
		igResp.Hashtags = append(igResp.Hashtags, extractHashtags(edge.Node.Text)...)
	}

}

// DownloadLinkGenerator generates download link and populates Response
func DownloadLinkGenerator(reelShortcode, folder string) (Response, error) {
	reelShortcode = path.Base(strings.TrimSuffix(reelShortcode, "/"))

	if reelShortcode == "" {
		return Response{}, errors.New("reelShortcode is empty")
	}

	resp, err := getInstagramPostGraphQL(reelShortcode)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return Response{}, errors.New("failed to read response body")
		}

		var igResp Response
		err = json.Unmarshal(bodyBytes, &igResp)
		FillData(igResp)
		if err != nil {
			return igResp, errors.New("failed to decode response")
		}

		if igResp.Data.XdtShortcodeMedia.VideoURL == "" {
			return igResp, errors.New("video URL not found")
		}

		if !igResp.Data.XdtShortcodeMedia.IsVideo {
			return igResp, errors.New("post is not a video")
		}

		fmt.Println(reelShortcode)
		if DownloadProxy(igResp.Data.XdtShortcodeMedia.VideoURL, fmt.Sprintf("%s/%s.mp4", folder, reelShortcode)) != nil {
			return igResp, errors.New("failed to download video")
		}

		return igResp, nil 

	case http.StatusNotFound:
		return Response{}, errors.New("post not found")
	case http.StatusTooManyRequests, http.StatusUnauthorized:
		return Response{}, errors.New("too many requests, try again later")
	default:
		return Response{}, errors.New("unexpected status code")
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
	req.Header.Set("Referer", fmt.Sprintf("https://www.instagram.com/%s/", shortcode))

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
