package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
	"time"
)

// Renders the index.html page.
//
// Requires that the user be authenticated.
func (h Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller         = "IndexHandler"
		err            error
		sessionID      string
		user           services.User
		dashboardTasks services.Dashboard
		dashboardData  views.DashboardViewData
	)

	// Check user authenticated and grab User
	if r.URL.Path != "/" {
		h.HandleError(w, r, fmt.Errorf("%s: %w %s", caller, ErrUnknownPath, r.URL.Path))
		return
	}

	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptySessionID))
		return
	}

	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", err))
		return
	}

	if dashboardTasks, err = h.services.GetDashboardData(user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", err))
		return
	}

	dashboardData.WeekOfTasks = h.parseTaskListToWeekOfTasks(dashboardTasks.WeekOfTasks, time.Now())
	dashboardData.OverdueTasks = h.parseTaskListToOverdueTasks(dashboardTasks.OverdueTasks)
	dashboardData.UnscheduledTasks = h.parseTaskListToUnscheduledTasks(dashboardTasks.UnscheduledTasks)
	dashboardData.StartDay = time.Now()

	w.Header().Set("Content-Type", "text/html")
	err = views.Layout(views.LayoutViewData{UserID: user.ID, Content: views.Index(dashboardData)}).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w Layout: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
