package prices

import (
	"fmt"
	"sem1-final-project-hard-level/internal/dto"
	archivehelpers "sem1-final-project-hard-level/internal/handlers/prices/archive_helpers"

	"gorm.io/gorm"
)

func handlePricesFile(db *gorm.DB, filePath string, dataFileName string, batchSize int, format dto.FormatType) (*dto.UploadPricesResult, error) {
	switch format {
	case dto.FormatZip:
		return archivehelpers.HandleZipFile(db, filePath, dataFileName, batchSize)
	case dto.FormatTar:
		return archivehelpers.HandleTarFile(db, filePath, dataFileName, batchSize)
	default:
		return nil, fmt.Errorf("invalid format: %s, allowed values: zip, tar", format.String())
	}
}
