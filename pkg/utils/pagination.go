package utils

type Paging struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=5,max=100"`
}
