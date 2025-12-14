package dto

type GetPricesQueryParamsDto struct {
	Start string `form:"start" validate:"required,dateformat"`
	End   string `form:"end"   validate:"required,dateformat,daterange=Start"`
	Min   *int   `form:"min"   validate:"omitempty,gt=0"`
	Max   *int   `form:"max"   validate:"omitempty,gt=0,minmax=Min"`
}
