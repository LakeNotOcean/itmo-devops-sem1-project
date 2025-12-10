package dto

type GetPricesQueryParamsDto struct {
	Start string `query:"start" validate:"required,dateformat"`
	End   string `query:"end"   validate:"required,dateformat,daterange=Start"`
	Min   *int   `query:"min"   validate:"omitempty,gt=0"`
	Max   *int   `query:"max"   validate:"omitempty,gt=0,minmax=Min"`
}
