package database

import (
	"fmt"
	"math/rand"
	"os"
)

func SaveData1(path string, data []byte) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(data)
	return err
}

func SaveData2(path string, data []byte) error {
	temp := fmt.Sprintf("%s.tmp.%d", path, rand.Int())
	fp, err := os.OpenFile(temp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	defer fp.Close()
	_, err = fp.Write(data)

	if err != nil {
		os.Remove(temp)
		return err
	}

	return os.Rename(temp, path)
}

func SaveData3(path string, data []byte) error {
	temp := fmt.Sprintf("%s.tmp.%d", path, rand.Int())
	fp, err := os.OpenFile(temp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	defer fp.Close()
	_, err = fp.Write(data)

	if err != nil {
		os.Remove(temp)
		return err
	}
	err = fp.Sync()
	if err != nil {
		os.Remove(temp)
		return err
	}

	return os.Rename(temp, path)
}
