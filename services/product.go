package services

import (
	"vision/daos"
	"vision/repositories"
	"vision/types"
)

type ProductService interface {
	GetProducts() []*daos.Product
	GetProductByID(productID *types.ID) (*daos.Product, error)
	CreateNewProduct(product *daos.Product) (*daos.Product, error)
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

func (p productService) GetProducts() []*daos.Product {
	//TODO implement me
	panic("implement me")
}

func (p productService) GetProductByID(productID *types.ID) (*daos.Product, error) {
	return p.ProductRepository.FindByID(productID)
}

func (p productService) CreateNewProduct(product *daos.Product) (*daos.Product, error) {
	product.Initialize()
	return p.ProductRepository.Save(product)
}

func (p productService) DeleteProduct(productID *types.ID) error {
	product := &daos.Product{Base: daos.Base{ID: productID}}
	product.Delete()
	return p.ProductRepository.Delete(product)
}

/*

func (s *CCTVService) GetAll(form *dtos.CCTVPagination) (*daos.CCTVList, error) {
	cctvList := make([]*daos.CCTV, 0)
	session := db.GetSession()
	query := session.Query(cctvTable.SelectAll()).BindStruct(&daos.CCTV{}).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&cctvList)
	return &daos.CCTVList{
		Count:         len(cctvList),
		Results:       cctvList,
		NextPageToken: iter.PageState(),
	}, err
}

func (s *CCTVService) Update(form *dtos.CCTVEdit) (*daos.CCTV, error) {
	cctv := daos.CCTV{
		Base:           daos.Base{UUID: form.UUID},
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
