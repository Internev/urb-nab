package main

import (
  "fmt"
  "log"
  "sync"
  "time"
  "strings"
  "net/http"
  "archive/zip"
  "io"
  "io/ioutil"
  // "encoding/json"
  // "path/filepath"
  "os"
  // "reflect"

  "github.com/PuerkitoBio/goquery"
)

var wg sync.WaitGroup
var wg2 sync.WaitGroup


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
  // go serveWeb()
  //
  // pollInterval := 12
  //
	// timerCh := time.Tick(time.Duration(pollInterval) * time.Hour)
  //
	// for range timerCh {
	// 	grab()
	// }
  grab()
}

func serveWeb() {
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/give", giveHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func grab() {
  links := prepLinks("https://www.urbandictionary.com")
  for _, link := range(links) {
    wg.Add(1)
    go saveDefinition(link)
  }

  wg.Wait()
  fmt.Println("Grabbed at", time.Now())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi, how's it going?")
}

func giveHandler(w http.ResponseWriter, r *http.Request) {
  newfile, err := os.Create("urb.zip")
  check(err)
  defer newfile.Close()

  zipWriter := zip.NewWriter(newfile)
  defer zipWriter.Close()

  files, err := ioutil.ReadDir("./scraped")
  check(err)
  for _, f := range files {
    makeZip(zipWriter, f.Name())
  }

  wd, err := os.Getwd()
  check(err)


  fmt.Println("Threads finished, we good.")
  fmt.Fprintf(w, "Have you some data.")
  w.Header().Set("Content-Type", "application/zip")
  w.Header().Set("Content-Disposition", "attachment; filename='urb.zip'")
  http.ServeFile(w, r, wd + "/urb.zip")
}

func makeZip(z *zip.Writer, fName string) {
  zipFile, err := os.Open("./scraped/" + fName)
  check(err)
  defer zipFile.Close()

  info, err := zipFile.Stat()
  check(err)

  header, err := zip.FileInfoHeader(info)
  check(err)

  header.Method = zip.Deflate

  writer, err := z.CreateHeader(header)
  check(err)

  _, err = io.Copy(writer, zipFile)
  check(err)

  fmt.Println(fName)
  return
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

  os.MkdirAll(path, os.ModePerm)

  f, err := os.Create(path + "/" + strings.Replace(t.term, "/", "-", -1) + ".txt")
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
  definition.save("/scraped")
}
