// Package docs contains swagger documentation
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "PolyGo API",
    "description": "High-performance Polymarket API proxy with caching and WebSocket support",
    "version": "1.0",
    "contact": {
      "name": "API Support",
      "email": "support@polygo.dev"
    },
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    }
  },
  "host": "{{.Host}}",
  "basePath": "{{.BasePath}}",
  "schemes": ["http", "https"]
}`

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "PolyGo API",
	Description:      "High-performance Polymarket API proxy",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
