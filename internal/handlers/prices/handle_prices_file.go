package prices

import (
	"fmt"
	"sem1-final-project-hard-level/internal/dto"

	"gorm.io/gorm"
)

func handlePricesFile(db *gorm.DB, filePath string, dataFileName string, batchSize int, format dto.FormatType) (*dto.UploadPricesResult, error) {
	switch format {
	case dto.FormatZip:
		return handleZipFile(db, filePath, dataFileName, batchSize)
	case dto.FormatTar:
		return handleTarFile(db, filePath, dataFileName, batchSize)
	default:
		return nil, fmt.Errorf("invalid format: %s, allowed values: zip, tar", format.String())
	}
}
