package database

import (
	"os"
)

func LogCreate(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0664)
}

func LogAppend(line string, fp *os.File) error {
	buf := []byte(line)
	buf = append(buf, '\n')

	_, err := fp.Write(buf)
	if err != nil {
		return err
	}
	return fp.Sync()
}
