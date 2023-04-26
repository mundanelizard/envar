package helpers

func CompressAndEncryptRepo(wd, repo, secret string) (string, string, error) {
	comDir, err := CompressEnvironment(wd, repo)
	if err != nil {
		return "", comDir, err
	}

	encDir, err := encryptCompressedEnvironment(comDir, secret)
	if err != nil {
		return comDir, encDir, err
	}

	return comDir, encDir, err
}
