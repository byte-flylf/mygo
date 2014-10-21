package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

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
		line = strings.TrimSuffix(line, "\n")
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
