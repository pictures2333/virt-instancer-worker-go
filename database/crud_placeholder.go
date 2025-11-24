package database

import (
	"context"

	"gorm.io/gorm"
)

// create a Placeholder
func CreatePlaceholder(filenameReal string) (err error) {
	ctx := context.Background()

	placeholder := Placeholder{
		FilenameReal: filenameReal,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[Placeholder](tx).Create(ctx, &placeholder); err != nil {
			return err
		}

		return nil
	})

	return err
}

// Delete a Placeholder
func DeletePlaceholder(filenameReal string) (err error) {
	ctx := context.Background()

	err = db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[Placeholder](tx).Where("filename_real = ?", filenameReal).Delete(ctx); err != nil {
			return err
		}

		return nil
	})

	return err
}
