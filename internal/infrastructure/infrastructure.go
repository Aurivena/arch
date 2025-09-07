package infrastructure

import "github.com/jmoiron/sqlx"

type Infrastructure struct {
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func New(sources *Sources) *Infrastructure {
	return &Infrastructure{}
}
