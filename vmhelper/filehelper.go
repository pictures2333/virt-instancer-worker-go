package vmhelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"Instancer-worker-go/schema"
	"Instancer-worker-go/utils"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func getLibvirtGid() (libvirtGid int, err error) {
	var libvirtGroup *user.Group
	if libvirtGroup, err = user.LookupGroup(config.LibvirtGroup); err != nil {
		return -1, err
	}

	if libvirtGid, err = strconv.Atoi(libvirtGroup.Gid); err != nil {
		return -1, err
	}

	return libvirtGid, nil
}

func VMMustMkdir(path string) (err error) {
	var libvirtGid int
	if libvirtGid, err = getLibvirtGid(); err != nil {
		return err
	}

	// mkdir
	if err = utils.MustMkdir(path); err != nil {
		return err
	}

	// chown
	if err = os.Chown(path, -1, libvirtGid); err != nil {
		return err
	}

	// chmod
	if err = os.Chmod(path, 0o770); err != nil {
		return err
	}

	return nil
}

func VMMustRmdir(path string) (err error) {
	return utils.MustRmdir(path)
}

func VMCopyFile(VMUUID string, file *schema.File) (err error) {
	var libvirtGid int
	if libvirtGid, err = getLibvirtGid(); err != nil {
		return err
	}

	// get filename_real
	objNameReplaced := strings.ReplaceAll(file.Filename, "/", "_")
	filename := fmt.Sprintf("%s_%s", file.Bucket, objNameReplaced)

	var filelinks []database.FileLink
	if filelinks, err = database.ReadFileLink(&filename); err != nil {
		return err
	}
	if len(filelinks) != 1 {
		return fmt.Errorf("File %s/%s has no Filelink (or two or above Filelinks)", file.Bucket, file.Filename)
	}
	filelink := filelinks[0]
	if filelink.FileObj.Status != "READY" {
		return fmt.Errorf("File %s/%s is not ready", file.Bucket, file.Filename)
	}
	filenameReal := filelink.FileObj.FilenameReal

	// placeholder
	if err = database.CreatePlaceholder(filenameReal); err != nil {
		return fmt.Errorf("Failed to allocate a placeholder for %s : %v", filenameReal, err)
	}
	defer func(err *error) {
		if err := database.DeletePlaceholder(filenameReal); err != nil {
			utils.Showerr(fmt.Sprintf("Failed to delete the placeholder of %s : %v", filenameReal, err), false)
		}
	}(&err)

	// copy
	src := fmt.Sprintf("%s/sync/%s", config.FileDir, filenameReal)
	dst := fmt.Sprintf("%s/vmfiles/%s/%s", config.FileDir, VMUUID, filename)
	if err = utils.CopyFile(src, dst); err != nil {
		return nil
	}
	// chown
	if err = os.Chown(dst, -1, libvirtGid); err != nil {
		return err
	}
	// chmod
	if err = os.Chmod(dst, 0o770); err != nil {
		return err
	}

	return nil
}
