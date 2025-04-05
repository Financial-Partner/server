package dto

type GachaResponse struct {
	ID     string `json:"id" example:"60d6ec33f777b123e4567890"`
	ImgSrc string `json:"img_src" example:"https://example.com/image.png"`
}

type DrawGachaRequest struct {
	Amount int64 `json:"amount" example:"100" binding:"required"`
}

type PreviewGachasResponse struct {
	Gachas []GachaResponse `json:"gachas"`
}
