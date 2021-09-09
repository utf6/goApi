// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2021-09-08 11:40:35.0663479 +0800 CST m=+35.725264701

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/getToken": {
            "post": {
                "tags": [
                    "权限管理"
                ],
                "summary": "获取token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/api/v1/articles": {
            "get": {
                "tags": [
                    "文章管理"
                ],
                "summary": "获取多个文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "标签id",
                        "name": "tag_id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "状态（0：删除，1：正常）",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "access_token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "文章管理"
                ],
                "summary": "新增文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "标签id",
                        "name": "tag_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "文章标题",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "文章描述",
                        "name": "desc",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "文章内容",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "access_token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/api/v1/articles/{id}": {
            "get": {
                "tags": [
                    "文章管理"
                ],
                "summary": "获取单个文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章id",
                        "name": "id",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "access_token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "put": {
                "tags": [
                    "文章管理"
                ],
                "summary": "修改文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "标签id",
                        "name": "tag_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "文章标题",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "文章描述",
                        "name": "desc",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "文章内容",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "access_token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "文章管理"
                ],
                "summary": "删除文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "access_token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/api/v1/tags": {
            "get": {
                "tags": [
                    "标签管理"
                ],
                "summary": "获取文章标签",
                "parameters": [
                    {
                        "type": "string",
                        "description": "标签名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "状态（0：禁用，1：正常）",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "access_token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "标签管理"
                ],
                "summary": "新增文章标签",
                "parameters": [
                    {
                        "type": "string",
                        "description": "标签名称",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "access_token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200, \"data\":{}, \"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/api/v1/tags/{id}": {
            "put": {
                "tags": [
                    "标签管理"
                ],
                "summary": "修改文章标签",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "标签名称",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "标签管理"
                ],
                "summary": "删除文章标签",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "gin.H": {
            "type": "object",
            "additionalProperties": true
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
