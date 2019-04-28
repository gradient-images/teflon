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

// The teflon package implements all the core system functionalities of Teflon,
// currently the memory and disk representation of teflon objects and prototyping.
//
// Glossary
//
// "teflon directory": Configuration directory for the `teflon` command.
//
// "meta directory": Directory containing metadata for the dir and its content. Always
// called `.teflon`.
//
// "show": A self containing administrative structure.
//
// "target": A relative or absolute filename in the file-system.
package teflon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	protobuf "github.com/golang/protobuf/proto"
)

// Full path to the Teflon directory, which stores configuration and show prototypes
// for the system. The teflon command sets its content from the $TEFLONDIR environment
// variable.
var TeflonDir string

const (
	ShowProtoDirName = "show_proto"
	ShowPrefix       = "file://"
	protoDirName     = "proto"
	metaDirName      = ".teflon"
	metaDirMetaName  = "_"
	metaExtension    = "._"
)

// TObject is the main type of teflon. All Teflon objects are represented by this
// struct in RAM.
type TObject struct {
	Show     string
	Path     string
	Parent   *TObject
	Children []*TObject
	FileInfo FileInfo
	PersistentMeta
}

// Helper struct type
type FileInfo struct {
	os.FileInfo
}

func (f FileInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Name":    f.Name(),
		"Size":    f.Size(),
		"Mode":    f.Mode(),
		"ModTime": f.ModTime(),
		"IsDir":   f.IsDir(),
	})
}

// Returns an unitialized TObject with only the full path set to the input is set.
func NewObject(path string) (*TObject, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &TObject{Path: path}, nil
}

// Returns a new initialized TObject.
func NewInitObject(path string) (*TObject, error) {
	o, err := NewObject(path)
	if err != nil {
		return nil, err
	}
	err = o.InitMeta()
	if err != nil {
		return nil, err
	}
	return o, nil
}

// InitMeta() initializes metadata of a TObject.
func (o *TObject) InitMeta() error {
	stat, err := os.Stat(o.Path)
	if err != nil {
		return err
	}
	o.FileInfo = FileInfo{stat}

	if _, err := os.Stat(o.MetaFile()); !os.IsNotExist(err) {
		in, err := ioutil.ReadFile(o.MetaFile())
		if err != nil {
			return err
		}
		err = protobuf.Unmarshal(in, o)
		if err != nil {
			return err
		}
	}
	if o.UserData == nil {
		o.UserData = make(map[string]string)
	}

	return nil
}

// MetaFile() returns the file path to the TObject's meta file. In the case of a
// file it is:
//   $DIR/.teflon/$FILE._
// In the case of a directory it is:
//   $DIR/.teflon/_
func (o *TObject) MetaFile() string {
	if o.FileInfo.IsDir() {
		return filepath.Join(o.Path, metaDirName, metaDirMetaName)
	}
	d, n := filepath.Split(o.Path)
	return filepath.Join(d, metaDirName, n+metaExtension)
}

func (o *TObject) SetMeta(key, value string) {
	o.UserData[key] = value
}

func (o *TObject) DelMeta(key string) {
	delete(o.UserData, key)
}

// SincMeta() writes metadata to disk.
func (o *TObject) SyncMeta() error {
	out, err := protobuf.Marshal(o)
	if err != nil {
		return err
	}

	o.createTeflonDir()

	err = ioutil.WriteFile(o.MetaFile(), out, 0644)
	if err != nil {
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

func IsShow(target string) bool {
	o, err := NewInitObject(target)
	if err != nil {
		return false
	}
	if strings.HasPrefix(o.Proto, ShowPrefix) {
		return true
	}
	return false
}

func FindProtoDirs(target string) []string {
	pdl := []string{}
	target = filepath.Clean(target)
	for {
		log.Println("DEBUG: Checking for proto:", target)
		d := filepath.Join(target, metaDirName, protoDirName)
		if IsDir(d) {
			pdl = append(pdl, d)
		}
		if IsShow(target) {
			log.Println("DEBUG: Reached Show root:", target)
			break
		}
		p := filepath.Dir(target)
		log.Println("DEBUG: Moving on to parent:", p)
		if p == target {
			break
		}
		target = p
	}
	return pdl
}

// FindShowRoot finds the show root of the given target.
func FindShowRoot(target string) string {
	if IsShow(target) {
		return target
	}

	parent := filepath.Dir(target)
	if parent == target {
		return ""
	}

	return FindShowRoot(parent)
}
