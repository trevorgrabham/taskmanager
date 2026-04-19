package routes

import (
	"net/http"
)

// Hits the backend db.DeleteTask() and redirects to "/"
func (h Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok bool
		id int
	)
	if ok = h.checkConnection(w); !ok {
		return
	}

	if id, ok = h.parseID(w, r.URL.Query().Get("id")); !ok {
		return
	}

	if ok = h.deleteTask(w, id); !ok {
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
