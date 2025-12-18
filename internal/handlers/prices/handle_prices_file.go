package prices

import (
	"fmt"
	"sem1-final-project-hard-level/internal/dto"
	"sem1-final-project-hard-level/internal/enum"
	archivehelpers "sem1-final-project-hard-level/internal/handlers/prices/helpers/archive_helpers"

	"gorm.io/gorm"
)

func handlePricesFile(db *gorm.DB, filePath string, dataFileName string, batchSize int, format enum.FormatType) (*dto.UploadPricesResult, error) {
	switch format {
	case enum.FormatZip:
		return archivehelpers.HandleZipFile(db, filePath, dataFileName, batchSize)
	case enum.FormatTar:
		return archivehelpers.HandleTarFile(db, filePath, dataFileName, batchSize)
	default:
		return nil, fmt.Errorf("invalid format: %s, allowed values: zip, tar", format.String())
	}
}
