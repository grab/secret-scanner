/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

// JSONFileStore is a JSON-based storage for scan histories
type JSONFileStore struct {
	DataFile *os.File
	Records  map[string]*History
}

// Initialize ...
func (fs *JSONFileStore) Initialize(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	recordBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fs.DataFile = file
	fs.Records = map[string]*History{}

	var records []*History
	err = json.Unmarshal(recordBytes, &records)
	if err != nil {
		return err
	}

	for _, record := range records {
		fs.Records[record.GetMapKey()] = record
	}

	return nil
}

// Close cleans up resources
func (fs *JSONFileStore) Close() {
	if fs.DataFile != nil {
		_ = fs.DataFile.Close()
	}
}

// GetDefaultStorePath returns default store file path,
// creates dir and file is not found
func (fs *JSONFileStore) GetDefaultStorePath() (string, error) {
	userHomeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	storeDirPath := path.Join(userHomeDir, DefaultStoreDir)
	storeFilePath := path.Join(storeDirPath, DefaultStoreFile)

	if _, err := os.Stat(storeFilePath); err != nil {
		if os.IsNotExist(err) {
			// attempt create dir and file
			err = fs.createDefaultStoreFile(storeDirPath, DefaultStoreFile)
			if err != nil {
				return "", err
			}

			return storeFilePath, nil
		}

		return "", err
	}

	return storeFilePath, nil
}

// Get retrieves history from store
func (fs *JSONFileStore) Get(gitprovider, repoID string) *History {
	val, exists := fs.Records[fmt.Sprintf("%s:%s", gitprovider, repoID)]
	if exists {
		return val
	}
	return nil
}

// Save persists records to file
func (fs *JSONFileStore) Save(history *History) error {
	fs.Records[history.GetMapKey()] = history

	var histories []*History

	for _, val := range fs.Records {
		histories = append(histories, val)
	}

	jsonBytes, err := json.Marshal(histories)
	if err != nil {
		return err
	}

	err = fs.DataFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = fs.DataFile.WriteAt(jsonBytes, 0)
	if err != nil {
		return err
	}

	return nil
}

func (fs *JSONFileStore) createDefaultStoreFile(dirPath, filename string) error {
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		return err
	}

	storeFile, err := os.Create(path.Join(dirPath, filename))
	if err != nil {
		return err
	}

	_, err = storeFile.Write([]byte("[]"))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(storeFile)

	return nil
}
