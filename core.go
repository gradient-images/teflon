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

// The teflon package implements all the system functionalities of Teflon.
//
// Unfortunately this document is polluted with struct fields and methods generated
// by `protoc` during compilation. Method names starting with `XXX` and `Get`
// belong to this category, please ignore them, as we don't use them and we are not
// able to document them by `godoc`. Real Teflon getters' names are identical to
// the field they are getting, without the 'Get' prefix, so look for them for
// information. Also note that the source is not polluted.
//
// Glossary
//
// "configuration directory": Configuration directory for the `teflon` command.
// Teflon stores the show prototypes here.
//
// "teflon directory": Directory containing Teflon related information for the
// current directory, like metadata and prototypes.
//
// "show": A self containing administrative structure.
//
// "target": A target is a filesystem path in Teflon's notation. The only
// distinction is that if the target string starts with '//', that means that it
// is "show-absolute", so the "//" points to the show root of the current
// directory.
package teflon

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// Full path to the configuration directory, which stores configuration and show
// prototypes for the system. The teflon command gets its value from the
// $TEFLONCONF environment variable.
var TeflonConf string

// Objects associates teflon object pointers to absolute file-system paths.
var Objects = map[string]*TeflonObject{}

// Shows are a list of all the show obejcts that the process knows about.
var Shows = []*TeflonObject{}

const (
	ShowProtoDirName = "show_proto"
	protoDirName     = "proto"
	teflonDirName    = ".teflon"
	metaDirMetaName  = "_"
	metaExtension    = "._"
)

// TeflonObject is the main type of teflon. All Teflon objects are represented by
// this struct in RAM.
type TeflonObject struct {
	Path     string
	Show     *TeflonObject
	FileInfo FileInfo
	Parent   *TeflonObject
	PersistentMeta
}

// Marshaling JSON manually to avoid recursion. There is probably a more elegant
// way of doing this.
func (o TeflonObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Show":      o.Show.Path,
		"Path":      o.Path,
		"Parent":    o.Parent.Path,
		"FileInfo":  o.FileInfo,
		"ShowRoot":  o.ShowRoot,
		"Proto":     o.Proto,
		"Instances": o.Instances,
		"UserData":  o.UserData,
	})
}

// NewTeflonObject creates a new initialized Teflon object.
//
// Initialization is always complete. It is not allowed to have half-baked objects
// in memory. Since it has to set the Show field to its correct value it first has
// to find the show root of the target. This is done recursively by creating the
// parent object. This means that not only the created object is fully initialized
// but there will be a complete chain of object leading from the target to the show
// root.
//
// If the target is show-absolute then the system first has to find the show
// root from the current directory to get the file-system path of the target. This
// means that the end result is two initialized chains to the same show root.
func NewTeflonObject(target string) (*TeflonObject, error) {

	// Convert target to file-system path
	fspath, err := FSPath(target)
	if err != nil {
		return nil, err
	}

	// Checks if it's in Objects
	o, ok := Objects[fspath]
	if ok {
		return o, nil
	} else {
		// Create the uninitialized object
		o = &TeflonObject{Path: fspath}
	}

	// Initialize metadata
	stat, err := os.Stat(o.Path)
	if err != nil {
		return nil, err
	}
	modtime, _ := ptypes.TimestampProto(stat.ModTime())
	o.FileInfo = FileInfo{
		Name:    stat.Name(),
		Size:    stat.Size(),
		Mode:    uint32(stat.Mode()),
		ModTime: modtime,
		IsDir:   stat.IsDir(),
	}
	m := o.MetaFile()

	// Read meta file if exists
	if _, err := os.Stat(m); !os.IsNotExist(err) {
		in, err := ioutil.ReadFile(m)
		if err != nil {
			return nil, err
		}
		err = protobuf.Unmarshal(in, o)
		if err != nil {
			return nil, err
		}
	}

	// Init UserData if not exists
	if o.UserData == nil {
		o.UserData = make(map[string]string)
	}

	// Check if it is show root
	if o.ShowRoot {
		o.Show = o
	} else {
		parent := filepath.Dir(o.Path)

		// Check if reached file-system root
		if parent != "/" {
			p, err := NewTeflonObject(parent)
			if err != nil {
				return nil, err
			}
			o.Parent = p
			o.Show = p.Show
		}
	}

	Objects[fspath] = o
	return o, nil
}

// Converts a target to a file-sytem absolute path.
func FSPath(target string) (string, error) {

	// Checks if target is show absolute.
	if strings.HasPrefix(target, "//") {
		o, err := NewTeflonObject(".")
		if err != nil {
			return "", err
		}
		if o.Show == nil {
			return "", errors.New("Couldn't resolve '//'.")
		}
		return filepath.Join(o.Show.Path, strings.TrimPrefix(target, "/")), nil
	}

	// Checks if target is file-system absolute.
	if strings.HasPrefix(target, "/") {
		return filepath.Clean(target), nil
	}

	// If neither of the above then it's relative.
	fspath, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	return fspath, nil
}

// MetaFile returns the file path to the TeflonObject's meta file. In the case of a
// file it is:
//   $DIR/.teflon/$FILE._
// In the case of a directory it is:
//   $DIR/.teflon/_
func (o *TeflonObject) MetaFile() string {
	if o.FileInfo.IsDir {
		return filepath.Join(o.Path, teflonDirName, metaDirMetaName)
	}
	d, n := filepath.Split(o.Path)
	return filepath.Join(d, teflonDirName, n+metaExtension)
}

// Sets an entry in the user section of the metadata.
func (o *TeflonObject) SetMeta(key, value string) {
	o.UserData[key] = value
}

// Deletes an entry from the user section of the metadata.
func (o *TeflonObject) DelMeta(key string) {
	delete(o.UserData, key)
}

// SincMeta() writes metadata to disk.
func (o *TeflonObject) SyncMeta() error {
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

// Creates the Teflon directory for the object's meta file.
func (o TeflonObject) createTeflonDir() error {
	err := os.Mkdir(filepath.Dir(o.MetaFile()), 0755)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

// Tells if a path is a dir or not.
func IsDir(fspath string) bool {
	fi, err := os.Stat(fspath)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

func IsShow(target string) bool {
	o, err := NewTeflonObject(target)
	if err != nil {
		return false
	}
	return o.ShowRoot
}

// Creates a list of objects containing proto directories.
func FindProtoDirs(target string) []string {
	pdl := []string{}
	target = filepath.Clean(target)
	for {
		d := filepath.Join(target, teflonDirName, protoDirName)
		if IsDir(d) {
			pdl = append(pdl, d)
		}
		if IsShow(target) {
			break
		}
		p := filepath.Dir(target)
		if p == target {
			break
		}
		target = p
	}
	return pdl
}
