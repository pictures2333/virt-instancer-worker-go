package database

import (
	"context"

	"gorm.io/gorm"
)

func ReadPort() (result []Port, err error) {
	ctx := context.Background()

	result, err = gorm.G[Port](db).Find(ctx)
	return result, err
}

func AllocatePort(VMUUID string, port int) (err error) {
	ctx := context.Background()

	portobj := Port{
		VMUUID: VMUUID,
		Port:   port,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return gorm.G[Port](tx).Create(ctx, &portobj)
	})

	return err
}

func FreePort(port int) (err error) {
	ctx := context.Background()

	err = db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[Port](tx).Where("port = ?", port).Delete(ctx); err != nil {
			return err
		}
		return nil
	})

	return err
}
