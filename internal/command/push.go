package command

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/lockfile"
	"github.com/mundanelizard/envi/internal/refs"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Push() *cli.Command {
	return &cli.Command{
		Name:   "push",
		Action: handlePush,
	}
}

func handlePush(values *cli.ActionArgs, args []string) {
	db := database.New(path.Join(wd, ".envi", "objects"))
	rs := refs.New(path.Join(wd, ".envi", "refs"))

	_, err := srv.RetrieveUser()
	if err != nil {
		logger.Fatal(err)
		return
	}

	repo, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	pid, err := rs.Read()
	if err != nil {
		logger.Fatal(err)
		return
	}

	obj, err := db.Read(pid)
	if err != nil {
		logger.Fatal(err)
		return
	}

	commit, err := database.NewCommitFromByteArray(pid, obj)
	if err != nil {
		logger.Fatal(err)
		return
	}

	secret, _ := values.GetString("secret")
	fmt.Println("Encrypting repository with the key:", secret)
	comDir, encDir, err := compressAndEncryptRepo(string(repo), secret)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = srv.PushRepo(string(repo), commit.TreeId(), encDir, secret)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = os.Remove(encDir)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = os.Remove(comDir)
	if err != nil {
		logger.Fatal(err)
		return
	}

	fmt.Println("Environment pushed!")
}

func compressAndEncryptRepo(repo, secret string) (string, string, error) {
	comDir, err := compressEnvironment(string(repo))
	if err != nil {
		return "", comDir, err
	}

	encDir, err := encryptCompressedEnvironment(comDir, secret)
	if err != nil {
		return comDir, encDir, err
	}

	return comDir, encDir, err
}

func encryptCompressedEnvironment(dir, secret string) (string, error) {
	in, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer in.Close()

	data, err := io.ReadAll(in)
	if err != nil {
		return "", err
	}

	outDir := path.Join(dir, ".enc")
	lock := lockfile.New(outDir)
	err = lock.Hold()
	if err != nil {
		return "", err
	}
	defer lock.Commit()

	cphr, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(data))
	iv := make([]byte, aes.BlockSize)
	cipher.NewCFBEncrypter(cphr, iv).XORKeyStream(cipherText, data)

	err = lock.Write(cipherText)
	if err != nil {
		return "", err
	}

	err = lock.Commit()
	if err != nil {
		return "", err
	}

	return outDir, nil
}

func compressEnvironment(repo string) (string, error) {
	zipDir := path.Join(os.TempDir(), path.Base(repo)+".envi.temp.zip")

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
