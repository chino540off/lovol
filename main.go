package main

import (
  "flag"
  "fmt"
  "log"
  "os"

  "github.com/docker/go-plugins-helpers/volume"
)

var (
  // These fields are populated by govvv
  BuildDate  string
  GitCommit  string
  GitBranch  string
  GitState   string
  GitSummary string
)

func main() {
  path := flag.String("path", "/run/docker/lovol/mnt", "Path where Quobyte is mounted on the host")
  //group := flag.String("group", "root", "Group to create the unix socket")
  show_version := flag.Bool("version", false, "Shows version string")
  flag.Parse()

  if *show_version {
    fmt.Printf("BuildDate=%s\n", BuildDate)
    fmt.Printf("GitCommit=%s\n", GitCommit)
    fmt.Printf("GitBranch=%s\n", GitBranch)
    fmt.Printf("GitState=%s\n", GitState)
    fmt.Printf("GitSummary=%s\n", GitSummary)
    return
  }

  if err := os.MkdirAll(*path, 0555); err != nil {
    log.Println(err.Error())
  }

  driver := newLoVolDriver(*path)
  handler := volume.NewHandler(driver)
  log.Println(handler.ServeUnix("loVol", 0))
}
