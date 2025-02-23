// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//nolint
package archive

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

const basePath = "/tmp"

const fileContent = "content"

type fileDef struct {
	name    string
	content string
}

type dirDef struct {
	path   string
	files  []fileDef
	subDir []dirDef
}

var filesToCreate = []dirDef{
	{
		path: "testDirA",
		files: []fileDef{
			{
				name:    "testFileA",
				content: fileContent,
			},
			{
				name:    "testFileB",
				content: fileContent,
			},
		},
		subDir: []dirDef{
			{
				path: "testDirC",
				files: []fileDef{
					{
						name:    "testFileA",
						content: fileContent,
					},
					{
						name:    "testFileB",
						content: fileContent,
					},
				},
			},
		},
	},
	{
		path: "testDirB",
		files: []fileDef{
			{
				name:    "testFileA",
				content: fileContent,
			},
			{
				name:    "testFileB",
				content: fileContent,
			},
		},
	},
}

func makeDir(root string, d dirDef) error {
	currentDir := filepath.Join(root, d.path)
	err := os.MkdirAll(currentDir, 0755)
	if err != nil {
		return err
	}

	for _, file := range d.files {
		_, err = os.Create(filepath.Join(currentDir, file.name))
		if err != nil {
			return err
		}
	}

	for _, sub := range d.subDir {
		err = makeDir(currentDir, sub)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestTarWithoutRootDir(t *testing.T) {
}

func TestTarWithRootDir(t *testing.T) {
	reader, err := TarWithRootDir("/Users/eric/Workspace/src/github.com/vbauerster/mpb", "/Users/eric/Workspace/src/github.com/vbauerster/empty")
	if err != nil {
		t.Error(err)
	}

	tmp, err := os.CreateTemp("/tmp", "tar")
	_, err = io.Copy(tmp, reader)
	if err != nil {
		t.Error(err)
	}
}
