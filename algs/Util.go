package algs

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ReadInts(path string) (numbers []int, err error) {
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var n int
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		parts := strings.Fields(strings.TrimSuffix(line, "\n"))
		if len(parts) == 0 {
			continue
		}
		for _, p := range parts {
			n, err = strconv.Atoi(p)
			if err != nil {
				return
			}
			numbers = append(numbers, n)
		}
	}
	return
}

func ReadLines(path string) (lines []string, err error) {
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	if err == io.EOF {
		err = nil
	}

	return
}

func WriteLines(lines []string, path string) (err error) {
	var file *os.File

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	for _, elem := range lines {
		_, err = file.WriteString(strings.TrimSpace(elem) + "\n")
		if err != nil {
			return
		}
	}
	return
}

func ReadUpToNLines(filename string, maxLines int) (int, []string) {
	var file *os.File
	var err error
	if file, err = os.Open(filename); err != nil {
		fmt.Println("failed to open the file: ", err)
	}
	defer file.Close()

	lines := make([]string, maxLines)
	reader := bufio.NewReader(file)
	i := 0
	for ; i < maxLines; i++ {
		line, err := reader.ReadString('\n')
		if line != "" {
			lines[i] = line
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("failed to finish reading the file: ", err)
		}
	} // Return the subslice actually used; could be < maxLines
	return i, lines[:i]
}

// check whether a directory denoted by a path
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	return false
}

// 复制文件到指定地方
func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer destFile.Close()
	return io.Copy(destFile, srcFile)
}

// 复制某个文件夹的文件到另一个文件夹
func CopyDir(srcDir, destDir string) {
	if IsDirExists(srcDir) {
		tmpSrc := strings.TrimSpace(srcDir)
		files, _ := ioutil.ReadDir(tmpSrc)
		for _, f := range files {
			CopyFile(srcDir+f.Name(), destDir+f.Name())
		}
	}
}

// check whether a file is exist
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.Mode().IsRegular()
	}
	return false
}

// reverse int slice
func ReverseIntSlice(s []int) {
	mid := len(s) / 2
	for i := 0; i < mid; i++ {
		j := len(s) - 1 - i
		s[i], s[j] = s[j], s[i]
	}
}
