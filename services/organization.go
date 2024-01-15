package services

import (
	"vision/daos"
	"vision/repositories"
	"vision/types"
)

type OrganizationService interface {
	GetOrganizations() ([]*daos.Organization, error)
	GetOrganizationByID(userID *types.ID) (*daos.Organization, error)
	CreateNewOrganization(data *daos.Organization) (*daos.Organization, error)
	DeleteOrganizationByID(userID *types.ID) error
}

type organizationService struct {
	OrganizationRepository repositories.OrganizationRepositoryInterface
}

func NewOrganizationService(
	organizationRepository repositories.OrganizationRepositoryInterface,
) OrganizationService {
	return &organizationService{
		OrganizationRepository: organizationRepository,
	}
}

func (o *organizationService) GetOrganizations() ([]*daos.Organization, error) {
	return o.OrganizationRepository.FindAll()
}

func (o *organizationService) GetOrganizationByID(organizationID *types.ID) (*daos.Organization, error) {
	return o.OrganizationRepository.FindByID(organizationID)
}

func (o *organizationService) CreateNewOrganization(organization *daos.Organization) (*daos.Organization, error) {
	organization.Initialize()
	err := o.OrganizationRepository.Save(organization)
	return organization, err
}

func (o *organizationService) DeleteOrganizationByID(organizationID *types.ID) error {
	organization := &daos.Organization{Base: daos.Base{ID: organizationID}}
	return o.OrganizationRepository.Delete(organization)
}

/*func (s *OrgService) GetAll(form *dtos.OrganizationPagination) (*daos.OrgList, error) {
	organizationList := make([]*daos.Org, 0)
	organization := daos.Org{
		Contact: form.Contact,
		Type:    form.Type,
		AdminID: form.AdminID,
	}
	session := db.GetSession()
	query := session.Query(orgTable.SelectAll()).BindStruct(&organization).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&organizationList)
	return &daos.OrgList{
		Count:         len(organizationList),
		Results:       organizationList,
		NextPageToken: iter.PageState(),
	}, err
}

func (s *OrgService) Update(form *dtos.OrganizationEdit) (*daos.Org, error) {
	organization := daos.Org{
		Base: daos.Base{ID: form.OrganizationID},
	}
	organization.Update()
	columns := []string{"name", "contact", "admin_id"}
	session := db.GetSession()
	err := session.Query(orgTable.Update(columns...)).BindStruct(&organization).ExecRelease()
	return &organization, err
}
*/
