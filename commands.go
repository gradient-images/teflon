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
func (o *TeflonObject) CreateShow(exs string) (no *TeflonObject, err error) {
	log.Printf("DEBUG: Inside CreateShow(): o.Path: %v  exs: %v", o.Path, exs)
	ex, err := NewExpr(exs)
	if err != nil {
		return nil, err
	}

	c := &Context{Dir: o}

	// NOTE: Eval's first return value needs clarification in the case of an error.
	res, err := ex.Generate(c)
	if err != nil {
		return nil, err
	}

	// Create display string of result (dres).
	dres, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatalln("FATAL: Couldnt marshal result JSON:", err)
	}
	log.Printf("DEBUG: dres: %s\n", dres)

	switch l := len(res); {
	case l == 0:
		return nil, errors.New("Pattern returned nothing:" + exs)
	case l > 1:
		return nil, errors.New("More than one value returned:" + exs)
	}
	return nil, err
}

// CreateObject() creates a new FS object and triggers a new event.
func CreateObject() {

}
