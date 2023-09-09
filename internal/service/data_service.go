package service

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewDataService
type DataService struct{}

func NewDataService(service *DataService) (*DataService, error) {
	return service, nil
}
