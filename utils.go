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
	"os"
	"path/filepath"
	"strings"
)

// Tells if a path is a dir or not.
func IsDir(fspath string) bool {
	fi, err := os.Stat(fspath)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

// Returns true if a file-system object exists at a give path.
func Exist(fspath string) bool {
	_, err := os.Stat(fspath)
	return !os.IsNotExist(err)
}

// Converts a target to a file-sytem absolute path.
func FSPath(target string) (string, error) {

	// Checks if target is show absolute.
	if strings.HasPrefix(target, "//") {
		o, err := NewTeflonObject(".")
		if err != nil {
			return "", err
		}
		if o.Show == nil {
			return "", errors.New("Couldn't resolve '//'.")
		}
		return filepath.Join(o.Show.Path, strings.TrimPrefix(target, "/")), nil
	}

	// Checks if target is file-system absolute.
	if strings.HasPrefix(target, "/") {
		return filepath.Clean(target), nil
	}

	// If neither of the above then it's relative.
	fspath, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	return fspath, nil
}

// Converts file-system absolute path to show-absolute notation. If the conversion
// can not be made, the function reurns the original string.
func ShowAbs(fspath string) string {
	o, err := NewTeflonObject(fspath)
	if err != nil {
		return fspath
	}
	return strings.Replace(fspath, o.Show.Path, "/", 1)
}
