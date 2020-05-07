package services

type IIndexService interface {
	GetIndex()(string,error)
}

type IndexService struct {

}

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (this *IndexService) GetIndex() (string,error)  {

	return "",nil
}

