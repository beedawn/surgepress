package assetcopy


import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyDir (srcDir, destDir string) error {
	info, err := os.Stat(srcDir)
	if errors.Is(err, fs.ErrNotExist){
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a diretory", srcDir)
	}

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		relativePath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, relativePath)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		return copyFile(path, target)

	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}