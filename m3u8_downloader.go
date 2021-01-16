package utils

import (
	"bytes"
	"errors"
	"strings"

	"github.com/levigross/grequests"
)

var (
	ErrRequestsFailed = errors.New("request not ok")

	ErrInvalidSeed = errors.New("seed is invalid")
)

func init() {

}

// download seed and parse it to urls
func downloadSeed(url string, ro *grequests.RequestOptions) []string {
	res, err := grequests.Get(url, ro)
	ErrPanic(err)
	if res == nil || !res.Ok {
		ErrPanic(ErrRequestsFailed)
	}
	content := res.String()
	lines := strings.Split(content, "\n")
	var urls []string
	for i := range lines {
		if len(lines[i]) != 0 && !strings.HasPrefix(lines[i], "#") {
			urls = append(urls, lines[i])
		}
	}
	if len(urls) == 0 {
		return urls
	}
	var completeUrls []string
	var baseUrl string

	for i := range urls {
		if strings.HasPrefix(urls[i], "http") {
			baseUrl = ""
		} else if strings.HasPrefix(urls[i], "/") {
			baseUrl = getHostname(url)
		} else {
			baseUrl = getPrefix(url)
		}
		completeUrls = append(completeUrls, joinUrl(baseUrl, urls[i]))
	}

	return completeUrls
}

// eg: https://github.com/1/1/2 --> https://github.com
func getHostname(url string) string {
	var index, count int
	for i := range url {
		if url[i] == '/' {
			count++
		}
		if count == 3 {
			index = i
			break
		}
	}
	if count == 3 {
		return url[:index]
	} else {
		return ""
	}
}

// eg: https://github.com/1/1/2 --> https://github.com/1/1/
func getPrefix(url string) string {
	var index int
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			index = i
			break
		}
	}
	if index == 0 {
		return ""
	}
	return url[:index+1]
}

func joinUrl(url1, url2 string) string {
	if len(url1) == 0 {
		return url2
	}
	if len(url2) == 0 {
		return url1
	}
	var buffer bytes.Buffer
	var b1 = url1[len(url1)-1] == '/'
	var b2 = url2[0] == '/'
	if (b1 && !b2) || (!b1 && b2) {
		buffer.WriteString(url1)
		buffer.WriteString(url2)
	} else if b1 && b2 {
		buffer.WriteString(url1)
		buffer.WriteString(url2[1:])
	} else {
		buffer.WriteString(url1)
		buffer.WriteString("/")
		buffer.WriteString(url2)
	}
	return buffer.String()
}

// get ts url according to a seed url, no matter how many seeds during the process!
// the process is not done concurrently
func getTsUrls(url string, ro *grequests.RequestOptions) []string {
	var originUrls, newUrls []string
	originUrls = append(originUrls, url)
	for {
		if len(originUrls) == 0 {
			ErrPanic(ErrInvalidSeed)
		}
		// not seed, just return
		if !strings.Contains(originUrls[0], "m3u8") {
			return originUrls
		}
		for i := range originUrls {
			temp := downloadSeed(originUrls[i], ro)
			for j := range temp {
				newUrls = append(newUrls, temp[j])
			}
		}
		if len(newUrls) == 0 {
			ErrPanic(ErrInvalidSeed)
		}
		var base string

		for i := range newUrls {
			if strings.HasPrefix(newUrls[i], "http") {
				base = ""
			} else if strings.HasPrefix(newUrls[i], "/") {
				base = getHostname(originUrls[0])
			} else {
				base = getPrefix(originUrls[0])
			}
			newUrls[i] = joinUrl(base, newUrls[i])
		}
		originUrls = newUrls
		var dummy []string
		newUrls = dummy
	}
}

func GetTsUrls(url string, ro *grequests.RequestOptions) []string {
	return getTsUrls(url, ro)
}
