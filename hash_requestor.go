package main

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"sync"
)

type hashRequestor struct {
	client *http.Client
	hasher hash.Hash
}

type requestResult struct {
	address string
	hash    string
	err     error
}

func NewRequestor(hasher hash.Hash) *hashRequestor {
	req := &hashRequestor{
		client: &http.Client{},
		hasher: hasher,
	}

	return req
}

func (r *hashRequestor) Process(parallelism uint, addressCh <-chan string, resultCh chan<- requestResult, wg *sync.WaitGroup) {
	for i := 0; i < int(parallelism); i++ {
		wg.Add(1)
		go r.processSingleAddress(addressCh, resultCh, wg)
	}
}

func (r *hashRequestor) processSingleAddress(addressCh <-chan string, resultCh chan<- requestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for addr := range addressCh {
		hashedResp, err := r.getHashedResponce(addr)

		resultCh <- requestResult{
			address: addr,
			hash:    hashedResp,
			err:     err,
		}
	}
}

func (r *hashRequestor) getHashedResponce(addr string) (string, error) {
	resp, err := r.makeRequest(addr)
	if err != nil {
		return "", err
	}

	hashStr := r.makeHash(resp)

	return hashStr, nil
}

func (r *hashRequestor) makeHash(data []byte) string {
	r.hasher.Write(data)
	checkSum := r.hasher.Sum(nil)
	return hex.EncodeToString(checkSum)
}

func (r *hashRequestor) makeRequest(addr string) ([]byte, error) {
	request, _ := http.NewRequest("GET", addr, nil)
	data, err := r.execRequest(request)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (r *hashRequestor) execRequest(request *http.Request) (data []byte, err error) {
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
