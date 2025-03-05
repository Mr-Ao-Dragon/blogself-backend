package main

import (
	"log"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Mr-Ao-Dragon/blogself-backend/graph"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func init() {
	if gin.Mode() == gin.ReleaseMode {

	}
}
func main() {

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router := gin.Default()
	port := func() string {
		switch gin.Mode() {
		case gin.DebugMode:
			return defaultPort
		case gin.ReleaseMode:
			return strconv.Itoa(3000)
		case gin.TestMode:
			return strconv.Itoa(8080)
		default:
			return defaultPort
		}
	}
	router.Any("/", func(c *gin.Context) {

		playgroundHandler := playground.Handler("GraphQL playground", "/graphql/playground")
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	})
	router.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})
	// 包装 srv.ServeHTTP 以适配 gin.HandlerFunc
	router.Any("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port())
	log.Fatal(router.Run(":" + port()))
}
