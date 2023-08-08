package main

import (
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strings"
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
