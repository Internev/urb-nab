package main

import (
  "testing"
)

func Test_GetStringFromUrl(t *testing.T) {
  result := getStringFromURL("https://www.internev.com/")

  if len(result) < 1 {
    t.Errorf("Result has no data!")
  }
}
