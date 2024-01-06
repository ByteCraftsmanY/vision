package forms

import "github.com/gocql/gocql"

type OrganizationNew struct {
	Name    string     `json:"name,omitempty" binding:"required"`
	Contact string     `json:"contact,omitempty" binding:"required"`
	Type    string     `json:"type,omitempty" binding:"required"`
	Address string     `json:"address,omitempty" binding:"required"`
	AdminID gocql.UUID `json:"admin_id,omitempty" binding:"required"`
}

type OrganizationEdit struct {
	OrganizationID gocql.UUID `json:"organization_id,omitempty" binding:"required"`
	Name           string     `json:"name,omitempty"`
	Contact        string     `json:"contact,omitempty"`
	AdminID        gocql.UUID `json:"admin_id,omitempty"`
}

type OrganizationPagination struct {
	Type          string     `json:"type"`
	AdminID       gocql.UUID `json:"admin_id"`
	Contact       string     `json:"contact"`
	Size          int        `json:"size"`
	NextPageToken []byte     `json:"next_page_token"`
}
