package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "os"
  "path/filepath"

  "github.com/docker/go-plugins-helpers/volume"
)

type loVolDriver struct {
	path string
}

func newLoVolDriver(path string) loVolDriver {
  driver := loVolDriver{
    path: path,
  }

  return driver
}

/**
* @brief Create method
*
* @param Request
*/
func (driver loVolDriver) Create(request volume.Request) volume.Response {
  log.Printf("Creating volume %s\n", request)

  count := "500" // ==> 500M
  if _count, ok := request.Options["count"]; ok {
    count = _count
  }

  basedir := filepath.Join(driver.path, request.Name)
  mppath := filepath.Join(basedir, "_mountpoint")
  imgpath := filepath.Join(basedir, "volume.img")

  if err := os.MkdirAll(mppath, os.ModeDir); err != nil {
    return volume.Response{Err: err.Error()}
  }

  dd_args := map[string]string{
    "if": "/dev/zero",
    "of": imgpath,
    "bs": "1M",
    "count": count,
  }

  if err := execute("dd", dd_args); err != nil {
    return volume.Response{Err: err.Error()}
  }

  if err := _execute("mkfs.ext4", imgpath); err != nil {
    return volume.Response{Err: err.Error()}
  }

  return volume.Response{Err: ""}
}

/**
* @brief Remove method
*
* @param Request
*/
func (driver loVolDriver) Remove(request volume.Request) volume.Response {
  log.Printf("Removing volume %s\n", request)

  basedir := filepath.Join(driver.path, request.Name)
  //mppath := filepath.Join(basedir, "_mountpoint")

  //_execute("umount", mppath)
  if err:= os.RemoveAll(basedir); err != nil {
    return volume.Response{Err: err.Error()}
  }

  return volume.Response{Err: ""}
}

/**
* @brief Mount method
*
* @param Request
*/
func (driver loVolDriver) Mount(request volume.MountRequest) volume.Response {
  log.Printf("Mount volume %s\n", request)

  mppath := filepath.Join(driver.path, request.Name, "_mountpoint")
  imgpath := filepath.Join(driver.path, request.Name, "volume.img")

  log.Printf("Mounting volume %s on %s\n", request.Name, mppath)

  if err := _execute("mount", "-o", "loop", imgpath, mppath); err != nil {
    return volume.Response{Err: err.Error()}
  }

  return volume.Response{Err: "", Mountpoint: mppath}
}

/**
* @brief Path method
*
* @param Request
*/
func (driver loVolDriver) Path(request volume.Request) volume.Response {
  log.Printf("Path volume %s\n", request)

  mppath := filepath.Join(driver.path, request.Name, "_mountpoint")

  return volume.Response{Mountpoint: mppath}
}

/**
* @brief Unmount method
*
* @param Request
*/
func (driver loVolDriver) Unmount(request volume.UnmountRequest) volume.Response {
  log.Printf("Umount volume %s\n", request)

  mppath := filepath.Join(driver.path, request.Name, "_mountpoint")

  if err:= _execute("umount", mppath); err != nil {
    return volume.Response{Err: err.Error()}
  }

  return volume.Response{}
}

/**
* @brief Get method
*
* @param Request
*/
func (driver loVolDriver) Get(request volume.Request) volume.Response {
  log.Printf("Get volume %s\n", request)

  // FIXME
  mPoint := filepath.Join(driver.path, request.Name)

  if fi, err := os.Lstat(mPoint); err != nil || !fi.IsDir() {
    log.Println(err)
    return volume.Response{Err: fmt.Sprintf("%v not mounted", mPoint)}
  }

  return volume.Response{Volume: &volume.Volume{Name: request.Name, Mountpoint: mPoint}}
}

/**
* @brief List method
*
* @param Request
*/
func (driver loVolDriver) List(request volume.Request) volume.Response {
  log.Printf("List volume %s\n", request)

  var vols []*volume.Volume

  files, err := ioutil.ReadDir(driver.path)
  if err != nil {
    log.Println(err)
    return volume.Response{Err: err.Error()}
  }

  for _, entry := range files {
    if entry.IsDir() {
      vols = append(vols, &volume.Volume{Name: entry.Name(), Mountpoint: filepath.Join(driver.path, entry.Name())})
    }
  }

  return volume.Response{Volumes: vols}
}

/**
* @brief Capabilities method
*
* @param Request
*/
func (driver loVolDriver) Capabilities(request volume.Request) volume.Response {
  log.Printf("Capabilities volume %s\n", request)

  return volume.Response{Capabilities: volume.Capability{Scope: "global"}}
}
