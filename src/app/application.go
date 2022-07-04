package app

import (
	"bookstore-oauth-api/clients/cassandra"
	"bookstore-oauth-api/domain/access_token"
	"bookstore-oauth-api/http"
	"bookstore-oauth-api/repository/db"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session := cassandra.GetSession()
	defer session.Close()

	dbRepository := db.NewRepository()
	accessTokenService := access_token.NewService(dbRepository)
	accessTokenHandler := http.NewHandler(accessTokenService)

	router.GET("/oauth/access_token/:access_token_id", accessTokenHandler.GetById)
	router.POST("/oauth/access_token", accessTokenHandler.Create)

	router.Run(":8080")
}
