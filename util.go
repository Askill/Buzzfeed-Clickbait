package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func printThesaurus(thes Thesaurus) {
	for i, s := range thes.pntrmap {
		fmt.Println(i, s.edges)
		fmt.Println("####################################################################")
	}
}

func trim(word string) string {
	word1 := strings.Replace(word, "\"", "", -1)
	word1 = strings.Replace(word1, "'", "", -1)

	word1 = strings.Replace(word1, ".", "", -1)
	word1 = strings.Replace(word1, "!", "", -1)
	word1 = strings.Replace(word1, "?", "", -1)
	word1 = strings.Replace(word1, ":", " ", -1)
	word1 = strings.Replace(word1, "#", " ", -1)
	word1 = strings.Replace(word1, ",", " ", -1)
	word1 = strings.Replace(word1, "(", " ", -1)
	word1 = strings.Replace(word1, ")", " ", -1)
	word1 = strings.Replace(word1, "”", " ", -1)
	word1 = strings.Replace(word1, "“", " ", -1)
	word1 = strings.Replace(word1, " ", "", -1)
	return word1
}
func contains(slc *[]string, str string) bool {
	for _, x := range *slc {
		if x == str {
			return true
		}
	}
	return false
}

var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
var lock sync.Mutex

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// Load loads the file at path into v.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, v)
}

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}
