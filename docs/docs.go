// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.github.com/avealice"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "paths": {
        "/api/actor": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Создает нового актера.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Создать актера.",
                "parameters": [
                    {
                        "description": "Данные нового актера",
                        "name": "actor",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.Actor"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Актер успешно создан",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/actor/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получает информацию об актере по его идентификатору.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Получить информацию об актере.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор актера",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Actor"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Обновляет информацию об актере по его идентификатору.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Обновить информацию об актере.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор актера",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Новые данные актера",
                        "name": "actor",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.ActorWithMovies"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Информация об актере успешно обновлена",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Удаляет актера по его идентификатору.",
                "tags": [
                    "actors"
                ],
                "summary": "Удалить актера.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор актера",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Актер успешно удален",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/actors": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получить всех актеров из базы данных.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Получить всех актеров.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handler.Actor"
                            }
                        }
                    }
                }
            }
        },
        "/api/movie": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Выполняет поиск фильмов по указанным критериям (название или актер).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Поиск фильмов",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Название фильма для поиска",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Имя актера для поиска",
                        "name": "actor",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список фильмов, удовлетворяющих критериям поиска",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.MovieWithActors"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Создает новый фильм.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Создать фильм",
                "parameters": [
                    {
                        "description": "Данные нового фильма",
                        "name": "movie",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.MovieWithActors"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Фильм создан успешно"
                    }
                }
            }
        },
        "/api/movie/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получает информацию о фильме по его идентификатору.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Получить информацию о фильме",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор фильма",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Информация о фильме",
                        "schema": {
                            "$ref": "#/definitions/model.MovieWithActors"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Обновляет информацию о фильме.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Обновить информацию о фильме",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор фильма",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Новые данные о фильме",
                        "name": "movie",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.MovieWithActors"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Фильм обновлен успешно"
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Удаляет фильм по его идентификатору.",
                "tags": [
                    "Movies"
                ],
                "summary": "Удалить фильм",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор фильма",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Фильм удален успешно"
                    }
                }
            }
        },
        "/api/movies": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получает список всех фильмов с возможностью сортировки.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Получить все фильмы",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Критерий сортировки (например, 'rating')",
                        "name": "sort_by",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Порядок сортировки (например, 'asc' или 'desc')",
                        "name": "sort_order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список фильмов",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Movie"
                            }
                        }
                    }
                }
            }
        },
        "/auth/sign-in": {
            "post": {
                "description": "Авторизует пользователя с заданными учетными данными и возвращает токен доступа",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Авторизация"
                ],
                "summary": "Авторизация пользователя",
                "parameters": [
                    {
                        "description": "Данные для входа",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SignInInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "map[token]: Токен доступа",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос или данные",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResponse"
                        }
                    }
                }
            }
        },
        "/auth/sign-up": {
            "post": {
                "description": "Регистрирует нового пользователя с заданными данными",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Регистрация"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "Данные нового пользователя",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "map[id]: ID нового пользователя",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос или данные",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.Actor": {
            "type": "object",
            "properties": {
                "birth_date": {
                    "description": "Birth date of the actor",
                    "type": "string"
                },
                "gender": {
                    "description": "Gender of the actor",
                    "type": "string"
                },
                "name": {
                    "description": "Name of the actor",
                    "type": "string"
                }
            }
        },
        "handler.ActorWithMovies": {
            "type": "object",
            "properties": {
                "birth_date": {
                    "description": "Birth date of the actor",
                    "type": "string"
                },
                "gender": {
                    "description": "Gender of the actor",
                    "type": "string"
                },
                "movies": {
                    "description": "Movies associated with the actor",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.Movie"
                    }
                },
                "name": {
                    "description": "Name of the actor",
                    "type": "string"
                }
            }
        },
        "handler.Movie": {
            "type": "object",
            "properties": {
                "description": {
                    "description": "Description of the movie",
                    "type": "string"
                },
                "rating": {
                    "description": "Rating of the movie",
                    "type": "integer"
                },
                "release_date": {
                    "description": "Release date of the movie",
                    "type": "string"
                },
                "title": {
                    "description": "Title of the movie",
                    "type": "string"
                }
            }
        },
        "handler.SignInInput": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handler.errorResponse": {
            "description": "JSON-структура ответа с сообщением об ошибке.",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "model.Actor": {
            "type": "object",
            "properties": {
                "birth_date": {
                    "description": "Birth date of the actor",
                    "type": "string"
                },
                "gender": {
                    "description": "Gender of the actor",
                    "type": "string"
                },
                "name": {
                    "description": "Name of the actor",
                    "type": "string"
                }
            }
        },
        "model.Movie": {
            "type": "object",
            "properties": {
                "description": {
                    "description": "Description of the movie",
                    "type": "string"
                },
                "rating": {
                    "description": "Rating of the movie",
                    "type": "integer"
                },
                "release_date": {
                    "description": "Release date of the movie",
                    "type": "string"
                },
                "title": {
                    "description": "Title of the movie",
                    "type": "string"
                }
            }
        },
        "model.MovieWithActors": {
            "type": "object",
            "properties": {
                "actors": {
                    "description": "Actors associated with the movie",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Actor"
                    }
                },
                "description": {
                    "description": "Description of the movie",
                    "type": "string"
                },
                "rating": {
                    "description": "Rating of the movie",
                    "type": "integer"
                },
                "release_date": {
                    "description": "Release date of the movie",
                    "type": "string"
                },
                "title": {
                    "description": "Title of the movie",
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "password": {
                    "description": "Password hash of the user",
                    "type": "string"
                },
                "role": {
                    "description": "Role of the user",
                    "type": "string"
                },
                "username": {
                    "description": "Username of the user",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8000",
	BasePath:         "/v2",
	Schemes:          []string{},
	Title:            "FilmHub API",
	Description:      "API для работы с фильмами и актерами в FilmHub.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
    
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
