// Author: Gergely Födémesi fgergo@gmail.com

/*
Module fswalker implements a concurrent filesystem reader to list all regular files or directories in a subtree.

	Synopsis

	// Reading from names results in all regular files and directories under /tmp
	names := fswalker.FileInfos("/tmp", fswalker.NeedFiles | fswalker.NeedDirs)

	...
	for p, ok := <-names; ok; p, ok = <-names {
		fmt.Println(p)
	}
*/
package fswalker

import (
	"os"
	"path/filepath"
)

const (
	NeedFiles = 1 << iota
	NeedDirs  = 1 << iota
)

// Reading from chan returned by fswalker.Names(path, opts)
// walks the filesystem using filepath.Walk() under subtree path
// resulting in names depending from opts.
// When no more files or directories are available the chan is closed.
// Errors are ignored.
func Names(path string, opts byte) chan string {
	names := make(chan string)

	go func() {
		filepath.Walk(path, func(wpath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			switch {
			case opts&NeedFiles == NeedFiles && info.Mode().IsRegular():
				names <- wpath
			case opts&NeedDirs == NeedDirs && info.Mode().IsDir():
				names <- wpath
			}

			return nil
		})
		close(names)
	}()

	return names
}
