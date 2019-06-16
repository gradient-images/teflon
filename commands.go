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
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

// Evaluates a Teflon expression and returns the result
func Get(dirs string, exs string) (res interface{}, err error) {
	log.Printf("DEBUG: Inside Get(): dir: %v  ex: %v", dirs, exs)

	ex, err := NewExpr(exs)
	if err != nil {
		return nil, err
	}

	dir, err := NewTeflonObject(dirs)
	if err != nil {
		return nil, err
	}

	c := &Context{Dir: dir}
	res, err = ex.Eval(c)

	return
}

// CreateShow() creates new Teflon show.
func (o *TeflonObject) CreateShow(exs string, protoName string) (oSl []*TeflonObject, err error) {
	log.Printf("DEBUG: Inside CreateShow(): o.Path: %v  exs: %v", o.Path, exs)

	ex, err := NewExpr(exs)
	if err != nil {
		return nil, err
	}

	c := &Context{Dir: o}

	res, err := ex.Generate(c)
	if err != nil {
		return nil, err
	}

	// Create display string of result (dres).
	dres, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatalln("ABORT: Couldnt marshal result JSON:", err)
	}
	log.Printf("DEBUG: dres: %s\n", dres)

	if len(res) == 0 {
		return nil, errors.New("Pattern returned nothing:" + exs)
	}

	proto := filepath.Join(TeflonConf, ShowProtoDirName, protoName)

	for _, fsp := range res {
		if _, err := os.Stat(fsp); !os.IsNotExist(err) {
			log.Println("WARNING: Target already exists:", fsp)
			continue
		}

		err = copy.Copy(proto, fsp)
		if err != nil {
			log.Fatalln("ABORT: Couldn't copy show proto:", err)
		}
		log.Printf("SUCCESS: Created new show: %s (%s)", fsp, protoName)

		o, err := NewTeflonObject(fsp)
		if err != nil {
			log.Fatalln("ABORT: Couldn't create object:", err)
		}

		o.ShowRoot = true
		o.Proto = proto

		if o.SyncMeta() != nil {
			log.Fatalln("ABORT: Couldn't write meta of newly created show:", err)
		}

		oSl = append(oSl, o)
	}
	return oSl, nil
}

// CreateObject() creates a new FS object and triggers a new event.
func CreateObject() {

}
