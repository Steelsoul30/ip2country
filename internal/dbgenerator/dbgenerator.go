package dbgenerator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

const (
	IPFile      = "GeoLite2-City-Blocks-IPv4.csv"
	CityFile    = "GeoLite2-City-Locations-en.csv"
	LocalDbFile = "dbLocal.dat"
)

type DbGenerator struct {
}

func NewDbGenerator() *DbGenerator {
	return &DbGenerator{}
}

func (s *DbGenerator) UnzipToMemory(zipFilePath string) (map[string][]byte, error) {
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	fileContents := make(map[string][]byte)

	for _, file := range zipReader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s from zip: %v", file.Name, err)
		}

		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, zippedFile)
		zippedFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s from zip: %v", file.Name, err)
		}
		fileContents[file.Name] = buffer.Bytes()
	}

	return fileContents, nil
}
