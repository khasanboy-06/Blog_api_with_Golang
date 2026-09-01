package dto

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200"`
	Content string `json:"content" binding:"required,min=10"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=3,max=200"`
	Content string `json:"content" binding:"omitempty,min=10"`
}