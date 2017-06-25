package web

import (
	"net/http"

	"github.com/if1live/kanae/histories/lendings"
)

func handlerLending(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Sync *lendings.Sync
		View *lendings.View
	}

	sync := svr.db.MakeLendingSync(nil)
	view := svr.db.MakeLendingView()

	ctx := Context{
		Sync: sync,
		View: view,
	}
	err := renderLayoutTemplate(w, "layout.html", "lending.html", ctx)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}
