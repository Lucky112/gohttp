package request

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"sync"
)

type HashRequestor struct {
	client *http.Client
	hasher hash.Hash
}

type HashResult struct {
	RawURL string
	Hash   string
	Err    error
}

func NewRequestor(hasher hash.Hash) *HashRequestor {
	req := &HashRequestor{
		client: &http.Client{},
		hasher: hasher,
	}

	return req
}

func (r *HashRequestor) Process(parallelism uint, urlCh <-chan string, resultCh chan<- HashResult, wg *sync.WaitGroup) {
	for i := 0; i < int(parallelism); i++ {
		wg.Add(1)
		go r.process(urlCh, resultCh, wg)
	}
}

func (r *HashRequestor) process(urlCh <-chan string, resultCh chan<- HashResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for rawURL := range urlCh {
		hashedResp, err := r.getHashedResponce(rawURL)

		resultCh <- HashResult{
			RawURL: rawURL,
			Hash:   hashedResp,
			Err:    err,
		}
	}
}

func (r *HashRequestor) getHashedResponce(rawURL string) (string, error) {
	resp, err := r.makeRequest(rawURL)
	if err != nil {
		return "", err
	}

	hashStr := r.makeHash(resp)

	return hashStr, nil
}

func (r *HashRequestor) makeHash(data []byte) string {
	r.hasher.Write(data)
	checkSum := r.hasher.Sum(nil)
	return hex.EncodeToString(checkSum)
}

func (r *HashRequestor) makeRequest(rawURL string) ([]byte, error) {
	request, _ := http.NewRequest("GET", rawURL, nil)
	data, err := r.execRequest(request)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (r *HashRequestor) execRequest(request *http.Request) (data []byte, err error) {
	resp, err := r.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("the request failed, err=%v", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
