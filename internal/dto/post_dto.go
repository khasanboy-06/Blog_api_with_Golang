package dto

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200"`
	Content string `json:"content" binding:"required,min=10"`
	Author  string `json:"author" binding:"required,min=2,max=100"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=3,max=200"`
	Content string `json:"content" binding:"omitempty,min=10"`
	Author  string `json:"author" binding:"omitempty,min=2,max=100"`
}