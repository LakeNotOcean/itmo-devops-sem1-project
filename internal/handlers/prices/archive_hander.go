package prices

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sem1-final-project-hard-level/internal/dto"

	"gorm.io/gorm"
)

// функции различаются, хотя, возможно, общую часть можно вынести

func handleTarFile(db *gorm.DB, filePath string, dataFileName string, batchSize int) (*dto.UploadPricesResult, error) {
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

		if header.Name == dataFileName || filepath.Base(header.Name) == dataFileName {
			return handleCSV(db, tarReader, batchSize)
		}
	}

	return nil, fmt.Errorf("%s not found in archive", dataFileName)
}

func handleZipFile(db *gorm.DB, filePath string, dataFileName string, batchSize int) (*dto.UploadPricesResult, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if f.Name == dataFileName || filepath.Base(f.Name) == dataFileName {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			return handleCSV(db, rc, batchSize)
		}
	}

	return nil, fmt.Errorf("%s not found in archive", dataFileName)
}
