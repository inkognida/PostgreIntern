package queries

import (
	"context"
	"postgreintern/internal/model"
)

const insertFileSystemFileEvent = `insert into system_file_events (event_type, path, file_name) values ($1, $2, $3) returning (id)`
func (q *Queries) SaveEvent(ctx context.Context, event model.FileEvent) (int, error) {
	q.mu.Lock()
	var id int
	err := q.pool.QueryRow(ctx, insertFileSystemFileEvent, event.EventType, event.Path, event.FileName).Scan(&id)
	if err != nil {
		q.mu.Unlock()
		return 0, err
	}
	q.mu.Unlock()
	return id, nil
}

const insertCommandEvent = `insert into commands_events (command, args, system_event_id) values ($1, $2, $3)`
func (q *Queries) SaveCommandExecution(ctx context.Context, event model.CmdEvent, id int) error {
	q.mu.Lock()
	_, err := q.pool.Exec(ctx, insertCommandEvent, event.Cmd, event.Args, id)
	if err != nil {
		q.mu.Unlock()
		return err
	}
	q.mu.Unlock()
	return nil
}
