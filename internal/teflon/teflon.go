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
  // "log"
  "os"
  "path/filepath"

  "github.com/gradient-images/teflon/internal/metadata"
)

var (
  TeflonRoot string
  TeflonDir string
)

const (
  TeflonDirName = ".teflon_root"
  ShowProtoDirName = "show_proto"
  ShowDirName = "show"
  ProtoDirName = "asset/proto"
)

type TObject interface {
  Name() string
  Meta() metadata.UserSection
}

type TeflonError struct {
  Message  string
}

func (err TeflonError) Error() string {
  return err.Message
}

func IsDir(target string) bool {
  fi, err := os.Stat(target)
  if err != nil { return false }
  if fi.IsDir() {return true }
  return false
}

func FindShowRoot(target string) (string, error) {
  dirName := filepath.Join(target, ShowDirName)
  if IsDir(dirName) {
    return target, nil
  }

  parent := filepath.Dir(target)
  if parent == target {
    return "", TeflonError{Message: "Show not found."}
  }

  return FindShowRoot(parent)
}
