package runs

import "time"

type QueryType int

const (
	QSuccesful QueryType = iota
	QFailed
	QBefore
	QAfter
)

type RunQuery struct {
	Result []*Run
	Type   QueryType
}

func newRunQuery(runs []*Run, kind QueryType) *RunQuery {
	return &RunQuery{
		Result: runs,
		Type:   kind,
	}
}

func QuerySuccessful(runs []*Run) *RunQuery {
	var queried []*Run
	for _, run := range runs {
		if run.Info.Success {
			queried = append(queried, run)
		}
	}
	return newRunQuery(queried, QSuccesful)
}

func QueryFailed(runs []*Run) *RunQuery {
	var queried []*Run
	for _, run := range runs {
		if !run.Info.Success {
			queried = append(queried, run)
		}
	}
	return newRunQuery(queried, QFailed)
}

func QueryBefore(runs []*Run, timestamp time.Time) *RunQuery {
	var queried []*Run
	for _, run := range runs {
		if timestamp.After(run.Timestamp) {
			queried = append(queried, run)
		}
	}
	return newRunQuery(queried, QBefore)
}

func QueryAfter(runs []*Run, timestamp time.Time) *RunQuery {
	var queried []*Run
	for _, run := range runs {
		if timestamp.Before(run.Timestamp) {
			queried = append(queried, run)
		}
	}
	return newRunQuery(queried, QAfter)
}
