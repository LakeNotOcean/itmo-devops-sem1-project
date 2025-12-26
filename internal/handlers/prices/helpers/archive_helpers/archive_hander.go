package archivehelpers

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"sem1-final-project-hard-level/internal/dto"
	csvhelpers "sem1-final-project-hard-level/internal/handlers/prices/helpers/csv_helpers"

	"gorm.io/gorm"
)

// функции различаются, хотя, возможно, общую часть можно вынести

func HandleTarFile(db *gorm.DB, filePath string, dataFileName string, batchSize int) (*dto.UploadPricesResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// ищем первый csv-файл
		if header.Typeflag == tar.TypeDir || !isCSVFile(header.Name) {
			continue
		}
		return csvhelpers.ProcessCSV(db, tarReader, batchSize)
	}

	return nil, fmt.Errorf("%s not found in archive", dataFileName)
}

func HandleZipFile(db *gorm.DB, filePath string, dataFileName string, batchSize int) (*dto.UploadPricesResult, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		// ищем первый csv-файл
		if f.FileInfo().IsDir() || !isCSVFile(f.Name) {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		return csvhelpers.ProcessCSV(db, rc, batchSize)
	}

	return nil, fmt.Errorf("%s not found in archive", dataFileName)
}

func isCSVFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".csv"
}
