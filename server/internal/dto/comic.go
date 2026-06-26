package dto

// CreateComicRequest 创建番剧
type CreateComicRequest struct {
	Name  string `json:"name" binding:"required"`
	URL   string `json:"url" binding:"required"`
	Cover string `json:"cover"`
	Sort  int    `json:"sort"`
}

// UpdateComicRequest 更新番剧
type UpdateComicRequest struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Cover string `json:"cover"`
	Sort  *int   `json:"sort"`
}

// ComicResponse 番剧响应
type ComicResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Cover     string `json:"cover"`
	Sort      int    `json:"sort"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListComicRequest 番剧列表请求
type ListComicRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}
