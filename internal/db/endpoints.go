package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Endpoint struct {
	ID             int64
	Slug           string
	CreatorIP      string
	ResponseStatus int
	ResponseDelay  int
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

func (q *Queries) CreateEndpoint(ctx context.Context, slug, creatorIP string) error {
	_, err := q.pool.Exec(ctx,
		`INSERT INTO endpoint (slug, creator_ip) VALUES ($1, $2)`,
		slug, creatorIP,
	)
	return err
}

func (q *Queries) GetEndpoint(ctx context.Context, slug string) (*Endpoint, error) {
	e := &Endpoint{}
	err := q.pool.QueryRow(ctx,
		`SELECT id, slug, creator_ip, response_status, response_delay_ms, request_count, created_at, last_activity_at
		 FROM endpoint WHERE slug = $1`,
		slug,
	).Scan(&e.ID, &e.Slug, &e.CreatorIP, &e.ResponseStatus, &e.ResponseDelay, &e.RequestCount, &e.CreatedAt, &e.LastActivityAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (q *Queries) UpdateEndpointConfig(ctx context.Context, slug string, status, delayMs int) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE endpoint SET response_status = $2, response_delay_ms = $3 WHERE slug = $1`,
		slug, status, delayMs,
	)
	return err
}

func (q *Queries) TouchEndpoint(ctx context.Context, slug string) (int, error) {
	var seq int
	err := q.pool.QueryRow(ctx,
		`UPDATE endpoint SET last_activity_at = now(), request_count = request_count + 1 WHERE slug = $1 RETURNING request_count`,
		slug,
	).Scan(&seq)
	return seq, err
}

func (q *Queries) CountEndpointsByIP(ctx context.Context, ip string) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM endpoint WHERE creator_ip = $1`,
		ip,
	).Scan(&count)
	return count, err
}

func (q *Queries) ListEndpointsByIP(ctx context.Context, ip string) ([]*Endpoint, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT id, slug, creator_ip, response_status, response_delay_ms, request_count, created_at, last_activity_at
		 FROM endpoint WHERE creator_ip = $1 ORDER BY last_activity_at DESC`,
		ip,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []*Endpoint
	for rows.Next() {
		e := &Endpoint{}
		if err := rows.Scan(&e.ID, &e.Slug, &e.CreatorIP, &e.ResponseStatus, &e.ResponseDelay, &e.RequestCount, &e.CreatedAt, &e.LastActivityAt); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, e)
	}
	return endpoints, rows.Err()
}

func (q *Queries) DeleteEndpoint(ctx context.Context, slug string) error {
	_, err := q.pool.Exec(ctx, `DELETE FROM endpoint WHERE slug = $1`, slug)
	return err
}

func (q *Queries) GetResponseBody(ctx context.Context, endpointID int64, status, delayMs int) (string, error) {
	var body string
	err := q.pool.QueryRow(ctx,
		`SELECT body FROM endpoint_response WHERE endpoint_id = $1 AND status = $2 AND delay_ms = $3`,
		endpointID, status, delayMs,
	).Scan(&body)
	if err != nil {
		return "", err
	}
	return body, nil
}

func (q *Queries) GetAllResponseBodies(ctx context.Context, endpointID int64) (map[string]string, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT status, delay_ms, body FROM endpoint_response WHERE endpoint_id = $1`,
		endpointID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bodies := make(map[string]string)
	for rows.Next() {
		var status, delayMs int
		var body string
		if err := rows.Scan(&status, &delayMs, &body); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%d-%d", status, delayMs)
		bodies[key] = body
	}
	return bodies, rows.Err()
}

func (q *Queries) UpsertResponseBody(ctx context.Context, endpointID int64, status, delayMs int, body string) error {
	_, err := q.pool.Exec(ctx,
		`INSERT INTO endpoint_response (endpoint_id, status, delay_ms, body)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (endpoint_id, status, delay_ms) DO UPDATE SET body = $4`,
		endpointID, status, delayMs, body,
	)
	return err
}

func (q *Queries) DeleteResponseBody(ctx context.Context, endpointID int64, status, delayMs int) error {
	_, err := q.pool.Exec(ctx,
		`DELETE FROM endpoint_response WHERE endpoint_id = $1 AND status = $2 AND delay_ms = $3`,
		endpointID, status, delayMs,
	)
	return err
}

func (q *Queries) DeleteExpiredEndpoints(ctx context.Context) (int64, error) {
	result, err := q.pool.Exec(ctx,
		`DELETE FROM endpoint WHERE last_activity_at < now() - INTERVAL '7 days'`,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
