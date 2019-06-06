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
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ENode is the building block of the AST. The meta selector and the object
// selector both implemeted as a chain of ENodes, only the evaluation is different.
// The meta selector part is evaluated bottom up as traditional C-like expressions
// do, while the object selector part is evaluated top down as traditional globbing
// does.
type ENode interface {
	Eval(*Context) (interface{}, error)
}

// ONode must be implemented by object selector nodes.
type ONode interface {
	// Match(string) bool
	NextMatch(*TeflonObject) *TeflonObject
	// NextGen(*TeflonObject) string
	SetNext(*ONode)
}

// MetaSelector Nodes

// AllMetaNode returns all metadata of an object
type AllMetaNode struct{}

// NumberNode represents a number literal
type NumberNode struct {
	Value float64
}

// StringNode represents a string literal
type StringNode struct {
	Value string
}

// MetaNode represents a metadata identifier
type MetaNode struct {
	NameList []string
}

// Adds numbers numberically and concatenate strings
type AddNode struct {
	first  ENode
	second ENode
}

// Adds numbers numberically and concatenate strings
type SubNode struct {
	first  ENode
	second ENode
}

// Adds numbers numberically and concatenate strings
type MulNode struct {
	first  ENode
	second ENode
}

// Adds numbers numberically and concatenate strings
type DivNode struct {
	first  ENode
	second ENode
}

//
// ObjectSelector nodes
//

type AbsPath struct {
	next   *ONode
	noMore bool
	count  int
}

type RelPath struct {
	next   *ONode
	noMore bool
	count  int
}

type ObjectName struct {
	next      *ONode
	noMore    bool
	name      string
	multi     bool
	pattern   *regexp.Regexp
	index     int
	lastMatch *TeflonObject
}

var multiPatt *regexp.Regexp = regexp.MustCompile(`.*[*].*`)

func NewObjectName(name string) (onn *ObjectName) {
	onn = &ObjectName{name: name}
	if multiPatt.MatchString(name) {
		patts := strings.ReplaceAll(name, "*", ".*")
		log.Println("DEBUG: patts:", patts)
		onn.pattern = regexp.MustCompile(patts)
		onn.multi = true
	}
	return
}

//
// ENode Implementations
//

func (amn *AllMetaNode) Eval(c *Context) (interface{}, error) {
	return c.IMap, nil
}

func (N *NumberNode) Eval(c *Context) (interface{}, error) {
	return N.Value, nil
}

func (S *StringNode) Eval(c *Context) (interface{}, error) {
	return S.Value, nil
}

func (M *MetaNode) Eval(c *Context) (interface{}, error) {
	var val interface{}
	v := c.IMap
	for i, n := range M.NameList {
		// Create lower map for case insensitive matching
		lm := map[string]string{}
		for k := range v {
			lm[strings.ToLower(k)] = k
		}

		var ok bool
		val, ok = v[lm[strings.ToLower(n)]]
		if !ok {
			return nil, errors.New("Couldn't find key in meta: " + n)
		}

		// If there is more name to come
		if i < len(M.NameList)-1 {
			switch val.(type) {
			case map[string]interface{}:
				// Convert next level to map
				v = val.(map[string]interface{})
			default:
				return nil, errors.New("Couldn't find key in meta: " + n)
			}
		}
	}
	return val, nil
}

func (a *AddNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f + s
		}
	}
	return v, nil
}

func (a *SubNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f - s
		}
	}
	return v, nil
}

func (a *MulNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f * s
		}
	}
	return v, nil
}

func (a *DivNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f / s
		}
	}
	return v, nil
}

//
// ONode Implementations
//

func (apn *AbsPath) NextMatch(o *TeflonObject) (res *TeflonObject) {
	if apn.noMore {
		return nil
	}

	var err error
	if apn.count == 1 {
		res, err = NewTeflonObject("/")
		if err != nil {
			return nil
		}
	} else {
		res, err = NewTeflonObject("//")
		if err != nil {
			return nil
		}
	}

	if apn.next != nil {
		res = (*apn.next).NextMatch(res)
	}

	if res == nil {
		apn.noMore = true
	}

	return res
}

func (apn *AbsPath) SetNext(node *ONode) {
	apn.next = node
}

func (rpn *RelPath) NextMatch(o *TeflonObject) (res *TeflonObject) {
	var err error
	if rpn.noMore {
		return nil
	}

	log.Printf("DEBUG: res: (%v, %T)", res, res)
	log.Printf("DEBUG: rpn.count: %v", rpn.count)

	// Give back o or traverse upvards.
	if rpn.count > 1 {
		fspath := o.Path
		for i := 1; i < rpn.count; i++ {
			fspath = filepath.Dir(fspath)
		}
		res, err = NewTeflonObject(fspath)
		if err != nil {
			log.Fatalln("FATAL: Couldn't create object:", fspath)
		}
	} else {
		res = o
	}

	if rpn.next != nil {
		log.Println("DEBUG: Traversing next from RelPath.")
		res = (*rpn.next).NextMatch(res)
	} else {
		rpn.noMore = true
	}

	if res == nil {
		rpn.noMore = true
	}
	return res
}

func (rpn *RelPath) SetNext(node *ONode) {
	rpn.next = node
}

func (onn *ObjectName) NextMatch(o *TeflonObject) (res *TeflonObject) {
	var err error
	if onn.noMore {
		return nil
	}

	if onn.multi {
		if onn.lastMatch == nil {
			children, err := filepath.Glob(filepath.Join(o.Path, "*"))
			if err != nil {
				onn.noMore = true
				return nil
			}
			for onn.index < len(children) {
				name := filepath.Base(children[onn.index])
				if onn.pattern.MatchString(name) {
					res, err = NewTeflonObject(children[onn.index])
					onn.index++
					if err != nil {
						log.Fatalln("FATAL: Couldn't create found object.", name)
					}
					break
				}
				onn.index++
			}
			if res == nil {
				onn.noMore = true
				return nil
			} else {
				onn.lastMatch = res
			}
			log.Println("DEBUG: Found match:", res.Path, onn.index)
		} else {
			res = onn.lastMatch
		}
	} else {
		// If not multi.
		res, err = NewTeflonObject(filepath.Join(o.Path, onn.name))
		if err != nil {
			onn.noMore = true
			return nil
		}
	}

	if onn.next != nil {
		res = (*onn.next).NextMatch(res)
		if res == nil {
			onn.lastMatch = nil
		}
	} else {
		onn.lastMatch = nil
	}

	if onn.next == nil && !onn.multi {
		onn.noMore = true
	}

	if res == nil {
		onn.lastMatch = nil
		onn.noMore = true
	}

	// log.Printf("DEBUG: Returning res.Path: %v", res.Path)

	return res
}

func (onn *ObjectName) SetNext(node *ONode) {
	onn.next = node
}

//
// Utility Functions
//

// NumberNode needs to be Stringer for string concatenation
func (N NumberNode) String() string {
	return strconv.FormatFloat(N.Value, 'G', -1, 64)
}
