package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3{
		fmt.Println("usage : ./出馬表変換 \"url\" (filename)")
		return
	}
	url := os.Args[1]
	var filename string
	if len(os.Args) == 3{
		filename = os.Args[2]
	}else{
		filename = "出馬表.csv"
	}
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"馬番","馬名","オッズ","人気","馬齢","斤量","騎手"}
	err = writer.Write(header)
	if err != nil {
		log.Fatal("Cannot write header to file", err)
	}
	data := get_horse_card(url)
	for _, row := range data {
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write row to file", err)
		}
	}
}

func get_horse_card(url string) [][]string{
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	utf8Reader, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		log.Fatal(err)
	}
	var names []string
	doc.Find(".name_line, .jockey").Each(func(i int, s *goquery.Selection) {
		s.Find(".name, .num ,.pop_rank ,.weight, .jockey, .age").Each(func(j int, t *goquery.Selection) {
			if t.HasClass("pop_rank") || t.HasClass("weight") {
				text := t.Text()
				re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)
				number := re.FindString(text)
				names = append(names, number)
			} else {
				name := t.Text()
				names = append(names, name)
			}
		})
	})
	var data [][]string
	count := 1
	for i := 0; i < len(names)-1; i += 6 {
		end := i + 6
		if end > len(names) {
			end = len(names)
		}
		row := []string{strconv.Itoa(count)}
		row = append(row, names[i:end]...)
    	data = append(data, row)
		count++
	}
	return(data)
}
