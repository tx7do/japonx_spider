package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	START_PAGE = 94
	END_PAGE   = 1402
)

/**
* 采集内容
 */
type Content struct {
	id          string ///< ID
	title       string ///< 标题
	desc        string ///< 描述
	content_url string ///< 内容入口地址
	cover_url   string ///< 封面地址
	thumb_url   string ///< 缩略图地址
	video_url   string ///< 视频地址
}

/// 采集单个页面的内容
func GetContent(page int, url string) Content {

	var content Content

	content.id = strconv.Itoa(page)
	content.content_url = url
	//fmt.Println(url)

	// 获取页面
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	// 封面图
	doc.Find("#bxslider img").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("src")
		content.thumb_url = link
	})
	//fmt.Println(content.thumb_url)

	// 标题
	tt := doc.Find("title")
	content.title = tt.Text()
	//fmt.Println(tt.Text())

	// 描述文字
	tx := doc.Find(".tx-comment")
	content.desc = tx.Text()
	//fmt.Println(tx.Text())

	// 视频地址
	s := doc.Find("script")
	//fmt.Println(s.Text())

	pat := `(\(function\(p).*(\))`
	re := regexp.MustCompile(pat)
	res := re.FindAllString(s.Text(), -1)
	//fmt.Println(res)

	vm := otto.New()
	value, err := vm.Run(res[0])
	if err != nil {
		log.Fatal(url, err)
	}
	//fmt.Println(value)

	pat = `(https)(.*)mp4`
	re = regexp.MustCompile(pat)
	res = re.FindAllString(value.String(), -1)
	if len(res) != 0 {
		content.video_url = res[0]
	}
	//fmt.Println(content.video_url)

	time.Sleep(time.Second * 1)

	//fmt.Println(content)
	return content
}

func main() {
	var contents []Content

	var startPage = START_PAGE
	var endPage = END_PAGE
	var url string
	for i := startPage; i <= endPage; i++ {
		url = fmt.Sprintf("https://www.japonx.net/portal/index/detail/id/%d.html", i)
		contents = append(contents, GetContent(i, url))
	}

	for _, s := range contents {
		fmt.Println("===========================")
		fmt.Println("页面地址:", s.content_url)
		fmt.Println("标题:", s.title)
		fmt.Println("缩略图:", s.thumb_url)
		fmt.Println("视频:", s.video_url)
	}

	fileName := "japonx.csv"
	buf := new(bytes.Buffer)
	r2 := csv.NewWriter(buf)
	for _, s := range contents {
		ss := make([]string, 4)
		ss[0] = s.id
		ss[1] = s.title
		ss[2] = s.thumb_url
		ss[3] = s.video_url
		r2.Write(ss)
		r2.Flush()
	}
	fout, err := os.Create(fileName)
	defer fout.Close()
	if err != nil {
		fmt.Println(fileName, err)
		return
	}
	fout.WriteString(buf.String())
}
