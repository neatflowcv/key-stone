package design

import (
	. "goa.design/goa/v3/dsl" //nolint:staticcheck
)

var _ = API("key-stone", func() {
	Title("Key Stone")
	Description("Key Stone is a platform for managing your keys")
	Version("v0.0.1")

	HTTP(func() {
		Path("/key-stone")
	})
})

var _ = Service("user", func() {
	HTTP(func() {
		Path("/users")
	})

	Error("BadRequest", ErrorResult, "Bad Request")
	Error("Unauthorized", ErrorResult, "Unauthorized")

	Method("create", func() {
		Payload(func() {
			Attribute("user", UserInput)

			Required("user")
		})

		HTTP(func() {
			POST("/")

			Response(StatusNoContent)
			Response("BadRequest", StatusBadRequest)
		})
	})
	Method("delete", func() {
		Payload(DeleteUserPayload)

		HTTP(func() {
			DELETE("/me")

			Header("Authorization", String, "The authorization header")

			Response(StatusNoContent)
			Response("Unauthorized", StatusUnauthorized)
		})
	})
})

var _ = Service("token", func() {
	HTTP(func() {
		Path("/auth")
	})

	Error("BadRequest", ErrorResult, "Bad Request")
	Error("Unauthorized", ErrorResult, "Unauthorized")

	Method("issue", func() {
		Payload(func() {
			Attribute("user", IssueInput)

			Required("user")
		})

		Result(TokenDetail)

		HTTP(func() {
			POST("/")

			Response(StatusOK)
			Response("BadRequest", StatusBadRequest)
			Response("Unauthorized", StatusUnauthorized)
		})
	})

	Method("refresh", func() {
		Payload(func() {
			Attribute("token", RefreshInput)

			Required("token")
		})

		Result(TokenDetail)

		HTTP(func() {
			POST("/refresh")

			Response(StatusOK)
			Response("BadRequest", StatusBadRequest)
			Response("Unauthorized", StatusUnauthorized)
		})
	})
})

var UserInput = Type("UserInput", func() { //nolint:gochecknoglobals
	Attribute("name", String, "The name of the user")
	Attribute("password", String, "The password of the user")
	Attribute("payload", MapOf(String, Any), "The payload of the user")

	Required("name", "password")
})

var DeleteUserPayload = Type("DeleteUserPayload", func() { //nolint:gochecknoglobals
	Attribute("Authorization", String, "The payload of the user")

	Required("Authorization")
})

var IssueInput = Type("IssueInput", func() { //nolint:gochecknoglobals
	Attribute("username", String, "The username of the user")
	Attribute("password", String, "The password of the user")

	Required("username", "password")
})

var TokenDetail = Type("TokenDetail", func() { //nolint:gochecknoglobals
	Attribute("access_token", String, "The access token of the user")
	Attribute("token_type", String, "The token type of the user")
	Attribute("expires_in", Int, "The expires in of the user")
	Attribute("refresh_token", String, "The refresh token of the user")

	Required("access_token", "token_type", "expires_in", "refresh_token")
})

var RefreshInput = Type("RefreshInput", func() { //nolint:gochecknoglobals
	Attribute("access_token", String, "The access token of the user")
	Attribute("refresh_token", String, "The refresh token of the user")

	Required("access_token", "refresh_token")
})
