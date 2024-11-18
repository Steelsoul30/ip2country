package dbgenerator

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/yl2chen/cidranger"

	"ip2country/pkg/store"
)

const (
	IPFile   = "GeoLite2-City-Blocks-IPv4.csv"
	CityFile = "GeoLite2-City-Locations-en.csv"
)

type DbGenerator struct {
	subnetInfo []store.SubnetInfo
	tree       cidranger.Ranger
}

type SubnetInfoCSV struct {
	Subnet      string `csv:"network"`
	CountryCode string `csv:"geoname_id"`
}

type CountryInfo struct {
	CountryCode string `csv:"geoname_id"`
	CountryName string `csv:"country_name"`
	CityName    string `csv:"city_name"`
}

func NewCustomRangerEntry(ipNet net.IPNet, data store.SubnetInfo) cidranger.RangerEntry {
	return &store.CustomTreeEntry{
		IPNet: ipNet,
		Info:  data,
	}
}

func NewDbGenerator() *DbGenerator {
	return &DbGenerator{}
}

func (s *DbGenerator) Close() {
	s.tree = nil
}

func (s *DbGenerator) DirectFromZip(zipFilePath string) (cidranger.Ranger, error) {
	err := s.UnzipAndPrepareData(zipFilePath)
	if err != nil {
		return nil, err
	}
	err = s.BuildCIDRTree()
	if err != nil {
		return nil, err
	}
	return s.tree, nil
}

func (s *DbGenerator) UnzipAndPrepareData(zipFilePath string) error {
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	var blocks []SubnetInfoCSV
	var locations []CountryInfo
	for _, file := range zipReader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s from zip: %v", file.Name, err)
		}

		if file.Name == IPFile {
			if err := gocsv.Unmarshal(zippedFile, &blocks); err != nil {
				return fmt.Errorf("failed to unmarshal blocks file: %v", err)
			}
		} else if file.Name == CityFile {
			if err := gocsv.Unmarshal(zippedFile, &locations); err != nil {
				return fmt.Errorf("failed to unmarshal blocks file: %v", err)
			}
		} else {
			err = fmt.Errorf("unknown file %s in zip", file.Name)
			slog.Error(err.Error())
			return err
		}
	}
	locationMap := make(map[string]CountryInfo)
	for _, location := range locations {
		locationMap[location.CountryCode] = location
	}
	subnetsInfo := make([]store.SubnetInfo, 0)
	for _, block := range blocks {
		if block.CountryCode == "" {
			continue
		}
		countryInfo, ok := locationMap[block.CountryCode]
		if !ok {
			slog.Error(fmt.Sprintf("Country code %s not found in locations", block.CountryCode))
			continue
		}
		subnetsInfo = append(subnetsInfo, store.SubnetInfo{
			Subnet:  block.Subnet,
			Country: countryInfo.CountryName,
			City:    countryInfo.CityName,
		})
	}
	s.subnetInfo = subnetsInfo
	return nil
}

// BuildCIDRTree Build a CIDR tree with additional data
func (s *DbGenerator) BuildCIDRTree() error {
	s.tree = cidranger.NewPCTrieRanger()
	for _, data := range s.subnetInfo {
		err := insertCIDR(s.tree, data.Subnet, data)
		if err != nil {
			slog.Error(fmt.Sprintf("Error inserting CIDR %s: %v\n", data.Subnet, err))
			return err
		}
	}
	return nil
}

func insertCIDR(tree cidranger.Ranger, cidr string, info store.SubnetInfo) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		slog.Error(fmt.Sprintf("invalid CIDR %s: %v", cidr, err))
		return fmt.Errorf("invalid CIDR %s: %v", cidr, err)
	}
	if err = tree.Insert(NewCustomRangerEntry(*ipNet, info)); err != nil {
		slog.Error(fmt.Sprintf("Error inserting CIDR %s: %v\n", cidr, err))
		return err
	}

	return nil
}

// SaveInfo saves the subnet info to a file
func (s *DbGenerator) SaveInfo(filename string) error {

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(s.subnetInfo); err != nil {
		return err
	}
	return os.WriteFile(filename, buf.Bytes(), 0600)
}

// LoadEntries loads the entries slice from a Gob file
func (s *DbGenerator) LoadEntries(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var entries []store.SubnetInfo
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&entries); err != nil {
		return err
	}
	s.subnetInfo = entries
	return nil
}

func (s *DbGenerator) TryLoadFromGob(filename string) (cidranger.Ranger, error) {
	err := s.LoadEntries(filename)
	if err != nil {
		return nil, err
	}
	info := s.subnetInfo
	if info == nil {
		return nil, fmt.Errorf("subnet info is nil")
	}
	err = s.BuildCIDRTree()
	if err != nil {
		return nil, err
	}
	return s.tree, nil
}
