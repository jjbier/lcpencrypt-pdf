package encrypt

import (
	"os"
	"testing"
)

func TestEncryptPdf(t *testing.T) {
	inputPath, outputPath := "../test/samples/sample.pdf", os.TempDir()+"/encriptedFile.pdf"

	encryptedPdf, error := EncryptPdf(inputPath, outputPath)

	if error != nil {
		t.Error(error)
		t.FailNow()
	}

	if _, err := os.Stat(encryptedPdf.Path); err != nil {
		t.Error("expected a path")

	}

	if len(encryptedPdf.Checksum) <= 0 {
		t.Error("expected a Checksum")
	}

	if encryptedPdf.EncryptionKey == nil {
		t.Error("expected a key")
	}

	if len(encryptedPdf.EncryptionKey) <= 0 {
		t.Error("expected a key is empty")
		t.Log(encryptedPdf.EncryptionKey)
	}

	if encryptedPdf.Size <= 0 {
		t.Error("expected a Size")
	}
}

func TestDecryptPdf(t *testing.T) {
	inputPath, encryptPdfPath := "../test/samples/sample.pdf", os.TempDir()+"/encriptedFile.pdf"
	decryptInputPdf := os.TempDir() + "/decryptFile.pdf"

	checksumOriginalPDF, error := GetChecksum(inputPath)

	if error != nil {
		t.Error(error)
		t.FailNow()
	}

	encryptedPdf, error := EncryptPdf(inputPath, encryptPdfPath)

	if error != nil {
		t.Error(error)
		t.FailNow()
	}

	if encryptedPdf.EncryptionKey == nil {
		t.Error("expected a key")
	}

	if len(encryptedPdf.EncryptionKey) <= 0 {
		t.Error("expected a key is empty")
	}

	decryptedPdf, error := DecryptPdf(encryptedPdf.Path, decryptInputPdf, encryptedPdf.EncryptionKey)

	if error != nil {
		t.Error(error)
		t.FailNow()
	}

	if decryptInputPdf != decryptedPdf.Path {
		t.Error("decrypted Pdf path is not expected")
	}

	if len(decryptedPdf.Checksum) <= 0 {
		t.Error("expected a key is empty")
	}

	if decryptedPdf.Size <= 0 {
		t.Error("expected a Size")
	}

	checksumDecrytedPDF, error := GetChecksum(decryptedPdf.Path)

	if error != nil {
		t.Error(error)
		t.FailNow()
	}

	if checksumOriginalPDF != checksumDecrytedPDF {
		t.Error("The hash for decrypted pdf not math")
	}
}
