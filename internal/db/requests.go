package db

import (
	"context"
	"time"
)

type Request struct {
	ID             int64
	EndpointID     string
	Seq            int
	Method         string
	Path           string
	QueryParams    string
	Headers        []byte
	Body           string
	BodySize       int
	ContentType    string
	SourceIP       string
	ResponseStatus *int
	ResponseBody   string
	ReceivedAt     time.Time
}

func (q *Queries) CreateRequest(ctx context.Context, req *Request) error {
	_, err := q.pool.Exec(ctx,
		`INSERT INTO requests (endpoint_id, seq, method, path, query_params, headers, body, body_size, content_type, source_ip, response_status, response_body)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		req.EndpointID, req.Seq, req.Method, req.Path, req.QueryParams, req.Headers, req.Body, req.BodySize, req.ContentType, req.SourceIP, req.ResponseStatus, req.ResponseBody,
	)
	return err
}

func (q *Queries) ListRequests(ctx context.Context, endpointID string, limit int) ([]Request, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT id, endpoint_id, seq, method, path, query_params, headers, body, body_size, content_type, source_ip, response_status, response_body, received_at
		 FROM requests WHERE endpoint_id = $1 ORDER BY received_at DESC LIMIT $2`,
		endpointID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var r Request
		err := rows.Scan(&r.ID, &r.EndpointID, &r.Seq, &r.Method, &r.Path, &r.QueryParams, &r.Headers, &r.Body, &r.BodySize, &r.ContentType, &r.SourceIP, &r.ResponseStatus, &r.ResponseBody, &r.ReceivedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}
	return requests, nil
}

func (q *Queries) GetRequest(ctx context.Context, endpointID string, seq int) (*Request, error) {
	r := &Request{}
	err := q.pool.QueryRow(ctx,
		`SELECT id, endpoint_id, seq, method, path, query_params, headers, body, body_size, content_type, source_ip, response_status, response_body, received_at
		 FROM requests WHERE endpoint_id = $1 AND seq = $2`,
		endpointID, seq,
	).Scan(&r.ID, &r.EndpointID, &r.Seq, &r.Method, &r.Path, &r.QueryParams, &r.Headers, &r.Body, &r.BodySize, &r.ContentType, &r.SourceIP, &r.ResponseStatus, &r.ResponseBody, &r.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (q *Queries) GetRequestCount(ctx context.Context, endpointID string) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx,
		`SELECT request_count FROM endpoints WHERE id = $1`,
		endpointID,
	).Scan(&count)
	return count, err
}

func (q *Queries) TrimRequests(ctx context.Context, endpointID string, keepCount int) error {
	_, err := q.pool.Exec(ctx,
		`DELETE FROM requests WHERE endpoint_id = $1 AND id NOT IN (
			SELECT id FROM requests WHERE endpoint_id = $1 ORDER BY received_at DESC LIMIT $2
		)`,
		endpointID, keepCount,
	)
	return err
}
