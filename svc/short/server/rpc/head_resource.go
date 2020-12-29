package rpc

import (
	"net/http"

	"dfl/svc/short/server/app"
	"dfl/svc/short/server/lib/fakehttp"
)

// HeadResource gets a resource and handles the response for it
func HeadResource(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fake := fakehttp.NewResponse()

		HandleResource(a)(fake, r)

		for key, value := range fake.Headers {
			for _, v := range value {
				w.Header().Add(key, v)
			}
		}

		if fake.Status >= 100 && fake.Status <= 999 {
			w.WriteHeader(fake.Status)
		}

		return
	}
}
