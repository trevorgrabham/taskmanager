package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

func (s Service) GetDashboardData(userID int) (dashboardData Dashboard, err error) {
	var (
		caller                                 = "GetDashboardData"
		weekData, overdueData, unscheduledData []sqlite.Task
	)
	if userID < 1 {
		return dashboardData, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if weekData, err = s.repo.GetWeekOfTasks(userID, time.Now()); err != nil {
		return dashboardData, err
	}

	if overdueData, err = s.repo.GetOverdueTasks(userID); err != nil {
		return dashboardData, err
	}

	if unscheduledData, err = s.repo.GetUnscheduledTasks(userID); err != nil {
		return dashboardData, err
	}

	dashboardData.WeekOfTasks = parseRepoTaskListToTaskList(weekData)
	dashboardData.UnscheduledTasks = parseRepoTaskListToTaskList(unscheduledData)
	dashboardData.OverdueTasks = parseRepoTaskListToTaskList(overdueData)
	return dashboardData, nil
}
