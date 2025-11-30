package main

import (
	"compress/gzip"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type ImportedCsvFile struct {
	Reader *csv.Reader
	Closer func() error
}

func NewFileReader(filename string) (*ImportedCsvFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var (
		csvReader *csv.Reader
		gzReader  *gzip.Reader
		reader    io.Reader = file
	)

	if strings.HasSuffix(filename, ".gz") {
		gzReader, err = gzip.NewReader(reader)
		if err != nil {
			if closeErr := file.Close(); closeErr != nil {
				err = errors.Join(err, closeErr)
			}

			return nil, err
		}
		reader = gzReader
	}

	csvReader = csv.NewReader(reader)
	if isTsv(filename) {
		csvReader.LazyQuotes = true
		csvReader.Comma = '\t'
		csvReader.FieldsPerRecord = -1
	}

	return &ImportedCsvFile{
		Reader: csvReader,
		Closer: func() (cerr error) {
			if gzReader != nil {
				cerr = errors.Join(cerr, gzReader.Close())
			}

			return errors.Join(cerr, file.Close())
		},
	}, nil
}

func (c *ImportedCsvFile) GetColumnNames() ([]string, error) {
	header, err := c.Reader.Read()
	if err != nil {
		return nil, err
	}

	columns := make([]string, len(header))
	for i, col := range header {
		columns[i] = strings.TrimSpace(col)
	}

	return columns, nil
}

func isTsv(filename string) bool {
	base := filename
	if strings.HasSuffix(strings.ToLower(filename), ".gz") {
		base = strings.TrimSuffix(filename, ".gz")
	}
	return strings.HasSuffix(strings.ToLower(base), ".tsv")
}
