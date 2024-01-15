package dtos

import "vision/types"

type ProductForm struct {
	UserName       string             `json:"username" binding:"required"`
	Password       string             `json:"password" binding:"required"`
	URL            string             `json:"url" binding:"required"`
	OrganizationID *types.ID          `json:"organization_id" binding:"required"`
	Extra          *map[string]string `json:"extra"`
}

type CCTVEdit struct {
	UUID           *types.ID         `json:"uuid" binding:"required"`
	UserName       string            `json:"username"`
	Password       string            `json:"password"`
	URL            string            `json:"url"`
	OrganizationID *types.ID         `json:"organization_id"`
	Extra          map[string]string `json:"extra"`
}

type CCTVPagination struct {
	Size          int    `json:"size"`
	NextPageToken []byte `json:"next_page_token"`
}
