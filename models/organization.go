package models

import (
	"github.com/gocql/gocql"
)

type Org struct {
	Base
	Name       string     `json:"name,omitempty"`
	Contact    string     `json:"contact,omitempty"`
	Type       string     `json:"type,omitempty"`
	BillAmount int64      `json:"bill_amount,omitempty"`
	Address    string     `json:"address,omitempty"`
	CCTVCount  int64      `json:"cctv_count,omitempty"`
	AdminID    gocql.UUID `json:"admin_id,omitempty"`
	UserCount  int64      `json:"user_count,omitempty"`
}

type OrgList struct {
	Count         int    `json:"count,omitempty"`
	Results       []*Org `json:"results,omitempty"`
	NextPageToken []byte `json:"next_page_token,omitempty"`
}
