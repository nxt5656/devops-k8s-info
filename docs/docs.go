package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "Zhiyi K8s Info API",
    "description": "API service for reading Kubernetes resources.",
    "version": "1.0"
  },
  "basePath": "{{.BasePath}}",
  "paths": {}
}`

// SwaggerInfo keeps global Swagger metadata.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1/k8s-info",
	Schemes:          []string{},
	Title:            "Zhiyi K8s Info API",
	Description:      "API service for reading Kubernetes resources.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

type swaggerDoc struct{}

func (s *swaggerDoc) ReadDoc() string {
	return SwaggerInfo.ReadDoc()
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), &swaggerDoc{})
}
