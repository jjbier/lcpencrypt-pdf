package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/readium/readium-lcp-server/crypto"
	"io/ioutil"
	"os"
)

// EncryptedPdf Encrypted pdf
type EncryptedPdf struct {
	Path          string
	EncryptionKey []byte
	Size          int64
	Checksum      string
}

// DecryptedPdf plain pdf
type DecryptedPdf struct {
	Path     string
	Size     int64
	Checksum string
}

func GetChecksum(path string) (string, error) {
	hasher := sha256.New()
	s, err := ioutil.ReadFile(path)
	_, err = hasher.Write(s)
	if err != nil {
		return "", errors.New("Unable to build checksum")
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil

}
func DecryptPdf(inputPath string, outputPath string, encryptionKey []byte) (DecryptedPdf, error) {
	if _, err := os.Stat(inputPath); err != nil {
		return DecryptedPdf{}, errors.New("Input file does not exists")
	}

	// Open Read file
	input, err := os.Open(inputPath)
	if err != nil {
		return DecryptedPdf{}, errors.New("Unable to read input file")
	}

	// Create output file
	output, err := os.Create(outputPath)
	if err != nil {
		return DecryptedPdf{}, errors.New("Unable to create output file")
	}

	// decrypt the pdf content, fill the output file
	encrypter := crypto.NewAESEncrypter_PUBLICATION_RESOURCES()
	decrypter, ok := encrypter.(crypto.Decrypter)
	if !ok {
		return DecryptedPdf{}, errors.New("Unable to create decrypter")

	}

	errDecrypted := decrypter.Decrypt(encryptionKey, input, output)

	if errDecrypted != nil {
		return DecryptedPdf{}, errors.New("Unable to decrypted file")
	}

	stats, err := output.Stat()
	if err != nil || (stats.Size() <= 0) {
		return DecryptedPdf{}, errors.New("Unable to output file")
	}

	checksum, err := GetChecksum(outputPath)

	if err != nil {
		return DecryptedPdf{}, errors.New("Unable to build checksum")
	}

	output.Close()
	input.Close()
	return DecryptedPdf{outputPath, stats.Size(), checksum}, nil
}

// EncryptPdf Encrypt input file to output file
func EncryptPdf(inputPath string, outputPath string) (EncryptedPdf, error) {
	if _, err := os.Stat(inputPath); err != nil {
		return EncryptedPdf{}, errors.New("Input file does not exists")
	}

	// Open Read file
	input, err := os.Open(inputPath)
	if err != nil {
		return EncryptedPdf{}, errors.New("Unable to read input file")
	}

	// Create output file
	output, err := os.Create(outputPath)
	if err != nil {
		return EncryptedPdf{}, errors.New("Unable to create output file")
	}

	// encrypt the pdf content, fill the output file
	encrypter := crypto.NewAESEncrypter_PUBLICATION_RESOURCES()
	encryptionKey, err := encrypter.GenerateKey()

	if err != nil {
		return EncryptedPdf{}, errors.New("Unable to create encryptionKey")
	}

	encryptError := encrypter.Encrypt(encryptionKey, input, output)

	if encryptError != nil {
		return EncryptedPdf{}, errors.New("Unable to encrypt file")
	}

	stats, err := output.Stat()
	if err != nil || (stats.Size() <= 0) {
		return EncryptedPdf{}, errors.New("Unable to output file")
	}

	checksum, err := GetChecksum(outputPath)

	if err != nil {
		return EncryptedPdf{}, errors.New("Unable to build checksum")
	}

	output.Close()
	input.Close()
	return EncryptedPdf{outputPath, encryptionKey, stats.Size(), checksum}, nil
}
