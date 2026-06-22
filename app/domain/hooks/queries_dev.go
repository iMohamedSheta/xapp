//go:build dev

package hooks

import (
	"context"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/iMohamedSheta/xqb"
)

// QueryLog holds both the SQL and timing details for a query
type QueryLog struct {
	BuildDuration string `json:"buildDuration"`
	ExecDuration  string `json:"execDuration"`
	Sql           string `json:"sql"`
	RawSql        string `json:"rawSql"`
	Bindings      []any  `json:"bindings"`
	SourceFile    string `json:"sourceFile"`
	SourceLine    int    `json:"sourceLine"`
	startTime     int64
}

var (
	mu           sync.Mutex
	queriesStore = make(map[string][]QueryLog) // requestID → []QueryLog
)

func captureCaller(skipInternal int) (file string, line int) {
	const maxDepth = 20
	pc := make([]uintptr, maxDepth)
	n := runtime.Callers(skipInternal, pc)
	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()
		// Skip our hook and runtime internals
		if !strings.Contains(frame.File, "xqb") &&
			!strings.Contains(frame.File, "/runtime/") &&
			!strings.Contains(frame.File, "hooks.go") {
			return frame.File, frame.Line
		}
		if !more {
			break
		}
	}
	return "", 0
}

// InitQueryHooks wires all query builder hooks
func InitQueryHooks() {
	// Called before query is built
	xqb.DefaultSettings().OnBeforeQuery(func(q *xqb.QueryBuilder) {
		reqID, _ := q.GetContext().Value(enums.ContextKeyRequestId.String()).(string)
		if reqID == "" {
			return
		}

		// skip 4 frames: runtime + hook wrapper
		file, line := captureCaller(4)
		mu.Lock()
		queriesStore[reqID] = append(queriesStore[reqID], QueryLog{
			startTime:  time.Now().UnixNano(),
			SourceFile: file,
			SourceLine: line,
		})
		mu.Unlock()
	})

	// Called after SQL is built (bindings still separate)
	xqb.DefaultSettings().OnAfterQuery(func(q *xqb.QueryExecuted) {
		reqID, ok := q.Context.Value(enums.ContextKeyRequestId.String()).(string)
		if reqID == "" || !ok {
			return
		}

		xqb.Dump(reqID)

		// Inject bindings to get final SQL
		boundSql, err := xqb.InjectBindings(q.Dialect, q.Sql, q.Bindings)
		if err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()

		logs := queriesStore[reqID]
		if len(logs) > 0 {
			last := &logs[len(logs)-1]
			last.BuildDuration = q.Time.String()
			last.Sql = boundSql
			last.RawSql = q.Sql
			last.Bindings = q.Bindings
			queriesStore[reqID] = logs
		}

		xqb.Dump(boundSql)
	})

	// Called after query execution
	xqb.DefaultSettings().OnAfterQueryExecution(func(c context.Context) {
		reqID, _ := c.Value(enums.ContextKeyRequestId.String()).(string)
		if reqID == "" {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		logs := queriesStore[reqID]
		if len(logs) > 0 {
			last := &logs[len(logs)-1]
			last.ExecDuration = time.Since(time.Unix(0, last.startTime)).String()
			queriesStore[reqID] = logs
		}
	})
}

// GetQueries returns all logged queries for a request
func GetQueries(reqID string) []QueryLog {
	mu.Lock()
	defer mu.Unlock()
	return queriesStore[reqID]
}

// ClearQueries removes all logged queries for a request
func ClearQueries(reqID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(queriesStore, reqID)
}
