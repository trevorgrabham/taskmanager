package task

import (
	"fmt"
	"strings"
)

type TaskMetaData struct {
	TaskName  string
	Count     int
	LastCompleted TaskDueDate
}

func (t TaskMetaData) String() string {
	return fmt.Sprintf("%s | %d times | last completed %s\n", t.TaskName, t.Count, t.LastCompleted.String())
}

// ============================================== TaskMetaDataList ==============================================

type TaskMetaDataList []TaskMetaData

func (t TaskMetaDataList) String() string {
	metaDataStrings := make([]string, len(t))
	for i := range t {
		metaDataStrings[i] = t[i].String()
	}
	return strings.Join(metaDataStrings, "\n")
}

// ============================================== TaskTime ==============================================

type TaskTime int

func (t TaskTime) String() string {
	if t > 60 {
		return fmt.Sprintf("%dh%dm", t/60, t%60)
	}
	return fmt.Sprintf("%dm", t)
}
