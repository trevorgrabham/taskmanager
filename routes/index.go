package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"local/taskmanager/views/dashboard"
	"net/http"
	"time"
)

// Renders the index.html page
func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller         = "IndexHandler"
		err            error
		sessionID      string
		user           services.User
		dashboardTasks services.Dashboard
		dashboardData  views.DashboardViewData
	)

	if r.URL.Path != "/" {
		h.HandleError(w, NewUnknownPathError(caller))
		return
	}

	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, NewSessionIDCookieParseError(caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, NewUnauthenticatedError(caller))
		return
	}

	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	if dashboardTasks, err = h.services.GetDashboardData(user.ID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	dashboardData.WeekOfTasks = h.parseTaskListToWeekOfTasks(dashboardTasks.WeekOfTasks, time.Now())
	dashboardData.OverdueTasks = h.parseTaskListToOverdueTasks(dashboardTasks.OverdueTasks)
	dashboardData.UnscheduledTasks = h.parseTaskListToUnscheduledTasks(dashboardTasks.UnscheduledTasks)

	w.Header().Set("Content-Type", "text/html")
	err = views.Layout(dashboard.Index(dashboardData)).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, NewTemplateRenderError(caller, "Layout", err))
		return
	}
}
