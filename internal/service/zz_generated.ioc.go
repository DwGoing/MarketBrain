//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package service

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &dataService_{}
		},
	})
	dataServiceStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &DataService{}
		},
		ConstructFunc: func(i interface{}, _ interface{}) (interface{}, error) {
			impl := i.(*DataService)
			var constructFunc DataServiceConstructFunc = NewDataService
			return constructFunc(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
	}
	singleton.RegisterStructDescriptor(dataServiceStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &fundsService_{}
		},
	})
	fundsServiceStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &FundsService{}
		},
		ConstructFunc: func(i interface{}, _ interface{}) (interface{}, error) {
			impl := i.(*FundsService)
			var constructFunc FundsServiceConstructFunc = NewFundsService
			return constructFunc(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
	}
	singleton.RegisterStructDescriptor(fundsServiceStructDescriptor)
}

type DataServiceConstructFunc func(impl *DataService) (*DataService, error)
type FundsServiceConstructFunc func(impl *FundsService) (*FundsService, error)
type dataService_ struct {
}

type fundsService_ struct {
}

type DataServiceIOCInterface interface {
}

type FundsServiceIOCInterface interface {
}

var _dataServiceSDID string

func GetDataServiceSingleton() (*DataService, error) {
	if _dataServiceSDID == "" {
		_dataServiceSDID = util.GetSDIDByStructPtr(new(DataService))
	}
	i, err := singleton.GetImpl(_dataServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*DataService)
	return impl, nil
}

func GetDataServiceIOCInterfaceSingleton() (DataServiceIOCInterface, error) {
	if _dataServiceSDID == "" {
		_dataServiceSDID = util.GetSDIDByStructPtr(new(DataService))
	}
	i, err := singleton.GetImplWithProxy(_dataServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(DataServiceIOCInterface)
	return impl, nil
}

type ThisDataService struct {
}

func (t *ThisDataService) This() DataServiceIOCInterface {
	thisPtr, _ := GetDataServiceIOCInterfaceSingleton()
	return thisPtr
}

var _fundsServiceSDID string

func GetFundsServiceSingleton() (*FundsService, error) {
	if _fundsServiceSDID == "" {
		_fundsServiceSDID = util.GetSDIDByStructPtr(new(FundsService))
	}
	i, err := singleton.GetImpl(_fundsServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*FundsService)
	return impl, nil
}

func GetFundsServiceIOCInterfaceSingleton() (FundsServiceIOCInterface, error) {
	if _fundsServiceSDID == "" {
		_fundsServiceSDID = util.GetSDIDByStructPtr(new(FundsService))
	}
	i, err := singleton.GetImplWithProxy(_fundsServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(FundsServiceIOCInterface)
	return impl, nil
}

type ThisFundsService struct {
}

func (t *ThisFundsService) This() FundsServiceIOCInterface {
	thisPtr, _ := GetFundsServiceIOCInterfaceSingleton()
	return thisPtr
}
