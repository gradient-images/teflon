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

package expr

import (
  "strings"
  "strconv"
  "errors"
  "log"

  "github.com/gradient-images/teflon/meta"
)


// String addressable version of the meta hierarchy.
type Context struct {
  IMap map[string]interface{}
}

// Expr is an object representing an expression.
type Expr struct {
  text string
  ast ENode
}

// Creates a new expression object from a string
func New(t string) *Expr {
  return &Expr{text: t}
}

func (e *Expr) Parse() error {
  ast, err := Parse(e.text, []byte(e.text))
  if err != nil {
    return err
  }
  e.ast = ast.(ENode)
  return nil
}

func (e *Expr) Eval(c *Context) (*meta.UserValue, error) {
  return e.ast.Eval(c)
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
  Eval(*Context) (*meta.UserValue, error)
}


// A Teflon expression is composed of a meta selector and an object selector
type ExprNode struct {
  MetaSelector ENode
  // ObjectSelector ENode
}

// N represents a number literal
type NumberNode struct {
  Value meta.UserValue_N
}

// S represents a string literal
type StringNode struct {
  Value meta.UserValue_S
}

// Identifier is a slice of strings that are used to search for the name in the
// Context.
type MetaNode struct {
  NameList []string
}

// Adds numbers numberically and concatenate strings
type AddNode struct {
  first ENode
  second ENode
}

func (Expr *ExprNode) Eval(c *Context) (*meta.UserValue, error) {
  return Expr.MetaSelector.Eval(c)
}

func (N *NumberNode) Eval(c *Context) (*meta.UserValue, error) {
  return &meta.UserValue{Value: &N.Value}, nil
}

func (S *StringNode) Eval(c *Context) (*meta.UserValue, error) {
  return &meta.UserValue{Value: &S.Value}, nil
}

func (M *MetaNode) Eval(c *Context) (*meta.UserValue, error) {
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

  // Since values are coming from the context it's enogh to handle JSON UserValues
  var uv *meta.UserValue
  log.Println("DEBUG: val:", val)
  switch v := val.(type) {
  case float64:
    uv = &meta.UserValue{Value: &meta.UserValue_N{v}}
  case string:
    uv = &meta.UserValue{Value: &meta.UserValue_S{v}}
  default:
    return nil, errors.New("Couldn't convert result to UserValue.")
  }
  return uv, nil
}

func (a *AddNode) Eval(c *Context) (*meta.UserValue, error) {
  log.Printf("DEBUG: Inside AddNode.Eval.")
  fp, err := a.first.Eval(c)
  if err != nil {
    return nil, err
  }
  f := *fp

  sp, err := a.second.Eval(c)
  if err != nil {
    return nil, err
  }
  s := *sp

  var v *meta.UserValue

  switch fv := f.Value.(type) {
  case *meta.UserValue_N:
    switch sv := s.Value.(type) {
    case *meta.UserValue_N:
      v = &meta.UserValue{Value: &meta.UserValue_N{fv.N + sv.N}}
    }
  }
  return v, nil
}


// Number needs a string conversion for string concatenation
func (N NumberNode) String() string {
  return strconv.FormatFloat(N.Value.N, 'G', -1, 64)
}
