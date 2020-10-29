package zipcompress

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"mapconvertor/datastruct"
	"os"
	"path/filepath"
)

// CompressMap is to compress inputMap to destZip with zip.
func CompressMap(inputMap []datastruct.MapContent, destZip string) (err error) {
	// var err error
	var buf *bytes.Buffer = new(bytes.Buffer)
	defer func() {
		err = ioutil.WriteFile(destZip, buf.Bytes(), 0644)
	}()
	var zipWriter *zip.Writer = zip.NewWriter(buf)
	defer zipWriter.Close()
	var writer io.Writer
	for _, m := range inputMap {
		if m.NAME != "" {
			writer, err = zipWriter.Create(m.NAME)
			_, err = writer.Write(m.CONTENT)
		}
	}
	return err
}

// CompressFile is to compress srcDir to destZip with zip.
func CompressFile(srcDir string, destZip string) (err error) {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	rootDir, err := ioutil.ReadDir(srcDir)
	for _, fi := range rootDir {
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		header.Name = fi.Name()
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(filepath.Join(srcDir, fi.Name()))
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
	}
	return err
}
