package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

// GetDashboardData returns the Dashboard data associated with the userID.
//
// If userID is empty, an ErrInvalidUserID is returned.
// Propogates errors returned by GetWeekOfTasks, GetOverdueTasks, and GetUnscheduledTasks.
func (s Service) GetDashboardData(userID int) (dashboardData Dashboard, err error) {
	var weekData, overdueData, unscheduledData []sqlite.Task
	if userID < 1 {
		return Dashboard{}, &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}

	if weekData, err = s.repo.GetWeekOfTasks(userID, time.Now()); err != nil {
		return Dashboard{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if overdueData, err = s.repo.GetOverdueTasks(userID); err != nil {
		return Dashboard{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if unscheduledData, err = s.repo.GetUnscheduledTasks(userID); err != nil {
		return Dashboard{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	dashboardData.WeekOfTasks = parseRepoTaskListToTaskList(weekData)
	dashboardData.UnscheduledTasks = parseRepoTaskListToTaskList(unscheduledData)
	dashboardData.OverdueTasks = parseRepoTaskListToTaskList(overdueData)
	return dashboardData, nil
}
