package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/cors"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/sync"
	gql "github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql/resolvers"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql/scalars"
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

		subnetRaw := map[string]string{
			// 109
			"20180205214124,3": "192.168.1.0/24",
			// 114
			"20180203234922,1": "192.168.2.0/24",
			// 115-1,115-2
			"20180204001228,2": "192.168.3.0/24",
			// 1f 管制群組
			"20190925091129,4": "192.168.0.0/16",
			// 206
			"20170629102323,1": "192.168.1.0/24",
			// 208
			"20170629102325,2": "192.168.2.0/24",
			// 311A
			"20170629102327,3": "192.168.3.0/24",
			// 2f 管制群組
			"20200427155607,4": "192.168.0.0/16",
			// 305-306
			"20170629105419,1": "192.168.1.0/24",
			// 307-310
			"20170629114059,2": "192.168.2.0/24",
			// 405
			"20170629200244,3": "192.168.3.0/24",
			// 405-1
			"20180202235136,2": "192.168.1.0/24",
			// 407-408
			"20170629110845,3": "192.168.2.0/24",
			// 410
			"20180202225649,1": "192.168.3.0/24",
			// 4f 管制群組
			"20190615104456,1": "192.168.0.0/16",
			// gieenm-tools-tokenStr
			"20210319014124,2": "192.168.2.0/16",
		}

		defaultSubnets := map[string]*scalars.IPAddr{}
		for key, raw := range subnetRaw {
			_, subnet, _ := net.ParseCIDR(raw)
			tmp := scalars.IPAddr(*subnet)
			defaultSubnets[key] = &tmp
		}

		err := sync.SyncNusoftToDatabase(defaultSubnets)
		if err != nil {
			panic(err)
		}
	default:
		panic("undefined environment mode")
	}

	r.Run(fmt.Sprintf(":%s", port))
}
