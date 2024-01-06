package forms

import "github.com/gocql/gocql"

type CCTV struct {
	UserName       string            `json:"username" binding:"required"`
	Password       string            `json:"password" binding:"required"`
	URL            string            `json:"url" binding:"required"`
	OrganizationID gocql.UUID        `json:"organization_id" binding:"required"`
	Extra          map[string]string `json:"extra"`
}

type CCTVEdit struct {
	UUID           gocql.UUID        `json:"uuid" binding:"required"`
	UserName       string            `json:"username"`
	Password       string            `json:"password"`
	URL            string            `json:"url"`
	OrganizationID gocql.UUID        `json:"organization_id"`
	Extra          map[string]string `json:"extra"`
}

type CCTVPagination struct {
	Size          int    `json:"size"`
	NextPageToken []byte `json:"next_page_token"`
}
