// Copyright (c) 2016 Readium Foundation
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation and/or
//    other materials provided with the distribution.
// 3. Neither the name of the organization nor the names of its contributors may be
//    used to endorse or promote products derived from this software without specific
//    prior written permission
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"encoding/json"
	"flag"
	"github.com/jjbier/lcpencrypt-pdf/consoleclient"
	"github.com/jjbier/lcpencrypt-pdf/encrypt"
	"github.com/jjbier/lcpencrypt-pdf/lcpclient"
	"github.com/readium/readium-lcp-server/lcpserver/api"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func showHelpAndExit() {
	log.Println("lcpencrypt-pdf protects an epub file for usage in an lcp environment")
	log.Println("-input         source epub file locator (file system or http GET)")
	log.Println("[-contentid]   optional content identifier, if omitted a new one will be generated")
	log.Println("[-output]      optional target location for protected content (file system or http PUT)")
	log.Println("[-lcpsv]       optional http endpoint for the License server")
	log.Println("[-console]     optional http endpoint for console server")
	log.Println("[-login]       login ( needed for License server) ")
	log.Println("[-password]    password ( needed for License server)")
	log.Println("[-help] :      help information")
	os.Exit(0)
	return
}

func exitWithError(lcpPublication apilcp.LcpPublication, err error, errorlevel int) {
	os.Stdout.WriteString(lcpPublication.ErrorMessage)
	os.Stdout.WriteString("\n")
	if err != nil {
		os.Stdout.WriteString(err.Error())
	}
	/* kept for future debug
	jsonBody, err := json.MarshalIndent(lcpPublication, " ", "  ")
	if err != nil {
		os.Stdout.WriteString("Error creating json lcpPublication\n")
		os.Exit(errorlevel)
	}
	os.Stdout.Write(jsonBody)
	os.Stdout.WriteString("\n")
	*/
	os.Exit(errorlevel)
}

func main() {
	var err error
	var addedPublication apilcp.LcpPublication
	var inputFilename = flag.String("input", "", "source pdf file locator (file system or http GET)")
	var contentid = flag.String("contentid", "", "optional content identifier; if omitted a new one is generated")
	var outputFilename = flag.String("output", "", "optional target location for the encrypted content (file system or http PUT)")
	var lcpsv = flag.String("lcpsv", "", "optional http endpoint of the License server (adds content)")
	var console = flag.String("console", "", "optional http endpoint of console server (adds content)")
	var username = flag.String("login", "", "login (License server)")
	var password = flag.String("password", "", "password (License server)")

	var help = flag.Bool("help", false, "shows information")

	if !flag.Parsed() {
		flag.Parse()
	}
	if *help {
		showHelpAndExit()
	}

	if *lcpsv != "" && (*username == "" || *password == "") {
		addedPublication.ErrorMessage = "incorrect parameters, lcpsv needs login and password, for more information type 'lcpencrypt-pdf -help' "
		exitWithError(addedPublication, nil, 80)
	}

	// check pdf input file content exits
	if _, err := os.Stat(*inputFilename); os.IsNotExist(err) {
		addedPublication.ErrorMessage = "Error opening input file, for more information type 'lcpencrypt-pdf -help' "
		exitWithError(addedPublication, err, 70)

	}

	if *contentid == "" { // contentID not set -> generate a new one
		uid, err_u := uuid.NewV4()
		if err_u != nil {
			exitWithError(addedPublication, err, 65)
		}
		*contentid = uid.String()
	}
	var basefilename string
	addedPublication.ContentId = *contentid

	// if the output file name not set,
	// then <content-id>.epub is created in the working directory
	if *outputFilename == "" {
		workingDir, _ := os.Getwd()
		*outputFilename = strings.Join([]string{workingDir, string(os.PathSeparator), *contentid, ".pdf"}, "")
		basefilename = filepath.Base(*inputFilename)
	} else {
		basefilename = filepath.Base(*outputFilename)
	}
	addedPublication.ContentDisposition = &basefilename
	// the output path must be accessible from the license server
	addedPublication.Output = *outputFilename

	encryptedPdf, err := encrypt.EncryptPdf(*inputFilename, *outputFilename)

	if err != nil || (encryptedPdf.Size == 0) {
		addedPublication.ErrorMessage = "Error encrypted the pdf file"
		exitWithError(addedPublication, err, 30)

	}

	addedPublication.Size = &encryptedPdf.Size
	addedPublication.Checksum = &encryptedPdf.Checksum
	addedPublication.ContentKey = encryptedPdf.EncryptionKey
	addedPublication.ContentDisposition = &encryptedPdf.Path

	//// notify the LCP Server
	if *lcpsv != "" {
		err = lcpclient.NotifyLcpServer(*lcpsv, *contentid, addedPublication, *username, *password)
		if err != nil {
			addedPublication.ErrorMessage = "Error notifying the License Server"
			exitWithError(addedPublication, err, 20)
		} else {
			os.Stdout.WriteString("License Server was notified\n")
		}
	}

	//// notify the Console Server
	if *console != "" {
		err = consoleclient.Notify(*console, *contentid, addedPublication)
		if err != nil {
			addedPublication.ErrorMessage = "Error notifying the Console Server"
			exitWithError(addedPublication, err, 20)
		} else {
			os.Stdout.WriteString("License Server was notified To Console\n")
		}
	}

	// write a json message to stdout for debug purpose
	jsonBody, err := json.MarshalIndent(addedPublication, " ", "  ")
	if err != nil {
		addedPublication.ErrorMessage = "Error creating json addedPublication"
		exitWithError(addedPublication, err, 10)
	}
	os.Stdout.Write(jsonBody)
	//os.Stdout.WriteString("\nEncryption was successful\n")
	os.Exit(0)
}
