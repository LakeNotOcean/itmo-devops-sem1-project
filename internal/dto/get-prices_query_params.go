package dto

type GetPricesQueryParamsDto struct {
	Start string `form:"start" validate:"omitempty,dateformat"`
	End   string `form:"end"   validate:"omitempty,dateformat,daterange=Start"`
	Min   *int   `form:"min"   validate:"omitempty,gt=0"`
	Max   *int   `form:"max"   validate:"omitempty,gt=0,minmax=Min"`
}
