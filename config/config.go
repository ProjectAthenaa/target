package config

import (
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"os"
	"strings"
)

var Module *sonic.Module

func init() {
	var name = "target"

	if podName := os.Getenv("POD_NAME"); podName != "" {
		name = strings.Split(podName, "-")[0]
	}

	fieldKey := "TARGET_link"

	Module = &sonic.Module{
		Name: name,
		Fields: []*sonic.ModuleField{
			{
				Validation: "https:\\/\\/www\\.target\\.com\\/p\\/.*\\/-\\/A-\\d+",
				Type:       sonic.FieldTypeText,
				Label:      "Product Link",
				FieldKey:   &fieldKey,
			},
		},
	}
}
