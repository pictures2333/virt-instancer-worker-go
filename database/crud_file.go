package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Create a FileLink and it's FileObj (new-create)
func CreateFileLink(
	// FileLink
	filename string,
	// FileObj
	filenameReal string,
	lastModified time.Time,
	status string,
) (err error) {
	ctx := context.Background()

	fileobj := FileObj{
		FilenameReal: filenameReal,
		LastModified: lastModified,
		Status:       status,
	}

	filelink := FileLink{
		Filename: filename,
		FileObj:  fileobj,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return gorm.G[FileLink](tx).Create(ctx, &filelink)
	})

	return err
}

// Create a FileObj and replace old FileObj in FileLink
func CreateFileObj(
	// FileLink
	filename string,
	// FileObj
	filenameReal string,
	lastModified time.Time,
	status string,
) (err error) {
	ctx := context.Background()

	// new fileobj
	fileobj := FileObj{
		Filename:     filename,
		FilenameReal: filenameReal,
		LastModified: lastModified,
		Status:       status,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 「裡面」的 err 跟「外面」的 err 沒有關係

		// unlink the old FileObj from the FileLink
		// set oldFileObj.Filename = nil -> orphan
		if _, err := gorm.G[FileObj](tx).Where("filename = ?", filename).Update(ctx, "filename", nil); err != nil {
			return err
		}

		// link the new FileObj to the FileLink
		if err := gorm.G[FileObj](tx).Create(ctx, &fileobj); err != nil {
			return err
		}

		return nil
	})

	return err
}

// Read FileLink
func ReadFileLink(filename *string) (result []FileLink, err error) {
	ctx := context.Background()

	query := gorm.G[FileLink](db).Preload("FileObj", func(db gorm.PreloadBuilder) error {
		return nil
	})

	if filename != nil {
		query = query.Where("filename = ?", filename)
	}

	result, err = query.Find(ctx)

	return result, err
}

// Read FileObj which do not link to Filelinks (orphan FileObj)
func ReadFileObjOrphan() (result []FileObj, err error) {
	ctx := context.Background()

	result, err = gorm.G[FileObj](db).Preload("Placeholders", func(db gorm.PreloadBuilder) error {
		return nil
	}).Where("filename IS NULL").Find(ctx)

	return result, err
}

// Update status of FileObj
func UpdateFileObjStatus(filenameReal string, status string) (err error) {
	ctx := context.Background()

	err = db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[FileObj](tx).Where("filename_real = ?", filenameReal).Update(ctx, "status", status); err != nil {
			return err
		}
		return nil
	})

	return err
}

// Delete a FileLink and turn a FileObj into orphan (old-delete)
func DeleteFileLink(filename string) (err error) {
	ctx := context.Background()

	err = db.Transaction(func(tx *gorm.DB) error {
		// unlink FileObj from the FileLink
		if _, err := gorm.G[FileObj](tx).Where("filename = ?", filename).Update(ctx, "filename", nil); err != nil {
			return err
		}

		// delete the FileLink
		if _, err := gorm.G[FileLink](tx).Where("filename = ?", filename).Delete(ctx); err != nil {
			return err
		}

		return nil
	})

	return err
}

// Delete a FileObj which
// (1) do not link to any placeholders
// (2) do not link to any FileLinks
func DeleteFileObjOrphan(filenameReal string) (err error) {
	ctx := context.Background()

	err = db.Transaction(func(tx *gorm.DB) error {
		// (2) check exists and it's orphan
		fileobjs, err := gorm.G[FileObj](tx).Where("filename IS NULL AND filename_real = ?", filenameReal).Find(ctx)
		if err != nil {
			return err
		}
		if len(fileobjs) != 1 {
			return fmt.Errorf("FileObj %s not found or it's not an orphan FileObj", filenameReal)
		}
		fileobj := fileobjs[0]

		// (1) check placeholders
		if len(fileobj.Placeholders) != 0 {
			return fmt.Errorf("FileObj %s is busy", filenameReal)
		}

		// delete
		if _, err := gorm.G[FileObj](tx).Where("filename IS NULL AND filename_real = ?", filenameReal).Delete(ctx); err != nil {
			return err
		}

		return nil
	})

	return err
}
