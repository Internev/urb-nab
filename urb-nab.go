package main

import (
  "fmt"
  "log"
  "sync"
  "time"
  // "net/http"
  // "encoding/json"
  // "io/ioutil"
  // "path/filepath"
  "os"
  // "reflect"

  "github.com/PuerkitoBio/goquery"
  // "github.com/jinzhu/gorm"
  // _ "github.com/jinzhu/gorm/dialects/postgres"
)

var wg sync.WaitGroup

type Entry struct {
  Term string
  Def string
  Example string
}

func check(e error) {
  if e != nil {
    log.Fatal(e)
  }
}

func main() {
  pollInterval := 12

	timerCh := time.Tick(time.Duration(pollInterval) * time.Hour)

	for range timerCh {
		grab()
	}
}

func grab() {

  links := prepLinks("https://www.urbandictionary.com")
  for _, link := range(links) {
    wg.Add(1)
    go saveDefinition(link)
  }

  wg.Wait()
  fmt.Println("Grabbed at", time.Now())

  // db, err := gorm.Open("postgres", "host=localhost port=5432 user=urb dbname=urbdic password=badger")
  // check(err)
  // defer db.Close()
  //
  // db.AutoMigrate(&Entry{})
  //
  // http.HandleFunc("/grab", grabHandler)
  // http.HandleFunc("/give", giveHandler)
  // log.Fatal(http.ListenAndServe(":8080", nil))
}

// func grabHandler(w http.ResponseWriter, r *http.Request) {
//   db, err := gorm.Open("postgres", "host=localhost port=5432 user=urb dbname=urbdic password=badger")
//   check(err)
//   defer db.Close()
// }
//
// func giveHandler(w http.ResponseWriter, r *http.Request) {
//   db, err := gorm.Open("postgres", "host=localhost port=5432 user=urb dbname=urbdic password=badger")
//   check(err)
//   defer db.Close()
//
//   var results []Entry
//   db.Find(&results)
//
//   b, err := json.Marshal(results)
//   check(err)
//
//   w.Header().Set("Content-Type", "application/json")
//   w.Write(b)
// }

func prepLinks(url string) []string {
  doc, err := goquery.NewDocument(url)
  check(err)

  var links []string

  doc.Find(".trending-link").Each(func(index int, s *goquery.Selection) {
      link, _ := s.Attr("href")
      links = append(links, url + link)
    })

  return links
}

type term struct {
  term string
  def string
  example string
}

func (t term) save(path string) {
  // Assemble Text
  var text string
  text += t.term + "\n"
  text += t.def + "\n"
  text += t.example + "\n"

  wd, err := os.Getwd()
  check(err)

  path = wd + path

  os.MkdirAll(path, os.ModePerm)

  f, err := os.Create(path + "/" + t.term + ".txt")
  check(err)

  defer f.Close()

  _, err = f.WriteString(text)
  check(err)

  f.Sync()
}

func saveDefinition(url string) {
  defer wg.Done()
  doc, err := goquery.NewDocument(url)
  check(err)

  def := doc.Find(".def-panel").First()
  termTitle := def.Find("a.word").Text()
  termDef := def.Find(".meaning").Text()
  termExample := def.Find(".example").Text()

  definition := term{termTitle, termDef, termExample}
  // dbEntry := Entry{Term: termTitle, Def: termDef, Example: termExample}
  //
  // db.Where(Entry{Term: termTitle}).Attrs(dbEntry).FirstOrCreate(&dbEntry)
  // db.Where(Entry{Term: termTitle}).Attrs(Entry{Term: termTitle, Def: termDef, Example: termExample}).FirstOrCreate(&Entry)
  // db.Create(&dbEntry)
  definition.save("/scraped")
}
