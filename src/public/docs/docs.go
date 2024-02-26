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
        "/api/v1/infer": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "InferenceController"
                ],
                "summary": "Filter inferences",
                "operationId": "filter-inference",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task IDs",
                        "name": "ids",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/resources.Inference"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "InferenceController"
                ],
                "summary": "Create an inference task",
                "operationId": "create-inference",
                "parameters": [
                    {
                        "type": "string",
                        "default": "animate_lcm",
                        "description": "AI model ID",
                        "name": "model",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "enum": [
                            "t2v",
                            "i2v",
                            "v2v",
                            "upscale"
                        ],
                        "type": "string",
                        "default": "t2v",
                        "description": "Infer type",
                        "name": "type",
                        "in": "formData"
                    },
                    {
                        "maxLength": 1000,
                        "minLength": 1,
                        "type": "string",
                        "description": "The prompt or prompts to guide image generation.",
                        "name": "prompt",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The prompt or prompts to guide what to not include in image generation.",
                        "name": "negative_prompt",
                        "in": "formData"
                    },
                    {
                        "maximum": 25,
                        "minimum": 1,
                        "type": "integer",
                        "default": 25,
                        "description": "More steps usually lead to a higher quality image at the expense of slower inference",
                        "name": "num_inference_steps",
                        "in": "formData"
                    },
                    {
                        "maximum": 16,
                        "minimum": 16,
                        "type": "integer",
                        "default": 16,
                        "description": "The number of video frames to generate. Default FPS: 8",
                        "name": "num_frames",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "default": 512,
                        "description": "The width in pixels of the generated image/video.",
                        "name": "width",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "default": 512,
                        "description": "The height in pixels of the generated image/video.",
                        "name": "height",
                        "in": "formData"
                    },
                    {
                        "maximum": 100,
                        "minimum": 0,
                        "type": "number",
                        "default": 2,
                        "description": "A higher guidance scale value encourages the model to generate images closely linked to the ` + "`" + `prompt` + "`" + ` at the expense of lower image quality.",
                        "name": "guidance_scale",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/resources.Inference"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/infer/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "InferenceController"
                ],
                "summary": "Get status of an inference task",
                "operationId": "get-inference",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/resources.Inference"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/models": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ModelController"
                ],
                "summary": "Get list supported models",
                "operationId": "get-models",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/resources.Model"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "resources.Inference": {
            "type": "object",
            "properties": {
                "deadline": {
                    "description": "Deadline for completing the task",
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_err": {
                    "type": "string"
                },
                "last_failed_at": {
                    "type": "integer"
                },
                "max_retry": {
                    "type": "integer"
                },
                "model": {
                    "type": "string"
                },
                "retried": {
                    "type": "integer"
                },
                "status": {
                    "description": "Status of the task. Values: active, pending, scheduled, retry, archived, completed",
                    "type": "string"
                },
                "target_file_url": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "resources.Model": {
            "type": "object",
            "properties": {
                "guidance_scale": {
                    "type": "number"
                },
                "height": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "logo_url": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "num_frames": {
                    "type": "integer"
                },
                "num_inference_steps": {
                    "type": "integer"
                },
                "recomd_prompts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resources.RecomdPrompt"
                    }
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "resources.RecomdPrompt": {
            "type": "object",
            "properties": {
                "asset_url": {
                    "type": "string"
                },
                "neg_prompt": {
                    "type": "string"
                },
                "prompt": {
                    "type": "string"
                },
                "thumbnail_url": {
                    "type": "string"
                }
            }
        },
        "response.Meta": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "meta": {
                    "$ref": "#/definitions/response.Meta"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "GoghAI API Public",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
