package models

import (
	"github.com/gocql/gocql"
)

type CCTV struct {
	Base
	Username       string            `json:"username,omitempty"`
	Password       string            `json:"-"`
	URL            string            `json:"url,omitempty"`
	OrganizationID gocql.UUID        `json:"organization_id,omitempty"`
	Extra          map[string]string `json:"extra,omitempty"`
}

type CCTVList struct {
	Count         int     `json:"count,omitempty"`
	Results       []*CCTV `json:"results,omitempty"`
	NextPageToken []byte  `json:"next_page_token,omitempty"`
}
