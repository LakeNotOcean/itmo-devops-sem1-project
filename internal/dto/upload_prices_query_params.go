package dto

type UploadPricesQueryParams struct {
	Format FormatType `form:"format" validate:"enum_value=zip;zip;tar"`
}
