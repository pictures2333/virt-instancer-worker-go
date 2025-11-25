package database

import (
	"context"

	"gorm.io/gorm"
)

func ReadInstanceMacAddressAll() (result *[]int, err error) {
	ctx := context.Background()

	var instances []Instance
	if instances, err = gorm.G[Instance](db).Find(ctx); err != nil {
		return nil, err
	}

	result = new([]int)
	for _, i := range instances {
		*result = append(*result, i.MacAddress)
	}

	return result, nil
}

func ReadInstanceSubnetAll() (result *[]string, err error) {
	ctx := context.Background()

	var instances []Instance
	if instances, err = gorm.G[Instance](db).Find(ctx); err != nil {
		return nil, err
	}

	result = new([]string)
	for _, i := range instances {
		*result = append(*result, i.Subnet)
	}

	return result, nil
}
