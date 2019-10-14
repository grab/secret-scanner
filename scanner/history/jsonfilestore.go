package history

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type JSONFileStore struct {
	DataFile    *os.File
	RecordBytes []byte
	Records     map[string]*ScanHistory
}

func (fs *JSONFileStore) Initialize(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	recordBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fs.DataFile = file
	fs.RecordBytes = recordBytes

	var records []*ScanHistory
	err = json.Unmarshal(recordBytes, &records)
	if err != nil {
		return err
	}

	for _, record := range records {
		fs.Records[fmt.Sprintf("%s:%s", record.GitProvider, record.RepoID)] = record
	}

	//err = fs.HydrateRecords()
	//if err != nil {
	//	return err
	//}

	return nil
}

//func (fs *CSVFileStore) Close() error {
//	if fs.DataFile != nil {
//		err := fs.DataFile.Close()
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

//func (fs *CSVFileStore) FetchAllRecords() (records [][]string, err error) {
//	records, err = csv.NewReader(fs.DataFile).ReadAll()
//	if err != nil {
//		return [][]string{}, err
//	}
//
//	return records, nil
//}

//func (fs *CSVFileStore) HydrateRecords() error {
//	if len(fs.RecordStrings) > 0 {
//		var histories []*ScanHistory
//
//		for _, line := range fs.RecordStrings {
//			history := &ScanHistory{
//				GitProvider: line[0],
//				RepoID:      line[1],
//				CommitHash:  line[2],
//				CreatedAt:   line[3],
//			}
//
//			histories = append(histories, history)
//		}
//
//		fs.Records = histories
//	}
//
//	return nil
//}

//func (fs *CSVFileStore) Get(gitprovider, repoID string) *ScanHistory {
//
//}
//
//func (fs *CSVFileStore) Save(history *ScanHistory) error {
//	record := []string{history.GitProvider, history.RepoID, history.CommitHash, history.CreatedAt}
//
//	err := fs.Writer.Write(record)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
