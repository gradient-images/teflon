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

package teflon

import (
  "log"
  "os"
  "path/filepath"
  "io/ioutil"

  "github.com/golang/protobuf/proto"
  "github.com/gradient-images/teflon/internal/metadata"
)

var (
  TeflonRoot string
  TeflonDir string
)

const (
  TeflonDirName = ".teflon_root"
  ShowProtoDirName = "show_proto"
  ProtoDirName = "proto"
  metaDirName = ".teflon"
	metaDirMetaName = "_"
	metaExtension = "._"
)

type TeflonError struct {
  Message  string
}

func (err TeflonError) Error() string {
  return err.Message
}

type TObject struct {
  Path, metaPath string
  meta *metadata.UserSection
}

func NewObject(path string) *TObject {
  return &TObject{Path: path}
}

func (o *TObject) InitMeta() error {
  stat, err := os.Stat(o.Path)
  if err != nil {
    return err
  }

  if stat.IsDir() {
    o.metaPath = filepath.Join(o.Path, metaDirName, metaDirMetaName)
  } else {
    d, n := filepath.Split(o.Path)
    o.metaPath = filepath.Join(d, metaDirName, n + metaExtension)
  }

  o.meta = &metadata.UserSection{}

  if _, err := os.Stat(o.metaPath); os.IsNotExist(err) {
    log.Println("Meta file doesn't exists.")
    o.meta.UserData = make(map[string]string)
  } else {
    log.Print("Meta file exists.")
    in, err := ioutil.ReadFile(o.metaPath)
    if err != nil {
      return err
    }
    err = proto.Unmarshal(in, o.meta)
    if err != nil {
      return err
    }
  }
  return nil
}

func (o TObject) GetMeta() (*metadata.UserSection, error) {
  if o.meta != nil {
    return o.meta, nil
  }

  err := o.InitMeta()
  if err != nil {
    return nil, err
  }

  return o.meta, nil
}

func (o TObject) SetMeta(key, value string) {
  o.meta.UserData[key] = value
}

func (o TObject) SyncMeta() error {
  out, err := proto.Marshal(o.meta)
	if err != nil {
		return err
	}

  o.createTeflonDir()

	if err := ioutil.WriteFile(o.metaPath, out, 0644); err != nil {
    return err
  }
  return nil
}

func (o TObject) createTeflonDir() error {
  err := os.Mkdir(filepath.Dir(o.metaPath), 0755)
  if err != nil {
    if os.IsExist(err){
      return nil
    }
    return err
  }
  return nil
}

func IsDir(target string) bool {
  fi, err := os.Stat(target)
  if err != nil { return false }
  if fi.IsDir() {return true }
  return false
}

func FindShowRoot(target string) (string, error) {
  dirName := filepath.Join(target, ShowProtoDirName)
  if IsDir(dirName) {
    return target, nil
  }

  parent := filepath.Dir(target)
  if parent == target {
    return "", TeflonError{Message: "Show not found."}
  }

  return FindShowRoot(parent)
}
