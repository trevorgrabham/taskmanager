package routes

import (
	"fmt"
	"local/taskmanager/internal/account"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/taskinfo"
	"log"
	"net/http"
)

func (h Handlers) ParseTaskForm(w http.ResponseWriter, r *http.Request) {
	var (
		ok            bool
		t             task.Task
		dueDateString string
		timeString    string
		userID int
	)
	if ok = h.checkHXRequest(w, r); !ok {
		return
	}

	if ok = h.checkConnection(w); !ok {
		return
	}

	if ok = h.checkMethod(w, r, http.MethodPost); !ok {
		return
	}

	if ok = h.parseForm(w, r); !ok {
		return
	}

	if t.Title, ok = h.parseTitle(w, r); !ok {
		return
	}

	if userID, ok = r.Context().Value(account.UserContextKey).(int); ok {
		t.UserID = userID
	}

	t.Category = r.FormValue("category")
	t.Description = r.FormValue("description")
	dueDateString = r.FormValue("due-date")
	if dueDateString != "" {
		if timeString == "" {
			if t.DueDate, ok = h.parseDueDate(w, fmt.Sprintf("%s %s", dueDateString, "00:00")); !ok {
				return
			}
		} else {
			if t.DueDate, ok = h.parseDueDate(w, fmt.Sprintf("%s %s", dueDateString, timeString)); !ok {
				return
			}
		}
	}

	if t.ID, ok = h.parseID(w, r.FormValue("id")); !ok {
		return
	}

	if t.RecurringPeriod, ok = h.parseRecurringPeriod(w, r.FormValue("period-value"), r.FormValue("period-unit")); !ok {
		return
	}

	switch {
	case t.ID == 0:
		if ok = h.addTask(w, t); !ok {
			return
		}

		w.Header().Set("HX-Redirect", "/")
	case t.ID > 1:
		if ok = h.updateTask(w, t); !ok {
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err := taskinfo.TaskInfo(t).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			log.Printf("parsing task form: %s\n", err)
			return
		}
	default:
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("parsing task form: bad id value of %d\n", t.ID)
		return
	}
}
