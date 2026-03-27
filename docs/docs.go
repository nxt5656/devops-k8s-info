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
  "paths": {
    "/namespaces": {
      "get": {
        "description": "Get all namespaces from the current Kubernetes cluster.",
        "produces": [
          "application/json"
        ],
        "tags": [
          "namespaces"
        ],
        "summary": "List namespaces",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "type": "object"
            }
          },
          "502": {
            "description": "Bad Gateway",
            "schema": {
              "type": "object"
            }
          }
        }
      }
    }
  }
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
