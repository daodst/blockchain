package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

/**

*/
func TestConnectivity(url string, timeout time.Duration) bool {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:HTTPTCP（http1.1，）
			MaxIdleConnsPerHost: 512,  //
		},
		Timeout: timeout, //Client,、response body;Timeout
	}
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	//log.Debug("cosmos",resp.Body)
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return true
}

//http
func HttpGet(url string, timeout time.Duration) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:HTTPTCP（http1.1，）
			MaxIdleConnsPerHost: 512,  //
		},
		Timeout: timeout, //Client,、response body;Timeout
	}
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	//log.Debug("cosmos",resp.Body)
	defer resp.Body.Close()
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respContent), err
}
