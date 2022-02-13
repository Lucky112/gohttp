package main

import (
	"crypto/md5"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeRequest(t *testing.T) {
	r := NewRequestor(nil)

	data, err := r.makeRequest("http://google.com")
	require.NoError(t, err)
	require.NotEmpty(t, data)
}

func TestFailedRequest(t *testing.T) {
	r := NewRequestor(nil)

	data, err := r.makeRequest("http://google")
	require.Error(t, err)
	require.Empty(t, data)

	data, err = r.makeRequest("http://empty.com")
	require.Error(t, err)
	require.Empty(t, data)
}

func TestHash(t *testing.T) {
	r := NewRequestor(md5.New())

	data := []byte(`<html> some data 93 </html>`)
	expectedCheckSum := "3e2e53c357b0b857a1851892e3da7ccc"

	checksum := r.makeHash(data)
	require.Equal(t, expectedCheckSum, checksum)
}

func TestGetResponse(t *testing.T) {
	r := NewRequestor(md5.New())

	data, err := r.getHashedResponce("http://google.com")
	require.NoError(t, err)
	require.NotEmpty(t, data)
}
