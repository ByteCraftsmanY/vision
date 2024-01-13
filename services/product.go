package services

import (
	"vision/entities"
	"vision/repositories"
	"vision/types"
)

type ProductService interface {
	GetProducts() []*entities.Product
	GetProductByID(productID *types.ID) (*entities.Product, error)
	CreateNewProduct(product *entities.Product) (*entities.Product, error)
	DeleteProduct(productID *types.ID) error
}

type productService struct {
	ProductRepository repositories.ProductRepositoryInterface
}

func NewProductService(
	productRepository repositories.ProductRepositoryInterface,
) ProductService {
	return &productService{
		ProductRepository: productRepository,
	}
}

func (p productService) GetProducts() []*entities.Product {
	//TODO implement me
	panic("implement me")
}

func (p productService) GetProductByID(productID *types.ID) (*entities.Product, error) {
	return p.ProductRepository.FindByID(productID)
}

func (p productService) CreateNewProduct(product *entities.Product) (*entities.Product, error) {
	product.Initialize()
	return p.ProductRepository.Save(product)
}

func (p productService) DeleteProduct(productID *types.ID) error {
	product := &entities.Product{Base: entities.Base{ID: productID}}
	product.Delete()
	return p.ProductRepository.Delete(product)
}

/*

func (s *CCTVService) GetAll(form *dtos.CCTVPagination) (*entities.CCTVList, error) {
	cctvList := make([]*entities.CCTV, 0)
	session := db.GetSession()
	query := session.Query(cctvTable.SelectAll()).BindStruct(&entities.CCTV{}).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&cctvList)
	return &entities.CCTVList{
		Count:         len(cctvList),
		Results:       cctvList,
		NextPageToken: iter.PageState(),
	}, err
}

func (s *CCTVService) Update(form *dtos.CCTVEdit) (*entities.CCTV, error) {
	cctv := entities.CCTV{
		Base:           entities.Base{UUID: form.UUID},
		Username:       form.UserName,
		Password:       form.Password,
		URL:            form.URL,
		OrganizationID: form.OrganizationID,
		Extra:          form.Extra,
	}
	cctv.Update()
	columns := []string{"username", "password", "url", "organization_id", "extra"}
	session := db.GetSession()
	err := session.Query(cctvTable.Update(columns...)).BindStruct(&cctv).ExecRelease()
	return &cctv, err
}
*/
