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
	e := ei.(*Expr)
	e.text = text
	return e, nil
}

// Evaluation starts with the object selector, since it provides the context for
// the meta selection.
func (e *Expr) Eval(c *Context) (res interface{}, err error) {
	if e.ObjectSelector == nil {
		if e.MetaSelector == nil {
			return c.Dir, nil
		} else {
			cc := &Context{Dir: c.Dir, IMap: c.Dir.IMap()}
			return e.MetaSelector.Eval(cc)
		}
	} else {
		rsl := []interface{}{}
		for {
			o := e.ObjectSelector.NextMatch(c.Dir)
			if o == nil {
				break
			}

			if e.MetaSelector == nil {
				rsl = append(rsl, o.Path)
			} else {
				cc := &Context{Dir: o, IMap: o.IMap()}
				m, err := e.MetaSelector.Eval(cc)
				if err != nil {
					return nil, err
				}
				rsl = append(rsl, m)
			}
		}
		return rsl, nil
	}
}

func (e *Expr) String() string {
	return e.text
}
