package dtos

import "vision/types"

type OrganizationForm struct {
	Name             string    `json:"name,omitempty" binding:"required"`
	Contact          string    `json:"contact,omitempty" binding:"required"`
	Type             string    `json:"type,omitempty" binding:"required"`
	Address          string    `json:"address,omitempty" binding:"required"`
	AssociatedUserID *types.ID `json:"associated_user_id,omitempty" binding:"required"`
}

/*func (o *OrganizationForm) ConvertToEntity() *daos.Organization {
	return &daos.Organization{
		Name:    o.Name,
		Contact: o.Contact,
		Type:    o.Type,
		Address: o.Address,
		AdminID: o.AssociatedUserID,
	}
}*/

type OrganizationUserDTO struct {
	Base
	Name           string   `json:"name,omitempty"`
	Contact        string   `json:"contact,omitempty"`
	AssociatedUser *UserDTO `json:"associated_user,omitempty"`
}

type OrganizationEdit struct {
	OrganizationID types.ID `json:"organization_id,omitempty" binding:"required"`
	Name           string   `json:"name,omitempty"`
	Contact        string   `json:"contact,omitempty"`
	AdminID        types.ID `json:"admin_id,omitempty"`
}

type OrganizationPagination struct {
	Type          string   `json:"type"`
	AdminID       types.ID `json:"admin_id"`
	Contact       string   `json:"contact"`
	Size          int      `json:"size"`
	NextPageToken []byte   `json:"next_page_token"`
}

type OrgList struct {
	Count int `json:"count,omitempty"`
	//Results       []*Organization `json:"results,omitempty"`
	NextPageToken []byte `json:"next_page_token,omitempty"`
}
