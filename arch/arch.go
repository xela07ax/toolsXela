package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
// Zips "./input" into "./output.zip"
//"/home/droid/Projects/todo2/src/","zip.zip"
func main()  {
	root := "/home/droid/Projects/todo2/src/"
	ZipDir(root,"r.ZIP")
}
func ZipDir(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	err = RecursiveZip(pathToZip,destinationFile)
	time.Sleep(300*time.Millisecond)
	return err
}
func RecursiveZip(pathToZip string, destinationW io.Writer) error {

	myZip := zip.NewWriter(destinationW)
	err := filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}
//
//func main()  {
//	ZipFolder("/home/droid/Projects/todo2/src/","zip.zip")
//	time.Sleep(400*time.Millisecond)
//
//}
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest) + string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func ZipFolder(baseFolder string, out string) error {
	// Get a Buffer to Write To
	outFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outFile.Close()
	// Create a new zip archive.
	w := zip.NewWriter(outFile)
	defer w.Close()
	ZipWriter(baseFolder,w)
	w.Close()
	return nil
}

func ZipWriter(baseFolder string, w *zip.Writer ) {
	// Add some files to the archive.
	addFiles(w, baseFolder, "")
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, file.Name() + "/")
		}
	}
}