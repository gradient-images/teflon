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

import "log"

type EventType int

const (
	PreNew EventType = iota
	PostNew
)

type Event struct {
	Object *TeflonObject
	Type   EventType
}

var events chan Event = make(chan Event)

func listen() {
	log.Println("DEBUG: Inside listen().")
  for {
    e := <- events
    log.Printf("DEBUG: Event received: %s (%d)", e.Object.GetPath(), e.Type)
  }
}
