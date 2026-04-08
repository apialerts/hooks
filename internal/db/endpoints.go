package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Endpoint struct {
	ID             string
	CreatorIP      string
	ResponseStatus int
	ResponseDelay  int
	ResponseBody   string
	RequestCount   int
	CreatedAt      time.Time
	LastActivityAt time.Time
}

type Queries struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Queries {
	return &Queries{pool: pool}
}

func (q *Queries) CreateEndpoint(ctx context.Context, id, creatorIP string) error {
	_, err := q.pool.Exec(ctx,
		`INSERT INTO endpoints (id, creator_ip) VALUES ($1, $2)`,
		id, creatorIP,
	)
	return err
}

func (q *Queries) GetEndpoint(ctx context.Context, id string) (*Endpoint, error) {
	e := &Endpoint{}
	err := q.pool.QueryRow(ctx,
		`SELECT id, creator_ip, response_status, response_delay_ms, response_body, request_count, created_at, last_activity_at
		 FROM endpoints WHERE id = $1`,
		id,
	).Scan(&e.ID, &e.CreatorIP, &e.ResponseStatus, &e.ResponseDelay, &e.ResponseBody, &e.RequestCount, &e.CreatedAt, &e.LastActivityAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (q *Queries) UpdateEndpointConfig(ctx context.Context, id string, status, delayMs int, body string) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE endpoints SET response_status = $2, response_delay_ms = $3, response_body = $4 WHERE id = $1`,
		id, status, delayMs, body,
	)
	return err
}

func (q *Queries) TouchEndpoint(ctx context.Context, id string) (int, error) {
	var seq int
	err := q.pool.QueryRow(ctx,
		`UPDATE endpoints SET last_activity_at = now(), request_count = request_count + 1 WHERE id = $1 RETURNING request_count`,
		id,
	).Scan(&seq)
	return seq, err
}

func (q *Queries) CountEndpointsByIP(ctx context.Context, ip string) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM endpoints WHERE creator_ip = $1`,
		ip,
	).Scan(&count)
	return count, err
}

func (q *Queries) DeleteEndpoint(ctx context.Context, id string) error {
	_, err := q.pool.Exec(ctx, `DELETE FROM endpoints WHERE id = $1`, id)
	return err
}

func (q *Queries) DeleteExpiredEndpoints(ctx context.Context) (int64, error) {
	result, err := q.pool.Exec(ctx,
		`DELETE FROM endpoints WHERE last_activity_at < now() - INTERVAL '14 days'`,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
