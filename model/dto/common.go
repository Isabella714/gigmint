package dto

type PagingRequest struct {
	Page uint32 `query:"page" validate:"required,gt=0"`
	Size uint32 `query:"size" validate:"required,gt=0"`
}
