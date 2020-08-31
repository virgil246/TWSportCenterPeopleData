package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	loc := []string{"cs", "ng", "xy", "da", "ws", "nh",
		"lzc", "tcc", "xzc", "yhc", "cmc", "zlc", "tyc", "zgc", "lkc"}
	// locPair := getchinese(loc)
	var data []interface{}

	for _, l := range loc {
		//TODO go routine
		data = append(data, fetch(l))

	}

	out, _ := json.MarshalIndent(data, "", " ")
	fmt.Println(string(out))
	timestamp := time.Now().Format("20060102_15-04-05")
	filename := timestamp + ".json"
	ioutil.WriteFile(filename, out, os.ModePerm)

}
func getchinese(loc []string) (place map[string]string) {
	place = make(map[string]string)
	var link url.URL
	link.Scheme = "https"
	for _, l := range loc {
		link.Host = l + "sc.cyc.org.tw"
		resp, err := http.Get(link.String())
		if err != nil {

		}

		place[l] = getTitle(resp)
	}

	return place
}
func getTitle(res *http.Response) (title string) {
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	title = doc.Find("title").Text()
	return title

}
func fetch(loc string) (people People) {
	var link url.URL
	link.Scheme = "https"
	link.Path = "api"

	link.Host = loc + "sc.cyc.org.tw"
	resp, err := http.Get(link.String())
	if err != nil {

	}
	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	// fmt.Println(string(sitemap))
	people.Loc = loc
	json.Unmarshal(sitemap, &people)

	return people
}

type People struct {
	Loc  string
	Gym  Nums `json:"gym"`
	Swim Nums `json:"swim"`
	Ice  Nums `json:"ice"`
}
type Nums struct {
	Now int
	Max int
	// Unknown int
}

func (p *Nums) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(err)
		return err
	}
	p.Now, _ = strconv.Atoi(v[0].(string))
	p.Max, _ = strconv.Atoi(v[1].(string))
	// p.Unknown, _ = strconv.Atoi(v[2].(string))

	return nil
}
