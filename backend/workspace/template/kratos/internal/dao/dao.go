package dao

import "context"

type Dao struct {
	// [DEFINITION]
}

func New() *Dao {
	d := &Dao{
		// [NEW]
	}
	// [INIT]
	return d
}

func (d *Dao) Ping(ctx context.Context) error {
	// [PING]
	return nil
}

func (d *Dao) Close() error {
	// [CLOSE]
	return nil
}