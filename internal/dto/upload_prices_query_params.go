package dto

import "sem1-final-project-hard-level/internal/enum"

type UploadPricesQueryParams struct {
	Format enum.FormatType `form:"type" validate:"enum_value=zip;zip;tar"`
}
