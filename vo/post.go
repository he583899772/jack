package vo

type CreatePostRequst struct {
	CategoryId uint `json:"category_id" binding:"required"`
	Title string `json:"title" binding:"required,max=100"`
	HeadImg string `json:"head_img"`
	Content string `json:"content" binding:"required"`
}
