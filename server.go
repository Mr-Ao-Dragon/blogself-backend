package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Mr-Ao-Dragon/blogself-backend/graph"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

var dataToken = cache.New(time.Hour, 24*time.Hour)

func init() {
	if gin.Mode() == gin.ReleaseMode {
		if err := dataToken.Add("AK", os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), cache.DefaultExpiration); err != nil {
			log.Panicln("Get Access key fail, crash now.")
		}
		if err := dataToken.Add("SK", os.Getenv("ALIBABA_CLOUD_SECRET_KEY_ID"), cache.DefaultExpiration); err != nil {
			log.Panicln("Get Secret key fail, crash now.")
		}
		if err := dataToken.Add("STS", os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN"), cache.DefaultExpiration); err != nil {
			log.Panicln("Get securityToken fail, crash now.")
		}
		if err := dataToken.Add("TARGET_DB", os.Getenv("OTS_NAME"), cache.DefaultExpiration); err != nil {
			log.Panicln("Get database name fail, crash now.")
		}
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

		playgroundHandler := playground.Handler("GraphQL playground", "/graphql")
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
