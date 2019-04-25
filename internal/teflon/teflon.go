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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var (
	TeflonRoot string
	TeflonDir  string
)

const (
	TeflonDirName    = ".teflon_root"
	ShowProtoDirName = "show_proto"
	ProtoDirName     = "proto"
	metaDirName      = ".teflon"
	metaDirMetaName  = "_"
	metaExtension    = "._"
)

type TeflonError struct {
	Message string
}

func (err TeflonError) Error() string {
	return err.Message
}

func NewObject(path string) (*TObject, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &TObject{File: &FileInfo{Name: path}}, nil
}

func (o *TObject) MetaFile() string {
	if o.File.IsDir {
		return filepath.Join(o.File.Name, metaDirName, metaDirMetaName)
	}
	d, n := filepath.Split(o.File.Name)
	return filepath.Join(d, metaDirName, n+metaExtension)
}

func (o *TObject) InitMeta() error {
	stat, err := os.Stat(o.File.Name)
	if err != nil {
		return err
	}

	// Init FileInfo. Later checking if something is changed
	// will come here.
	o.File.Size = stat.Size()
	o.File.ModTime = &timestamp.Timestamp{
		Seconds: stat.ModTime().Unix(),
		Nanos:   int32(stat.ModTime().UnixNano()),
	}
	o.File.IsDir = stat.IsDir()

	if _, err := os.Stat(o.MetaFile()); os.IsNotExist(err) {
		log.Println("Meta file doesn't exists.")
		o.UserData = make(map[string]string)
	} else {
		log.Print("Meta file exists.")
		in, err := ioutil.ReadFile(o.MetaFile())
		if err != nil {
			return err
		}
		err = proto.Unmarshal(in, o)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (o TObject) GetAllMeta() (*UserSection, error) {
// 	if o.UserData != nil {
// 		return o.UserData, nil
// 	}
//
// 	err := o.InitMeta()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return o.UserData, nil
// }

func (o *TObject) SetMeta(key, value string) {
	o.UserData[key] = value
}

func (o *TObject) DelMeta(key string) {
	delete(o.UserData, key)
}

func (o *TObject) SyncMeta() error {
	out, err := proto.Marshal(o)
	if err != nil {
		return err
	}

	o.createTeflonDir()

	if err := ioutil.WriteFile(o.MetaFile(), out, 0644); err != nil {
		return err
	}
	return nil
}

func (o TObject) createTeflonDir() error {
	err := os.Mkdir(filepath.Dir(o.MetaFile()), 0755)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func IsDir(target string) bool {
	fi, err := os.Stat(target)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
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
