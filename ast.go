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

import (
	"errors"
	"strconv"
	"strings"
)

// String addressable version of the meta hierarchy.
type Context struct {
	IMap map[string]interface{}
}

// Expr is an object representing an expression.
type Expr struct {
	text           string
	MetaSelector   ENode
	ObjectSelector ENode
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

func (e *Expr) Eval(c *Context) (interface{}, error) {
	return e.MetaSelector.Eval(c)
}

func (e *Expr) String() string {
	return e.text
}

// ENode is the building block of the AST. The meta selector and the object
// selector both implemeted as a chain of ENodes, only the evaluation is different.
// The meta selector part is evaluated bottom up as traditional C-like expressions
// do, while the object selector part is evaluated top down as traditional globbing
// does.
type ENode interface {
	Eval(*Context) (interface{}, error)
}

// MetaSelector Nodes

// N represents a number literal
type NumberNode struct {
	Value float64
}

// S represents a string literal
type StringNode struct {
	Value string
}

// Identifier is a slice of strings that are used to search for the name in the
// Context.
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

// ObjectSelector nodes

type ShowSelector struct {
	next  *ENode
	count int
}

type ObjectName struct {
	next *ENode
	name string
}

// Eval Definitions

// MetaSelector Evals

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

// ObjectSelector Evals

func (rs *ShowSelector) Eval(c *Context) (interface{}, error) {
	return nil, nil
}

// Number needs a string conversion for string concatenation
func (N NumberNode) String() string {
	return strconv.FormatFloat(N.Value, 'G', -1, 64)
}
