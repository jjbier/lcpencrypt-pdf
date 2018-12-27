package main

import (
	"encoding/base64"
	"github.com/readium/readium-lcp-server/crypto"
	"log"
	"os"
	"strings"
)

func main() {


	os.Stdout.WriteString("\nEncryption PDF Tool \n")
	os.Stdout.WriteString("\n=================== \n")

	encrypter := crypto.NewAESEncrypter_PUBLICATION_RESOURCES()

	key, err := encrypter.GenerateKey()
	if err != nil {
		log.Println("Error generating a key")
		os.Exit(-1)

	}

	os.Stdout.WriteString("\n Key:" + base64.StdEncoding.EncodeToString(key))

	// current dir
	workingDir, _ := os.Getwd()

	inputFilename  := strings.Join([]string{workingDir, string(os.PathSeparator), "go-in-action.pdf"}, "")
	input, err := os.Open(inputFilename)
	if err != nil {
		log.Println("Error Loading file " + inputFilename)
		os.Exit(-1)

	}
	os.Stdout.WriteString("\n Input file name: " + inputFilename)
	outputEncryptedFilename := strings.Join([]string{workingDir, string(os.PathSeparator), "contentid", "-encrypted.pdf"}, "")

	// create an output file
	outputEncryptedFile, err := os.Create(outputEncryptedFilename)
	if err != nil {
		os.Exit(40)

	}

	os.Stdout.WriteString("\n Encrypted file name: " + outputEncryptedFilename)

	errorEncrypred := encrypter.Encrypt(key, input, outputEncryptedFile)

	if errorEncrypred != nil {
		log.Println("Error on encrypted file")
		os.Exit(-2)
	} else {
		os.Stdout.WriteString("\n Encrypted file  Success.\n")
	}


	outputDesEncryptedFilename := strings.Join([]string{workingDir, string(os.PathSeparator), "contentid", "-desencrypted.pdf"}, "")


	//// create an output file
	outputDesEncryptedFile, err := os.Create(outputDesEncryptedFilename)
	if err != nil {
		os.Exit(40)

	}




	os.Stdout.WriteString("\n DesEncrypted file name :" + outputDesEncryptedFilename)
	decrypter, ok := encrypter.(crypto.Decrypter)
	if !ok {
		os.Exit(41)

	}

	inputEncryptedFile, err := os.Open(outputEncryptedFilename)
	if err != nil {
		log.Println("Error Loading file " + inputFilename)
		os.Exit(-1)

	}

	errDecrypter := decrypter.Decrypt(key, inputEncryptedFile, outputDesEncryptedFile)

	if errDecrypter != nil {
		log.Println("Error on decrypter")
		os.Exit(43)

	}

	os.Stdout.WriteString("\n DesEncrypted file  Success.\n")
	os.Exit(0)
}
