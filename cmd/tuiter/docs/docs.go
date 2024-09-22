// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "email": "madepietro@unlam.edu.ar"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/likes": {
            "post": {
                "description": "Add a like to a tuit",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Add a like to a tuit",
                "parameters": [
                    {
                        "description": "Like",
                        "name": "like",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.like"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.userPostPayload"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove a like from a tuit",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Remove a like from a tuit",
                "parameters": [
                    {
                        "description": "Like",
                        "name": "like",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.like"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.userPostPayload"
                        }
                    }
                }
            }
        },
        "/tuits": {
            "get": {
                "description": "Search tuits",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "posts"
                ],
                "summary": "Search tuits",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Page ID",
                        "name": "page_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.tuitPayload"
                            }
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "create a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "create a new user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.userCreatePayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/handlers.userPayload"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Get a user by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.userPayload"
                        }
                    }
                }
            }
        },
        "/users/{id}/tuits": {
            "get": {
                "description": "Search Users Tuits will return a list of tuits from the user perspective. This means that the user will",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tuits"
                ],
                "summary": "Search Users' tuits",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/userpost.Feed"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.like": {
            "type": "object",
            "properties": {
                "tuit_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "handlers.tuitPayload": {
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/handlers.userPayload"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "likes": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "parent_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "handlers.userCreatePayload": {
            "type": "object",
            "required": [
                "email",
                "name"
            ],
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handlers.userPayload": {
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handlers.userPostPayload": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "avatar_url": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "liked": {
                    "type": "boolean"
                },
                "likes": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "parent_id": {
                    "type": "integer"
                }
            }
        },
        "userpost.Feed": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "avatar_url": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "liked": {
                    "type": "boolean"
                },
                "likes": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "parent_id": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{"https"},
	Title:            "Tuiter API",
	Description:      "This is the API for Tuiter, a Twitter clone.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
