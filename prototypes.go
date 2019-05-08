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
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ListProtos lists all prototypes seen in the context of the object.
func (o *TeflonObject) ListProtos() (map[string]string, error) {
	if o.Show == nil {
		return nil, errors.New("Prototyping not supported outside shows.")
	}
	var protoMap = map[string]string{}
	if !o.FileInfo.IsDir {
		o = o.Parent
	}
	for {
		d := filepath.Join(o.Path, teflonDirName, protoDirName)
		protoList, err := ioutil.ReadDir(d)
		if err != nil {
			if os.IsNotExist(err) {
				o = o.Parent
				continue
			}
			return nil, err
		}
		for _, p := range protoList {
			if _, ok := protoMap[p.Name()]; !ok {
				saProto, _ := ShowAbs(filepath.Join(d, p.Name()))
				protoMap[p.Name()] = saProto
			}
		}
		if o.ShowRoot {
			return protoMap, nil
		}
		o = o.Parent
	}
}

// FindProto finds prototype by it's exact name in the context of object 'o'.
func (p *TeflonObject) FindProto(proto string) (string, error) {
	for {
		// Create candidate.
		c := filepath.Join(p.Path, teflonDirName, protoDirName, proto)
		// If proto exists.
		if Exist(c) {
			return c, nil
		}
		// If reached show root.
		if p.ShowRoot {
			return "", errors.New("Couldn't find proto: " + proto)
		}
		p = p.Parent
	}
}

// Find prototype for a given target in the context of object 'o'. There are
// three kinds of prototypes Teflon looks for: 'full', 'prefix' and 'extension'.
// The function traverses the show hierarchy upwards looking for these returning
// the first matching prototype it founds.
func (o *TeflonObject) FindProtoForTarget(name string) (string, error) {
	// Determine extension prototype name.
	ext := filepath.Ext(name)
	if ext != "" {
		ext = "_" + ext[1:]
	}

	// Determine prefix prototype name.
	pre := strings.SplitN(name, "_", 2)[0]
	if pre == name {
		pre = ""
	} else {
		pre = pre + "_"
	}

	// Traverse directory tree upwards until show root.
	for {
		d := filepath.Join(o.Path, teflonDirName, protoDirName)
		if c := filepath.Join(d, name); Exist(c) {
			return c, nil
		}
		if c := filepath.Join(d, pre); pre != "" && Exist(c) {
			return c, nil
		}
		if c := filepath.Join(d, ext); ext != "" && Exist(c) {
			return c, nil
		}
		// If reached show root.
		if o.ShowRoot {
			return "", errors.New("No appropriate proto for: " + name)
		}
		o = o.Parent
	}
}

// Sets prototype of the object 'o' to the given proto. It also sets the
// prototype's instance list to include the object.
func (o *TeflonObject) SetProto(proto string) error {
	// Convert proto and object path to show-absolute notation.
	saProto, err := ShowAbs(proto)
	if err != nil {
		return err
	}
	saTarget, err := ShowAbs(o.Path)
	if err != nil {
		return err
	}

	// Remove old proto from old proto's instance list.
	oldProto := o.Proto
	if o.Proto != "" {
		op, err := NewTeflonObject(oldProto)
		if err != nil {
			return err
		}
		for i := range op.Instances {
			if op.Instances[i] == saTarget {
				op.Instances = append(op.Instances[:i], op.Instances[i+1:]...)
				if op.SyncMeta() != nil {
					return err
				}
				break
			}
		}
	}

	// Set new prototype for the object.
	o.Proto = saProto
	o.Instances = []string{}
	if o.SyncMeta() != nil {
		return err
	}

	// Set instances field on proto object.
	p, err := NewTeflonObject(proto)
	if err != nil {
		log.Fatalln("ABORT: Couldn't create object:", err)
	}
	p.Instances = append(p.Instances, saTarget)

	if p.SyncMeta() != nil {
		return err
	}
	return nil
}
