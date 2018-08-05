package main

import (
  "fmt"
  "log"
  // "io/ioutil"
  // "path/filepath"
  "os"

  "github.com/PuerkitoBio/goquery"
)

func check(e error) {
  if e != nil {
    log.Fatal(e)
  }
}

func main() {
  links := prepLinks("https://www.urbandictionary.com")
  fmt.Println("links:", links)
  for _, link := range(links) {
    fmt.Println("link:", link)
    saveDefinition(link)
  }

  fmt.Println("Done?!")
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

func saveDefinition(url string) {
  doc, err := goquery.NewDocument(url)
  check(err)

  def := doc.Find(".def-panel").First()
  termTitle := def.Find("a.word").Text()
  termDef := def.Find(".meaning").Text()
  termExample := def.Find(".example").Text()

  definition := term{termTitle, termDef, termExample}
  fmt.Println("I have made a struct:", definition)
  definition.save("/scraped")
}
