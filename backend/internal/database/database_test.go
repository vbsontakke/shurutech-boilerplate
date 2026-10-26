package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vbsontakke/boilerplate/internal/database"
	"github.com/vbsontakke/boilerplate/internal/testutil"
)

// Integration tests: skipped automatically when the test DB is not running.
// Start it with `make db-up`.

func TestConnectAndMigrate(t *testing.T) {
	pool := testutil.NewTestDB(t)

	var n int
	err := pool.QueryRow(context.Background(), `SELECT count(*) FROM schema_migrations`).Scan(&n)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, n, 1)
}

func TestWithTx_RollsBackOnError(t *testing.T) {
	pool := testutil.NewTestDB(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS _tx_probe (v INT UNIQUE)`)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DROP TABLE IF EXISTS _tx_probe`) })

	boom := errors.New("boom")
	err = database.WithTx(ctx, pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO _tx_probe (v) VALUES (1)`); err != nil {
			return err
		}
		return boom
	})
	require.ErrorIs(t, err, boom)

	var count int
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM _tx_probe`).Scan(&count))
	assert.Equal(t, 0, count, "insert should have been rolled back")
}

func TestIsUniqueViolation(t *testing.T) {
	pool := testutil.NewTestDB(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS _uniq_probe (v INT UNIQUE)`)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DROP TABLE IF EXISTS _uniq_probe`) })

	_, err = pool.Exec(ctx, `INSERT INTO _uniq_probe (v) VALUES (1)`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO _uniq_probe (v) VALUES (1)`)
	assert.True(t, database.IsUniqueViolation(err))
}
