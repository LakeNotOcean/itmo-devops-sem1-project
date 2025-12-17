package dto

type UploadPricesQueryParams struct {
	Format FormatType `form:"type" validate:"enum_value=zip;zip;tar"`
}
