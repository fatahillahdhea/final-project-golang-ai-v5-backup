package service

import (
	"encoding/csv"
	"errors"
	"strings"
	"fmt"
	"strconv"
	
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	if fileContent == "" {
		return nil, errors.New("file content is empty")
	}

	// Map untuk menampung hasil pemrosesan
	result := make(map[string][]string)

	// Membaca konten sebagai CSV
	reader := csv.NewReader(strings.NewReader(fileContent))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("invalid CSV format")
	}

	// Pastikan ada header
	if len(rows) < 1 {
		return nil, errors.New("CSV file is empty or missing header")
	}

	// Ambil header
	headers := rows[0]
	for _, header := range headers {
		result[header] = []string{}
	}

	// Isi data
	for _, row := range rows[1:] {
		for i, value := range row {
			if i < len(headers) {
				result[headers[i]] = append(result[headers[i]], strings.TrimSpace(value))
			}
		}
	}

	return result, nil
}

func (s *FileService) ProcessEnergyFile(fileContent string) (map[string]float64, error) {
    if fileContent == "" {
        return nil, errors.New("file content is empty")
    }

    result := make(map[string]float64)
    reader := csv.NewReader(strings.NewReader(fileContent))
    rows, err := reader.ReadAll()
    if err != nil {
        return nil, errors.New("invalid CSV format")
    }

    if len(rows) < 1 {
        return nil, errors.New("CSV file is empty or missing header")
    }

    // Header diindikasikan sebagai: "Device,Usage"
    for _, row := range rows[1:] {
        if len(row) != 2 {
            continue // Skip baris yang tidak valid
        }
        device := strings.TrimSpace(row[0])
        usage, err := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
        if err != nil {
            return nil, fmt.Errorf("invalid usage data for device %s", device)
        }
        result[device] = usage
    }

    return result, nil
}

