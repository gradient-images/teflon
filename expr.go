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

//go:generate pigeon -o expr_parser.go expr_parser.peg

package teflon

import "errors"

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
	GenerateAll([]string) []string
	SetNext(*ONode)
}

// String addressable version of the meta hierarchy.
type Context struct {
	IMap map[string]interface{}
	Dir  *TeflonObject
}

// Expr is an object representing a Teflon expression.
type Expr struct {
	text           string
	MetaSelector   ENode
	ObjectSelector ONode
}

// Creates a new expression object from a string
func NewExpr(text string) (*Expr, error) {
	ei, err := Parse("", []byte(text))
	if err != nil {
		return nil, err
	}
	ex := ei.(*Expr)
	ex.text = text
	return ex, nil
}

// Evaluation starts with the object selector, since it provides the context for
// the meta selection.
func (ex *Expr) Eval(c *Context) (res interface{}, err error) {
	if ex.ObjectSelector == nil {
		if ex.MetaSelector == nil {
			return c.Dir.IMap(), nil
		} else {
			cc := &Context{Dir: c.Dir, IMap: c.Dir.IMap()}
			return ex.MetaSelector.Eval(cc)
		}
	} else {
		rsl := []interface{}{}
		for {
			o := ex.ObjectSelector.NextMatch(c.Dir)
			if o == nil {
				break
			}

			if ex.MetaSelector == nil {
				rsl = append(rsl, o.Path)
			} else {
				cc := &Context{Dir: o, IMap: o.IMap()}
				m, err := ex.MetaSelector.Eval(cc)
				if err != nil {
					return nil, err
				}
				rsl = append(rsl, m)
			}
		}
		return rsl, nil
	}
}

// Generation is the process of generating a []string from an object selector.
func (ex *Expr) Generate(c *Context) (res []string, err error) {
	if ex.MetaSelector != nil {
		return nil, errors.New("Meta selector is not allowed in generator expressions.")
	}
	res = ex.ObjectSelector.GenerateAll([]string{c.Dir.Path})
	return res, nil
}

func (ex *Expr) String() string {
	return ex.text
}
