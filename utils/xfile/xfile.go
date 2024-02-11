package xfile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blackRice-Tu/golib/utils/xos"
)

// GetExeDirectory ...
func GetExeDirectory() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// GetFileInfo get file info
func GetFileInfo(file string) (os.FileInfo, bool) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			return fileInfo, true
		}
		if os.IsNotExist(err) {
			return nil, false
		}
		return nil, false
	}
	return fileInfo, true
}

// CpFile cp file
func CpFile(oPath, dPath string) error {
	if len(dPath) == 0 || len(oPath) == 0 {
		return fmt.Errorf("CpFile: dPath is nil")
	}
	oFile, ok := GetFileInfo(oPath)
	if !ok {
		return fmt.Errorf("CpFile: oPath not exist")
	}
	if dFile, ok := GetFileInfo(dPath); ok {
		if dFile.Size() == oFile.Size() {
			return nil
		}
		shell := fmt.Sprintf("rm %s", dPath)
		_, err := xos.ExecShell(shell)
		if err != nil {
			return err
		}
	}
	// cp new file
	shell := fmt.Sprintf("cp -r %s %s", oPath, dPath)
	_, err := xos.ExecShell(shell)
	if err != nil {
		return err
	}
	// check md5
	dFile, ok := GetFileInfo(dPath)
	if ok {
		if dFile.Size() == oFile.Size() {
			return nil
		}
	}
	return fmt.Errorf("CpFileFailed:  %s [%d] == %s [%d] ", oPath, oFile.Size(), dFile, dFile.Size())
}

// DeleteFile delete file
func DeleteFile(oPath string) error {
	if len(oPath) == 0 || len(oPath) == 0 {
		return nil
	}
	if !FileIsExist(oPath) {
		return nil
	}
	shell := fmt.Sprintf("rm %s", oPath)
	_, err := xos.ExecShell(shell)
	if err != nil {
		return err
	}
	return nil
}

// FileIsExist check file is exist
func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func CreateDirectory(path string) error {
	if FileIsExist(path) {
		return nil
	}
	return os.MkdirAll(path, os.ModePerm)
}
