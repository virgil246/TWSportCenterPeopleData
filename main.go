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

	queryeach(loc)

}
func queryeach(loc []string) {
	var data []interface{}
	locPair := getchinese(loc)

	for _, l := range loc {
		//TODO go routine

		people, ok := fetch(l)
		if ok {
			// fmt.Println(people)
			people.ChLoc = locPair[l]
			data = append(data, people)
		} else {
			continue
		}

	}
	// fmt.Println(data)
	out, _ := json.MarshalIndent(data, "", " ")
	fileWrite(out)

}
func fileWrite(data []byte) {
	// fmt.Println(string(data))
	// timestamp := time.Now().Format("20060102_15-04-05")
	filename := "output.json"
	ioutil.WriteFile(filename, data, os.ModePerm)
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
func fetch(loc string) (People, bool) {
	var people People
	var link url.URL
	link.Scheme = "https"
	link.Path = "api"

	link.Host = loc + "sc.cyc.org.tw"
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(link.String())
	if err != nil {
		fmt.Println(err.Error())
		return people, false
	}
	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println(string(sitemap))
	people.Loc = loc
	json.Unmarshal(sitemap, &people)

	return people, true
}

type People struct {
	Loc   string
	ChLoc string
	Gym   Nums `json:"gym"`
	Swim  Nums `json:"swim"`
	Ice   Nums `json:"ice"`
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
