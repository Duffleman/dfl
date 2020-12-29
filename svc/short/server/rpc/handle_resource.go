package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/short/server/app"
	"dfl/svc/short/server/rpc/pipeline"

	"github.com/go-chi/chi"
)

func HandleResource(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		qi := a.ParseQueryType(chi.URLParam(r, "query"))

		pipe := pipeline.New(ctx, &pipeline.Creator{
			App: a,
			R:   r,
			W:   w,
			QI:  qi,
		})

		pipe.Steps([]pipeline.HandlerType{
			// ordered
			pipeline.ListFromDB,
			pipeline.HandleDeleted,
			pipeline.MakeContext,
			pipeline.ValidateRequest,
			pipeline.HandleNSFWPrimer,
			pipeline.HandleURLRedirect,
			pipeline.LoadFilesFromFS,
			pipeline.SyntaxHighlighter,
			pipeline.RenderMD,
			pipeline.MakeTarball,
			pipeline.FilterMultiFile,
			pipeline.ServeContent,
		})

		rpc.HandleError(w, r, pipe.Run())
	}
}
