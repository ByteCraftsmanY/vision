package services

import (
	"vision/entities"
	"vision/repositories"
	"vision/types"
)

type OrganizationService interface {
	GetOrganizations() ([]*entities.Organization, error)
	GetOrganizationByID(userID *types.ID) (*entities.Organization, error)
	CreateNewOrganization(data *entities.Organization) (*entities.Organization, error)
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

func (o *organizationService) GetOrganizations() ([]*entities.Organization, error) {
	return o.OrganizationRepository.FindAll()
}

func (o *organizationService) GetOrganizationByID(organizationID *types.ID) (*entities.Organization, error) {
	return o.OrganizationRepository.FindByID(organizationID)
}

func (o *organizationService) CreateNewOrganization(organization *entities.Organization) (*entities.Organization, error) {
	organization.Initialize()
	err := o.OrganizationRepository.Save(organization)
	return organization, err
}

func (o *organizationService) DeleteOrganizationByID(organizationID *types.ID) error {
	organization := &entities.Organization{Base: entities.Base{ID: organizationID}}
	return o.OrganizationRepository.Delete(organization)
}

/*func (s *OrgService) GetAll(form *dtos.OrganizationPagination) (*entities.OrgList, error) {
	organizationList := make([]*entities.Org, 0)
	organization := entities.Org{
		Contact: form.Contact,
		Type:    form.Type,
		AdminID: form.AdminID,
	}
	session := db.GetSession()
	query := session.Query(orgTable.SelectAll()).BindStruct(&organization).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&organizationList)
	return &entities.OrgList{
		Count:         len(organizationList),
		Results:       organizationList,
		NextPageToken: iter.PageState(),
	}, err
}

func (s *OrgService) Update(form *dtos.OrganizationEdit) (*entities.Org, error) {
	organization := entities.Org{
		Base: entities.Base{ID: form.OrganizationID},
	}
	organization.Update()
	columns := []string{"name", "contact", "admin_id"}
	session := db.GetSession()
	err := session.Query(orgTable.Update(columns...)).BindStruct(&organization).ExecRelease()
	return &organization, err
}
*/
