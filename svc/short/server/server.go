package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"dfl/lib/cache"
	"dfl/lib/cher"
	"dfl/lib/key"
	dflmw "dfl/lib/middleware"
	"dfl/lib/ptr"
	"dfl/svc/short/server/app"
	"dfl/svc/short/server/db"
	"dfl/svc/short/server/lib/storageproviders"
	"dfl/svc/short/server/rpc"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq" // required for the PGSQL driver to be loaded
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/speps/go-hashids"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEMc43MCK/BZKT+tfeOre/tAhNzi2OvrEy
ivDKI4ixaJhKSdjGKLvkoz3O5tGUgPWItugWuUMohzYDqOnDVUAaza2HeO3SUOGG
KeKIzFcp4NSDmjdLUjt+8a6/wY3EAqxM
-----END PUBLIC KEY-----`

type Config struct {
	Logger *logrus.Logger

	Port      int    `envconfig:"port"`
	DSN       string `envconfig:"dsn"`
	PublicKey string `envconfig:"public_key"`

	Salt    string `envconfig:"salt"`
	RootURL string `envconfig:"root_url"`

	StorageProvider string `envconfig:"storage_provider"`

	LFSFolder      string      `envconfig:"lfs_folder"`
	LFSPermissions os.FileMode `envconfig:"lfs_permissions"`

	AWSRegion string `envconfig:"aws_region"`
	AWSBucket string `envconfig:"aws_bucket"`
	AWSRoot   string `envconfig:"aws_root"`

	RedisURI string `envconfig:"redis_uri"`
}

func DefaultConfig() Config {
	return Config{
		Logger: logrus.New(),

		Port:      3000,
		DSN:       "postgresql://postgres@localhost/dflimg?sslmode=disable",
		PublicKey: publicKey,

		Salt:    "savour-shingle-sidney-rajah-punk-lead-jenny-scot",
		RootURL: "http://localhost:3000",

		StorageProvider: "lfs",

		LFSFolder:      "/Users/duffleman/Downloads/short",
		LFSPermissions: 0644,

		AWSRegion: "eu-west-1",
		AWSBucket: "s3.duffleman.co.uk",
		AWSRoot:   "i.dfl.mn",

		RedisURI: "redis://localhost:6379",
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
		sp, err = storageproviders.NewAWSProviderFromCfg(cfg.AWSRegion, cfg.AWSBucket, cfg.AWSRoot)
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

	public := key.MustParseECDSAPublic([]byte(cfg.PublicKey))

	db := db.New(pgDb)

	hd := hashids.NewData()
	hd.Salt = cfg.Salt
	hd.MinLength = 3

	hasher, err := hashids.NewWithData(hd)
	if err != nil {
		return err
	}

	redisOpts, err := redis.ParseURL(cfg.RedisURI)
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
	router.Use(dflmw.AuthMiddleware(public, []dflmw.HTTPResource{
		{Verb: ptr.String(http.MethodGet)},
		{Verb: ptr.String(http.MethodHead)},
		{Verb: ptr.String(http.MethodOptions)},
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "") // Needed for redirect to work
		http.Redirect(w, r, "https://duffleman.co.uk", http.StatusMovedPermanently)
	})

	router.Get("/robots.txt", rpc.Robots())

	router.Post("/resave_hashes", rpc.ResaveHashes(app))
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
