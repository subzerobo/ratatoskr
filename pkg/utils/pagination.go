package utils

type Paging struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=5,max=100"`
}

type MorePaging struct {
	LastID uint `validate:"required,min=0"`
	Size   int `validate:"required,min=5,max=300"`
}
