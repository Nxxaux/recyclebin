// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !windows

package recyclebin

import (
	"path"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

func getTrashedFilename(trashPath string, filename string) string {
	extension := path.Ext(filename)
	bareName := strings.TrimSuffix(filename, extension)
	newFilename := filename
	index := -1
	isDuplicateFilename := true
	for isDuplicateFilename {
		index += 1
		newFilename = bareName + strconv.Itoa(index) + extension //Trash file
		existsTrashFile, _ := afero.Exists(fs, buildTrashFilePath(trashPath, newFilename))
		existsTrashInfo, _ := afero.Exists(fs, buildTrashInfoPath(trashPath, newFilename))
		isDuplicateFilename = existsTrashFile || existsTrashInfo
	}
	return newFilename
}

func buildTrashFilePath(trashPath string, filename string) string {
	return trashPath + "/files/" + filename
}
