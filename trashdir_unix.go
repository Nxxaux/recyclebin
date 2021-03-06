// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !windows

package recyclebin

import (
	"errors"
	"github.com/spf13/afero"
	"os"
	"strconv"
)

type envStorage interface {
	Getenv(name string) string
}

type osEnvStorage int

func (_ osEnvStorage) Getenv(name string) string {
	return os.Getenv(name)
}

func getTrashDirectory(filepath string, envStorage envStorage, uid int) (string, error) {
	if isExternalDevice(filepath) {
		deviceTrashPath, err := getDeviceTrashDirectory(filepath, uid)
		if err != nil {
			return "", err
		}
		return deviceTrashPath, nil
	}

	dataHomeDirectory := getDataHomeDirectory(envStorage)
	homeTrashPath, err := getHomeTrashDirectory(dataHomeDirectory)
	if err != nil {
		return "", errors.New("cannot find or create any trash directory")
	}
	return homeTrashPath, nil
}

func isExternalDevice(filepath string) bool {
	return false
}

func getHomeTrashDirectory(dataHomeDirectory string) (string, error) {
	homeTrashPath := dataHomeDirectory + "/Trash"
	hasHomeTrash, _ := afero.DirExists(fs, homeTrashPath)
	if hasHomeTrash {
		//add by Nxxaux-----------------------
		/*
			.local/share/Trash/files: permission denied fixed !!
		*/
		fs := afero.NewOsFs()
		_, err := fs.Stat(homeTrashPath)
		if err != nil {
			return "", err
		}
		err = fs.Chmod(homeTrashPath, 0700) //drwx --- ---
		if err != nil {
			return "", err
		}
		//------------------------------------
		return homeTrashPath, nil
	}
	err := fs.MkdirAll(homeTrashPath, os.ModeDir)
	return homeTrashPath, err
}

func getDataHomeDirectory(envStorage envStorage) string {
	XDG_DATA_HOME := envStorage.Getenv("XDG_DATA_HOME")
	if XDG_DATA_HOME == "" {
		HOME := envStorage.Getenv("HOME")
		return HOME + "/.local/share"
	}
	return XDG_DATA_HOME
}

func getDeviceTrashDirectory(partitionRootPath string, uid int) (string, error) {
	topTrashPath := partitionRootPath + "/.Trash"
	hasTrash, _ := afero.DirExists(fs, topTrashPath)
	if !hasTrash {
		topTrashUidPath := ".Trash-" + strconv.Itoa(uid)
		if err := fs.Mkdir(topTrashUidPath, os.ModeDir); err != nil {
			return "", err
		}
		return topTrashUidPath, nil
	}

	if isSymlink(topTrashPath) {
		return "", errors.New("device's top .Trash directory is a symbolic link")
	}

	uidTrashPath := topTrashPath + strconv.Itoa(uid)
	hasUidTrash, _ := afero.DirExists(fs, uidTrashPath)
	if !hasUidTrash {
		if err := fs.Mkdir(uidTrashPath, os.ModeDir); err != nil {
			return "", err
		}
	}
	return uidTrashPath, nil
}

func isSymlink(path string) bool {
	file, _ := fs.Stat(path)
	return file.Mode()&os.ModeSymlink == 0
}
