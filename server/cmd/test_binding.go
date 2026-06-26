package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type TestUpdateFriendRequest struct {
	Name        string  `json:"name" binding:"max=50"`
	URL         string  `json:"url" binding:"omitempty,url,max=255"`
	Description *string `json:"description" binding:"max=500"`
	Avatar      *string `json:"avatar" binding:"omitempty,url,max=255"`
	Screenshot  *string `json:"screenshot" binding:"omitempty,url,max=255"`
	Sort        *int    `json:"sort" binding:"omitempty,min=1,max=10"`
	IsInvalid   *bool   `json:"is_invalid"`
	IsPending   *bool   `json:"is_pending"`
	TypeID      *uint   `json:"type_id" binding:"omitempty"`
	RSSUrl      *string `json:"rss_url" binding:"omitempty,url,max=500"`
	Accessible  *int    `json:"accessible"`
	IsFeatured  *bool   `json:"is_featured"`
}

func main() {
	v := validator.New()
	req := TestUpdateFriendRequest{}
	req.IsFeatured = new(bool)
	*req.IsFeatured = true
	err := v.Struct(req)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("OK")
	}
}
