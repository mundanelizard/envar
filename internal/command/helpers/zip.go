package helpers

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
)

func DecompressEnvironment(zipDir, dest string) error {
	zipFile, err := zip.OpenReader(zipDir)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		filePath := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targeFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, file.Mode())
		if err != nil {
			return nil
		}

		_, err = io.Copy(targeFile, fileReader)
		if err != nil {
			return err
		}

		targeFile.Close()
		fileReader.Close()
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
