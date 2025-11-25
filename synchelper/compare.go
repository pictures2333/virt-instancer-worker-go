package synchelper

import (
	"Instancer-worker-go/database"
	"Instancer-worker-go/schema"
	"fmt"
	"strings"
)

// find items in MinIO but not in DB (new files)
func find_new_create(a *[]schema.FileMINIO, b *[]database.FileLink) *[]schema.FileMINIO {
	m := make(map[string]bool)
	for _, f := range *b {
		m[f.Filename] = true
	}

	result := new([]schema.FileMINIO)
	for _, f := range *a {
		tmp := fmt.Sprintf("%s_%s", f.Bucket, strings.ReplaceAll(f.Filename, "/", "_"))
		if !m[tmp] {
			*result = append(*result, f)
		}
	}

	return result
}

// find items in db bit not in MinIO (old files)
func find_old_delete(a *[]database.FileLink, b *[]schema.FileMINIO) *[]database.FileLink {
	m := make(map[string]bool)
	for _, f := range *b {
		tmp := fmt.Sprintf("%s_%s", f.Bucket, strings.ReplaceAll(f.Filename, "/", "_"))
		m[tmp] = true
	}

	result := new([]database.FileLink)
	for _, f := range *a {
		if !m[f.Filename] {
			*result = append(*result, f)
		}
	}

	return result
}

// find items in MinIO and db (need to check update)
func find_same(a *[]schema.FileMINIO, b *[]database.FileLink) *[]schema.FileMINIO {
	m := make(map[string]bool)
	for _, f := range *b {
		m[f.Filename] = true
	}

	result := new([]schema.FileMINIO)
	for _, f := range *a {
		tmp := fmt.Sprintf("%s_%s", f.Bucket, strings.ReplaceAll(f.Filename, "/", "_"))
		if m[tmp] {
			*result = append(*result, f)
		}
	}

	return result
}
