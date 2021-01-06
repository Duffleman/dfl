package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"dfl/lib/key"
	dflmw "dfl/lib/middleware"
	"dfl/lib/ptr"
	rpclib "dfl/lib/rpc"
	"dfl/svc/auth/server/app"
	"dfl/svc/auth/server/db"
	"dfl/svc/auth/server/rpc"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq" // required for the PGSQL driver to be loaded
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const privateKey = `-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDAWpmGQ3J04zCgXdYgw/o8CIsB+9aG1b/UxKP0pU0ws4JyZ7EjXvbJo
/t+HptXPs7ugBwYFK4EEACKhZANiAAQxzjcwIr8FkpP61946t7+0CE3OLY6+sTKK
8MojiLFomEpJ2MYou+SjPc7m0ZSA9Yi26Ba5QyiHNgOo6cNVQBrNrYd47dJQ4YYp
4ojMVyng1IOaN0tSO37xrr/BjcQCrEw=
-----END EC PRIVATE KEY-----`

const publicKey = `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEMc43MCK/BZKT+tfeOre/tAhNzi2OvrEy
ivDKI4ixaJhKSdjGKLvkoz3O5tGUgPWItugWuUMohzYDqOnDVUAaza2HeO3SUOGG
KeKIzFcp4NSDmjdLUjt+8a6/wY3EAqxM
-----END PUBLIC KEY-----`

type Config struct {
	Logger *logrus.Logger

	Port int    `envconfig:"port"`
	DSN  string `envconfig:"dsn"`

	PrivateKey string `envconfig:"private_key"`
	PublicKey  string `envconfig:"public_key"`
	JWTIssuer  string `envconfig:"jwt_issuer"`

	WebAuthn WebAuthn `envconfig:"webauthn"`
}

type WebAuthn struct {
	ID          string `envconfig:"rpid"`
	Origin      string `envconfig:"rporigin"`
	DisplayName string `envconfig:"rpdisplayname"`
}

func DefaultConfig() Config {
	return Config{
		Logger: logrus.New(),
		Port:   3000,

		DSN: "postgresql://postgres@localhost/dflauth?sslmode=disable",

		PrivateKey: privateKey,
		PublicKey:  publicKey,
		JWTIssuer:  "localhost:3000",

		WebAuthn: WebAuthn{
			ID:          "localhost",
			DisplayName: "DFL Auth",
			Origin:      "http://localhost:3000",
		},
	}
}

func Run(cfg Config) error {
	cfg.Logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
	}

	pgDb, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return err
	}

	db := db.New(pgDb)

	sk := app.NewSK(
		key.MustParseECDSAPrivate([]byte(cfg.PrivateKey)),
		key.MustParseECDSAPublic([]byte(cfg.PublicKey)),
	)

	web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.WebAuthn.DisplayName,
		RPID:          cfg.WebAuthn.ID,
		RPOrigin:      cfg.WebAuthn.Origin,
	})
	if err != nil {
		return err
	}

	app := &app.App{
		Logger:    cfg.Logger,
		WA:        web,
		SK:        sk,
		DB:        db,
		JWTIssuer: cfg.JWTIssuer,
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	router.Use(dflmw.AuthMiddleware(sk.Public(), []dflmw.HTTPResource{
		{Verb: ptr.String(http.MethodGet), Path: ptr.String("/")},
		{Verb: ptr.String(http.MethodGet), Path: ptr.String("/authorize")},
		{Verb: ptr.String(http.MethodGet), Path: ptr.String("/register")},
		{Verb: ptr.String(http.MethodGet), Path: ptr.String("/u2f_manage")},
		{Verb: ptr.String(http.MethodPost), Path: ptr.String("/authorize_confirm")},
		{Verb: ptr.String(http.MethodPost), Path: ptr.String("/authorize_prompt")},
		{Verb: ptr.String(http.MethodPost), Path: ptr.String("/register_confirm")},
		{Verb: ptr.String(http.MethodPost), Path: ptr.String("/register_prompt")},
		{Verb: ptr.String(http.MethodPost), Path: ptr.String("/token")},
	}))

	// Internal
	router.Get("/", wrap(app, rpc.Index))
	router.Get("/register", wrap(app, rpc.RegisterGet))
	router.Get("/u2f_manage", wrap(app, rpc.U2FManageGet))
	router.Post("/create_client", wrap(app, rpc.CreateClient))
	router.Post("/create_invite_code", wrap(app, rpc.CreateInviteCode))
	router.Post("/create_key_confirm", wrap(app, rpc.CreateKeyConfirm))
	router.Post("/create_key_prompt", wrap(app, rpc.CreateKeyPrompt))
	router.Post("/delete_key", wrap(app, rpc.DeleteKey))
	router.Post("/list_u2f_keys", wrap(app, rpc.ListU2FKeys))
	router.Post("/register_confirm", wrap(app, rpc.RegisterConfirm))
	router.Post("/register_prompt", wrap(app, rpc.RegisterPrompt))
	router.Post("/sign_key_confirm", wrap(app, rpc.SignKeyConfirm))
	router.Post("/sign_key_prompt", wrap(app, rpc.SignKeyPrompt))
	router.Post("/whoami", wrap(app, rpc.WhoAmI))

	// OAuth
	router.Get("/authorize", wrap(app, rpc.AuthorizeGet))
	router.Post("/authorize_confirm", wrap(app, rpc.AuthorizeConfirm))
	router.Post("/authorize_prompt", wrap(app, rpc.AuthorizePrompt))
	router.Post("/token", wrap(app, rpc.Token))

	cfg.Logger.Infof("Server running on port %d", cfg.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}

func wrap(a *app.App, fn func(*app.App, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(a, w, r)
		if err != nil {
			rpclib.HandleError(w, r, err, a.Logger)
			return
		}
	}
}
