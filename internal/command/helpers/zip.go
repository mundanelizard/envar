package helpers

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
)

func DecompressEnvironment(zipDir, dest string) error {
	err := os.MkdirAll(dest, 0655)
	if err != nil {
		return err
	}
	zipFile, err := zip.OpenReader(zipDir)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		filePath := filepath.Join(dest, file.Name)

		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, file.Mode())
		if err != nil && !os.IsExist(err) {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, file.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}

		err = targetFile.Close()
		if err != nil {
			return err
		}

		err = fileReader.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func CompressEnvironment(wd, repo string) (string, error) {
	zipFileName := path.Base(repo) + ".env.zip"
	zipDir := path.Join(os.TempDir(), zipFileName)

	zipFile, err := os.Create(zipDir)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	dirname := path.Join(wd, ".envi")
	err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open the file to add to the ZIP archive
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a new file header for the file in the ZIP archive
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the name of the file in the ZIP archive
		header.Name = path[len(dirname)+1:]

		// Add the file header to the ZIP archive
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy the contents of the file to the ZIP archive
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	return zipDir, err
}
