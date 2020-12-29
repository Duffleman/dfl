package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"dfl/lib/cache"
	"dfl/lib/cher"
	"dfl/lib/config"
	"dfl/svc/short/server/app"
	"dfl/svc/short/server/db"
	dflmw "dfl/svc/short/server/lib/middleware"
	"dfl/svc/short/server/lib/storageproviders"
	"dfl/svc/short/server/rpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq" // required for the PGSQL driver to be loaded
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/speps/go-hashids"
)

type Config struct {
	Logger *logrus.Logger

	Port int    `json:"port"`
	DSN  string `json:"dsn"`

	Salt    string            `json:"salt"`
	RootURL string            `json:"root_url"`
	Users   map[string]string `json:"users"`

	StorageProvider string `json:"storage_provider"`

	LFSFolder      string      `json:"lfs_folder"`
	LFSPermissions os.FileMode `json:"lfs_permissions"`

	AWS       *aws.Config `json:"aws"`
	AWSBucket string      `json:"aws_bucket"`
	AWSRoot   string      `json:"aws_root"`

	Redis config.Redis `json:"redis"`
}

func DefaultConfig() Config {
	return Config{
		Logger: logrus.New(),

		Port: 3000,
		DSN:  "postgresql://postgres@localhost/dflimg?sslmode=disable",

		Salt:    "savour-shingle-sidney-rajah-punk-lead-jenny-scot",
		RootURL: "http://localhost:3000",
		Users: map[string]string{
			"Duffleman": "test",
		},

		StorageProvider: "lfs",

		LFSFolder:      "/Users/duffleman/Downloads/short",
		LFSPermissions: 0644,

		AWS: &aws.Config{
			Region: aws.String("eu-west-1"),
		},
		AWSBucket: "s3.duffleman.co.uk",
		AWSRoot:   "i.dfl.mn",

		Redis: config.Redis{
			URI: "redis://localhost:6379",
		},
	}
}

func Run(cfg Config) error {
	cfg.Logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
	}

	var err error
	var sp storageproviders.StorageProvider

	switch cfg.StorageProvider {
	case "aws":
		sp, err = storageproviders.NewAWSProviderFromCfg(cfg.AWS, cfg.AWSBucket, cfg.AWSRoot)
		if err != nil {
			return err
		}
	case "lfs":
		sp, err = storageproviders.NewLFSProviderFromCfg(cfg.LFSFolder, cfg.LFSPermissions)
		if err != nil {
			return err
		}
	default:
		return cher.New("unsupported_provider", nil)
	}

	pgDb, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return err
	}

	db := db.New(pgDb)

	hd := hashids.NewData()
	hd.Salt = cfg.Salt
	hd.MinLength = 3

	hasher, err := hashids.NewWithData(hd)
	if err != nil {
		return err
	}

	redisOpts, err := redis.ParseURL(cfg.Redis.URI)
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(redisOpts)
	_, err = redisClient.Ping().Result()
	if err != nil {
		return err
	}

	redis := cache.NewCache(redisClient)

	app := &app.App{
		SP:      sp,
		DB:      db,
		Hasher:  hasher,
		Redis:   redis,
		RootURL: cfg.RootURL,
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	router.Use(dflmw.AuthMiddleware(cfg.Users))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "") // Needed for redirect to work
		http.Redirect(w, r, "https://duffleman.co.uk", http.StatusMovedPermanently)
	})

	router.Get("/robots.txt", rpc.Robots())

	router.Get("/resave_hashes", rpc.ResaveHashes(app))
	router.Post("/add_shortcut", rpc.AddShortcut(app))
	router.Post("/create_signed_url", rpc.CreateSignedURL(app))
	router.Post("/delete_resource", rpc.DeleteResource(app))
	router.Post("/list_resources", rpc.ListResources(app))
	router.Post("/remove_shortcut", rpc.RemoveShortcut(app))
	router.Post("/set_nsfw", rpc.SetNSFW(app))
	router.Post("/shorten_url", rpc.ShortenURL(app))
	router.Post("/upload_file", rpc.UploadFile(app))
	router.Post("/view_details", rpc.ViewDetails(app))

	router.Get("/{query}", rpc.HandleResource(app))
	router.Head("/{query}", rpc.HeadResource(app))

	cfg.Logger.Infof("Server running on port %d", cfg.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}
