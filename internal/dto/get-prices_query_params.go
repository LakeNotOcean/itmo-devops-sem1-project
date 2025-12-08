package dto

type GetPricesQueryParamsDto struct {
	Start string `query:"start" validate:"omitempty,dateformat"`
	End   string `query:"end"   validate:"omitempty,dateformat,daterange=Start"`
	Min   *int   `query:"min"   validate:"omitempty,gt=0"`
	Max   *int   `query:"max"   validate:"omitempty,gt=0,minmax=Min"`
}
