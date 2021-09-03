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

	fieldKey := "LOOKUP_link"

	Module = &sonic.Module{
		Name: name,
		Fields: []*sonic.ModuleField{
			{
				Validation: "https://www.target.*?",
				Type:       sonic.FieldTypeText,
				Label:      "Product Link",
				FieldKey:   &fieldKey,
			},
		},
	}
}
