package database

import (
	"Instancer-worker-go/schema"
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

// create
func CreateInstance(
	VMUUID string,
	networkUUID string, networkBrname string, subnet string,
	vmMacAddress int,
	data *schema.InstanceData,
) (err error) {
	ctx := context.Background()

	var dataBytes []byte
	if dataBytes, err = json.Marshal(data); err != nil {
		return err
	}
	dataStr := string(dataBytes)

	instance := Instance{
		// vm
		VMUUID: VMUUID,
		// network - bridge
		NetworkUUID:       networkUUID,
		NetworkBridgeName: networkBrname,
		Subnet:            subnet,
		// network - vm
		MacAddress: vmMacAddress,
		// data
		Data: dataStr,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return gorm.G[Instance](tx).Create(ctx, &instance)
	})

	return err
}

// read
func ReadInstance(VMUUID *string) (result []Instance, err error) {
	ctx := context.Background()

	query := gorm.G[Instance](db).Preload("Ports", func(db gorm.PreloadBuilder) error {
		return nil
	})

	if VMUUID != nil {
		query = query.Where("vm_uuid = ?", VMUUID)
	}

	result, err = query.Find(ctx)

	return result, err
}

// update

// delete
func DeleteInstance(VMUUID string) (err error) {
	ctx := context.Background()

	err = db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[Instance](tx).Where("vm_uuid = ?", VMUUID).Delete(ctx); err != nil {
			return err
		}
		return nil
	})

	return err
}
