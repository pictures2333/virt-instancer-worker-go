package vmhelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/schema"
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type FileInfo struct {
	FullPath string
	Type     string
}

// be called in template
func GetFileInfo(data *schema.InstanceData, key string) (result FileInfo) {
	result.FullPath = ""
	result.Type = ""

	file, ok := data.Config.Files[key]
	if !ok {
		return result
	}

	objNameReplaced := strings.ReplaceAll(file.Filename, "/", "_")
	filename := fmt.Sprintf("%s_%s", file.Bucket, objNameReplaced)
	result.FullPath = fmt.Sprintf("%s/vmfiles/%s/%s", config.FileDir, data.VMUUID, filename)
	result.Type = file.Type

	return result
}

// render XML for Libvirt
func XMLTemplate(data *schema.InstanceData, templateType string) (result string, err error) {
	var (
		templateStr *string
		t           *template.Template
		buf         bytes.Buffer
	)

	// check templateType
	switch templateType {
	case "vm":
		templateStr = &data.Config.VMXML
	case "network":
		templateStr = &data.Config.NetworkXML
	default:
		return "", fmt.Errorf("invalid templateType")
	}

	// funcmap
	funcMap := template.FuncMap{
		"GetFileInfo": GetFileInfo,
	}

	// render
	if t, err = template.New(templateType).Funcs(funcMap).Parse(*templateStr); err != nil {
		return "", err
	}

	if err = t.Execute(&buf, data); err != nil {
		return "", err
	}

	result = buf.String()

	return result, nil
}
