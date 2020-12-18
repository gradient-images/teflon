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

// Commented out for faster compilation. Uncomment it to update.
// //go:generate stringer -type EventType

package teflon

import (
  "log"
  "path/filepath"
  )

// EventType is an enum describing the type of a Teflon event.
type EventType int

const (
	PreNew EventType = iota
	PostNew
)

type Event struct {
	Object *TeflonObject
	Type   EventType
  Result chan *TeflonObject
}

var Events chan Event = make(chan Event)
var Done chan bool = make(chan bool)

func listen() {
  for evt := range Events {
    log.Printf("EVENT: Event received: %s (%s)", evt.Object.GetPath(), evt.Type)

    p := evt.Object
    for {
      p = p.Parent
      cfs := filepath.Join(p.Path, contractDirName)
      var c *TeflonObject
      var err error
      if Exist(cfs) {
        c, err = NewTeflonObject(cfs)
        if err != nil {
          log.Fatalln("ABORT: Can't create contract dir object:", err)
        }
        log.Println("DEBUG: Created contract dir object:", c.Path)
      }
      if p.ShowRoot {
        break
      }
    }

    if evt.Result != nil {
      evt.Result <- evt.Object
    }
  }
  Done <- true
}
