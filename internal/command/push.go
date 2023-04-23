package command

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/mundanelizard/envi/internal/server"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Push() *cli.Command {
	return &cli.Command{
		Name:   "push",
		Action: handlePush,
	}
}

func handlePush(values *cli.ActionArgs, args []string) {
	ok, err := server.RetrieveUser()
	if err != nil {
		logger.Fatal(err)
		return
	}

	if !ok {
		fmt.Println("Authenticate with a server inorder to create a new envi repository")
		return
	}

	raw, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	repo := string(raw)

	ok, err = server.CheckAccess(repo)
	if err != nil {
		logger.Fatal(err)
		return
	}

	if !ok {
		fmt.Println("Authorised access to repository. You can't update repository")
		return
	}

	pushes, err := server.PushCount(repo)
	if err != nil {
		logger.Fatal(err)
		return
	}

	if pushes == 0 {
		err = handleInitialPush(repo)
	} else {
		err = handleSubsequentPush(repo)
	}

	if err != nil {
		logger.Fatal(err)
		return
	}
}

func handleInitialPush(repo string) error {
	zipDir := path.Join(os.TempDir(), path.Base(repo)+".envi.temp.zip")
	err := compressEnviDir(zipDir)
	if err != nil {
		return err
	}

	// todo: encrypt zip file

	return uploadZipToServer(zipDir)
}

func uploadZipToServer(zipDir string) error {
	url := "http://localhost:9000"

	file, err := os.Open(zipDir)
	if err != nil {
		return err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "envi.zip")
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	data := make(map[string]interface{})
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return err
	}

	fmt.Println(data["message"])
	return nil
}

func compressEnviDir(zipDir string) error {
	zipFile, err := os.Create(zipDir)
	if err != nil {
		return err
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

	return err
}

func handleSubsequentPush(repo string) error {

	return nil
}
