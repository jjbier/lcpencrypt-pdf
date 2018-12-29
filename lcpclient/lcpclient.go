package lcpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/readium/readium-lcp-server/lcpserver/api"
	"net/http"
)

// notification of newly added content (Publication)
func NotifyLcpServer(lcpService, contentid string, lcpPublication apilcp.LcpPublication, username string, password string) error {
	//exchange encryption key with lcp service/content/<id>,
	//Payload:
	//  content-id: unique id for the content
	//  content-encryption-key: encryption key used for the content
	//  protected-content-location: full path of the encrypted file
	//  protected-content-length: content length in bytes
	//  protected-content-sha256: content sha
	//  protected-content-disposition: encrypted file name
	//fmt.Printf("lcpsv = %s\n", *lcpsv)
	var urlBuffer bytes.Buffer
	urlBuffer.WriteString(lcpService)
	urlBuffer.WriteString("/contents/")
	urlBuffer.WriteString(contentid)

	jsonBody, err := json.Marshal(lcpPublication)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", urlBuffer.String(), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if (resp.StatusCode != 302) && (resp.StatusCode/100) != 2 { //302=found or 20x reply = OK
		return fmt.Errorf("lcp server error %d", resp.StatusCode)
	}

	return nil
}
