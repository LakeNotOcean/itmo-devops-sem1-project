package handlers

import (
	"os"
	"sem1-final-project-hard-level/internal/dto"
)

func HandlePricesFile(tempFile *os.File, format dto.FormatType) (*dto.UploadPricesResult, error) {
	switch format {
	case dto.FormatZip:
	case dto.FormatTar:
		()
	default:
		return fmt.Errorf("invalid format: %s, allowed values: zip, tar", str)
	}
}

func processTarFile() {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	// Ищем файл data.csv в архиве
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if header.Name == "data.csv" || filepath.Base(header.Name) == "data.csv" {
			return h.processCSV(tarReader)
		}
	}

	return nil, fmt.Errorf("data.csv not found in archive")
}
func handleCSV() {

}
