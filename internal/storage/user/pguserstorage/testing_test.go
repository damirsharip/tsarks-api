package pguserstorage

import (
	"testing"

	"github.com/pashagolub/pgxmock"
)

type storageFixtures struct {
	storage     *storage
	pgxPoolMock pgxmock.PgxPoolIface
}

func setUp(t *testing.T) storageFixtures {
	var fixture storageFixtures

	pool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("errored during db initialization: %v", err)
	}

	fixture.pgxPoolMock = pool
	fixture.storage = NewStorage(pool)

	return fixture
}

func (f *storageFixtures) tearDown() {
	f.pgxPoolMock.Close()
}
