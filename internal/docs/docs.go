package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "schemes": {{ marshal .Schemes }},
  "swagger": "2.0",
  "info": {
    "description": "{{escape .Description}}",
    "title": "{{.Title}}",
    "contact": {},
    "version": "{{.Version}}"
  },
  "host": "{{.Host}}",
  "basePath": "{{.BasePath}}",
  "paths": {
    "/links": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "links"
        ],
        "summary": "Add links batch",
        "parameters": [
          {
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/handlers.LinksRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/handlers.LinksResponse"
            }
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "/links/report": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/pdf"
        ],
        "tags": [
          "report"
        ],
        "summary": "Generate report",
        "parameters": [
          {
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/handlers.ReportRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "PDF"
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          },
          "404": {
            "description": "Not Found",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "handlers.LinksRequest": {
      "type": "object",
      "properties": {
        "links": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "handlers.LinksResponse": {
      "type": "object",
      "properties": {
        "links": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "links_num": {
          "type": "integer"
        }
      }
    },
    "handlers.ReportRequest": {
      "type": "object",
      "properties": {
        "links_list": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        }
      }
    }
  }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Link Status Checker API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
