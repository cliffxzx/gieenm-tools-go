package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/cors"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	gql "github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql/resolvers"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/seeder"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/validator"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/gin-contrib/gzip"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func migrationReset() error {
	driver, err := postgres.WithInstance(database.GetDB().DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://pkg/gieenm-system/migrations", "gieenm_tools", driver)
	if err != nil {
		return err
	}

	version, _, _ := m.Version()
	if err := m.Drop(); err != nil && version != 0 {
		return err
	}

	database.GetDB().Exec("create table if not exists schema_migrations (version int8 not null primary key, dirty bool not null);")

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		log.Println(err)
	}

	return nil
}

func main() {
	//Start the default gin server
	r := gin.Default()

	//Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please create one in the root directory")
	}

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(utils.StoreGinContextToSelf)

	r.Use(cors.CORS)
	r.Use(validator.Validator)

	//Start PostgreSQL database
	database.Init()

	// Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	database.InitRedis("1")

	err = firewall.InitFirewalls()
	if err != nil {
		panic(err)
	}

	r.POST("/graphql", gin.WrapH(handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &resolvers.Resolver{}}))))

	port := utils.MustGetEnv("PORT")

	switch strings.ToUpper(utils.MustGetEnv("ENV")) {
	case "PRODUCTION":
		gin.SetMode(gin.ReleaseMode)
	case "DEVELOPMENT":
		r.GET("/graphql", gin.WrapH(playground.Handler("GraphQL", "/graphql")))

		if err := migrationReset(); err != nil {
			panic(err)
		}

		seeder.Init()
		if err := seeder.Seeder(); err != nil {
			panic(err)
		}
	default:
		panic("undefined environment mode")
	}

	r.Run(fmt.Sprintf(":%s", port))
}
