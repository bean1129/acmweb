package text

import (
	"os"
	"path/filepath"
)

type ZFile struct{}

func NewFile() *ZFile {
	return &ZFile{}
}

func (c *ZFile) Remove(path string) error {
	return os.RemoveAll(path)
}

func (c *ZFile) Create(path string) (*os.File, error) {
	dir := c.Dir(path)
	if !c.Exists(dir) {
		if err := c.Mkdir(dir); err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}

func (c *ZFile) Mkdir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (c *ZFile) Dir(path string) string {
	if path == "." {
		return filepath.Dir(c.RealPath(path))
	}
	return filepath.Dir(path)
}

func (c *ZFile) RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !c.Exists(p) {
		return ""
	}
	return p
}

func (c *ZFile) Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}
