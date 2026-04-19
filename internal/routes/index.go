package routes

import (
	"local/taskmanager/internal/views"
	"local/taskmanager/internal/views/dashboard"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
	"time"
)

// Renders the index.html page
func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok   bool
		info = views.PageInfo{StartDay: time.Now()}
	)

	if ok = h.checkKnownPath(w, r); !ok {
		return
	}

	if ok = h.checkConnection(w); !ok {
		return
	}

	if info.WeekOfTasks, ok = h.getUpcomingWeek(w); !ok {
		return
	}

	info.FavCategoryTasks = h.getFavCategory(info.WeekOfTasks)

	if info.OverdueTasks, ok = h.getOverdueTasks(w); !ok {
		return
	}

	if info.UnscheduledTasks, ok = h.getUnscheduledTasks(w); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := shared.Layout(dashboard.Index(info)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("index: %s\n", err)
		return
	}
}
