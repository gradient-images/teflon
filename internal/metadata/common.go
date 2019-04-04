package metadata

import (
  "os"
  "log"
  "path/filepath"
)

const (
	teflonDirName = ".teflon"
	teflonDirMetaName = "_"
	teflonMetaExt = "._"
)

func Get(baseName string) *UserSection {
  baseInfo, err := os.Stat(baseName)
  if err != nil {
    log.Fatal(err)
  }

  var metaName string

  if baseInfo.IsDir() {
    metaName = filepath.Join(baseName, teflonDirName, teflonDirMetaName)
  } else {
    d, n := filepath.Split(baseName)
    metaName = filepath.Join(d, teflonDirName, n + teflonMetaExt)
  }

  log.Print(baseName, metaName)

  us := UserSection{}

  if _, err := os.Stat(metaName); os.IsNotExist(err) {
    log.Print("Meta file doesn't exists.")
    us.UserData = make(map[string]string)
  } else {
    log.Print("Meta file exists.")
  }

  return &us
}
