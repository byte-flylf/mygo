// generate cscope.out file for vim
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"algs"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: mkcscope <source-dir>")
		os.Exit(0)
	}

	var path, base string
	var err error

	if path, err = filepath.Abs(os.Args[1]); err != nil {
		log.Fatal(err)
	}
	base = filepath.Base(os.Args[1])

	var cwd string
	if cwd, err = os.Getwd(); err != nil {
		log.Fatal(err)
	}

	var p string = os.Getenv("HOME") + "/tmp/cscope/" + base
	if !algs.IsDirExists(p) {
		os.MkdirAll(p, 0700)
	}
	os.Chdir(p)
	tmpfile := "/tmp/cscope.files"
	if !algs.IsFile(tmpfile) {
		os.Create(tmpfile)
	}

	fmtstr := `find %s -name "*.[chSs]" -o -name "*.cpp" -o -name "*.cc" -o -name "*.hpp" > %s && cscope  -bkq -i %s`
	cmdstr := fmt.Sprintf(fmtstr, path, tmpfile, tmpfile)
	if err = exec.Command("bash", "-c", cmdstr).Run(); err != nil {
		log.Fatal(err)
	}

	os.Chdir(cwd)
}
