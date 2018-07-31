package project

import (
	"archive/zip"
	"fmt"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/utils/log"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"bytes"
)

type FileLike interface {
	Name() string
	Stat() (os.FileInfo, error)
	Read([]byte) (int, error)
	Close() error
}

type SymLink struct {
	buf *bytes.Buffer
	info os.FileInfo
}

func (s SymLink) Name() string {
	return s.info.Name()
}

func (s SymLink) Read(p []byte) (int, error) {
	return s.buf.Read(p)
}

func (s SymLink) Close() error {
	return nil
}

func (s SymLink) Stat() (info os.FileInfo, err error) {
	return s.info, nil

}


// Archive project files according to the vgo project standard
func archiveProject(writer io.Writer, sourcePath, module, version string, excludePathsRegExp *regexp.Regexp) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}

		if excludePathsRegExp.FindString(path) != "" {
			log.Debug(fmt.Sprintf("Excluding path '%s' from zip archive.", path))
			return nil
		}
		fileName := getFileName(sourcePath, path, module, version)

		var file FileLike

		if info.Mode()&os.ModeSymlink != 0 {
			dst, err := os.Readlink(path)
			if err != nil {
				return err
			}
			i, err := os.Lstat(path)
			if err != nil {
				return err
			}
			file = SymLink{
				buf: bytes.NewBufferString(dst),
				info: i,
			}
		} else {
			file, err = os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
		}

		zipFile, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.CopyN(zipFile, file, info.Size())
		return err
	})
}

// getFileName composes filename for zip to match standard specified as
// module@version/{filename}
func getFileName(sourcePath, filePath, moduleName, version string) string {
	filename := strings.TrimPrefix(filePath, sourcePath)
	filename = strings.TrimLeft(filename, string(os.PathSeparator))
	moduleID := fmt.Sprintf("%s@%s", moduleName, version)

	return filepath.Join(moduleID, filename)
}
