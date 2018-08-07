package main

import (
  "fmt"
  "log"
  "sync"
  // "io/ioutil"
  // "path/filepath"
  "os"
  // "reflect"

  "github.com/PuerkitoBio/goquery"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

var wg sync.WaitGroup

type termDB struct {
  gorm.Model
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
  db, err := gorm.Open("postgres", "host=localhost port=5432 user=urb dbname=urbdic password=badger")
  check(err)
  defer db.Close()

  db.AutoMigrate(&termDB{})

  links := prepLinks("https://www.urbandictionary.com")
  for _, link := range(links) {
    wg.Add(1)
    go saveDefinition(link, db)
  }

  wg.Wait()
  fmt.Println("All Done.")
}

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

  f, err := os.Create(path + "/" + t.term + ".txt")
  check(err)

  defer f.Close()

  _, err = f.WriteString(text)
  check(err)

  f.Sync()
}

func saveDefinition(url string, db *gorm.DB) {
  defer wg.Done()
  doc, err := goquery.NewDocument(url)
  check(err)

  def := doc.Find(".def-panel").First()
  termTitle := def.Find("a.word").Text()
  termDef := def.Find(".meaning").Text()
  termExample := def.Find(".example").Text()

  definition := term{termTitle, termDef, termExample}
  db.Create(&termDB{Term: termTitle, Def: termDef, Example: termExample})
  definition.save("/scraped")
}
