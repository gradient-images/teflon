// Copyright © 2019 Máté Birkás <gadfly16@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadata

import (
  "os"
  "log"
  "path/filepath"
  "io/ioutil"

  "github.com/golang/protobuf/proto"
)

const (
	teflonDirName = ".teflon"
	teflonDirMetaName = "_"
	teflonMetaExt = "._"
)

type Metadata struct {
  baseName, metaName string
  UserSection UserSection
}

func Get(baseName string) *Metadata {
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

  log.Println(baseName, metaName)

  us := UserSection{}

  if _, err := os.Stat(metaName); os.IsNotExist(err) {
    log.Print("Meta file doesn't exists.")
    us.UserData = make(map[string]string)
  } else {
    log.Print("Meta file exists.")
    in, err := ioutil.ReadFile(metaName)
    if err != nil {
      log.Fatalln("Error reading meta file:", err)
    }
    if err := proto.Unmarshal(in, &us); err != nil {
      log.Fatalln("Failed to parse meta file:", err)
    }
  }

  return &Metadata{baseName, metaName, us}
}

func (md Metadata) Sync() {
  out, err := proto.Marshal(&md.UserSection)
	if err != nil {
		log.Fatalln("Failed to encode metadata:", err)
	}
  createTeflonDir(md.metaName)
	if err := ioutil.WriteFile(md.metaName, out, 0644); err != nil {
    log.Fatalln("Failed to write meta file:", err)
  }
}

func createTeflonDir(metaName string) error {
  err := os.Mkdir(filepath.Dir(metaName), 0755)
  if err != nil {
    if os.IsExist(err){
      return nil
    }
    log.Fatalln(err)
  }
  return nil
}
