// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/videos": {
            "post": {
                "description": "A POST function which generates video based on video sources and themes (background, cover and music) selected. All params should be absolute path",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Accept user-provided videos, images and themes, and generate video for user to download",
                "parameters": [
                    {
                        "description": "GenerateVideoBody",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/videos.GenerateVideoBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/videos/cover": {
            "post": {
                "description": "A POST function which generates cover videos based on user input, e.g. font, size, styles etc.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Generates cover videos based on user input, e.g. font, size, styles etc.",
                "parameters": [
                    {
                        "description": "GenerateCoverPageBody",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/videos.GenerateCoverPageBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad requests",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "videos.GenerateCoverPageBody": {
            "type": "object",
            "properties": {
                "coverPage": {
                    "type": "string"
                },
                "destPath": {
                    "type": "string"
                },
                "endTime": {
                    "type": "string"
                },
                "fadeInDuration": {
                    "type": "string"
                },
                "fontColor": {
                    "type": "string"
                },
                "fontSize": {
                    "type": "string"
                },
                "startTime": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "x": {
                    "type": "string"
                },
                "y": {
                    "type": "string"
                }
            }
        },
        "videos.GenerateVideoBody": {
            "type": "object",
            "properties": {
                "bgmMusic": {
                    "type": "string"
                },
                "coverPage": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "videoSrcList": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
