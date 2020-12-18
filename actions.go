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
	"log"
	"os"
	"path/filepath"

	"github.com/gradient-images/teflon/internal/meta"

	"github.com/otiai10/copy"
)

// Get() evaluates a Teflon expression and returns the result.
func (o *TeflonObject) Get(exs string) (res interface{}, err error) {
	log.Printf("DEBUG: Inside Get(): o.Path: %v  ex: %v", o.Path, exs)

	ex, err := NewExpr(exs)
	if err != nil {
		return nil, err
	}

	c := &Context{Dir: o}
	res, err = ex.Eval(c)

	return res, nil
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

		if o.SyncMeta() != nil {
			log.Fatalln("ABORT: Couldn't write meta of newly created show:", err)
		}

		oSl = append(oSl, o)
	}
	return oSl, nil
}

// CreateObject() creates a new FS object and triggers a new event.
func (o *TeflonObject) CreateObject(exs string, file bool) (oSl []*TeflonObject, err error) {
	log.Printf("DEBUG: Inside CreateObject(): o.Path: %v  exs: %v", o.Path, exs)

	ex, err := NewExpr(exs)
	if err != nil {
		return nil, err
	}

	c := &Context{Dir: o}

	res, err := ex.Generate(c)
	if err != nil {
		return nil, err
	}

	for _, fsp := range res {
		if _, err := os.Stat(fsp); !os.IsNotExist(err) {
			log.Println("WARNING: Target already exists:", fsp)
			continue
		}

		rch := make(chan *TeflonObject)
		Events <- Event{o, PreNew, rch}
		<-rch

		if file {
			f, err := os.Create(fsp)
			if err != nil {
				log.Println("WARNING: Couldn't create file:", fsp, err)
				continue
			}
			f.Close()
		} else {
			err := os.Mkdir(fsp, 0755)
			if err != nil {
				log.Println("WARNING: Couldn't create directory:", fsp, err)
				continue
			}
		}

		o, err := NewTeflonObject(fsp)
		if err != nil {
			return nil, err
		}

		oSl = append(oSl, o)

		Events <- Event{o, PostNew, nil}

		log.Println("SUCCESS: Created:", fsp)
	}
	return oSl, nil
}

func (o *TeflonObject) SetContractPattern(exs string, pat string) (oSl []*TeflonObject, err error) {
	ex, err := NewExpr(exs)
	if err != nil {
		return nil, err
	}

	c := &Context{Dir: o}

	res, err := ex.Generate(c)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("Pattern returned nothing:" + exs)
	}

	for _, fsp := range res {
		o, err := NewTeflonObject(fsp)
		if err != nil {
			return nil, err
		}
		o.Contract = &meta.Contract{Pattern: pat}
		err = o.SyncMeta()
		if err != nil {
			return oSl, err
		}
		oSl = append(oSl, o)
	}
	return oSl, nil
}
