package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	buf := new(bytes.Buffer)
	out := io.MultiWriter(buf, os.Stdout)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles) //
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(in io.Writer, path string, printFiles bool) error {
	return subDirTree(in, path, "", printFiles)
}

func subDirTree (in io.Writer, path string, prefix string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	if !printFiles {
		files = dirList(files)
	}
	for num, name := range files {
		subPrefix := prefix
		if printFiles || name.IsDir() {
			if num != len(files)-1 {
				fmt.Fprintln(in, subPrefix + "├───" + name.Name() + postName(&name))
				subPrefix += "│\t"
			} else {
				fmt.Fprintln(in, subPrefix + "└───" + name.Name() + postName(&name))
				subPrefix += "\t"
			}
			err := subDirTree(in, path+string(os.PathSeparator)+name.Name(), subPrefix, printFiles)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func postName(name fs.FileInfo) string {
	if !name.IsDir() {
		if name.Size() > 0 {
			return " (" + strconv.FormatInt(name.Size(), 10) + "b)"

		} else {
			return " (empty)"
		}
	}
	return ""
}

func dirList(files []fs.FileInfo) []fs.FileInfo {
	dirs := make([]fs.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}
	return dirs
}