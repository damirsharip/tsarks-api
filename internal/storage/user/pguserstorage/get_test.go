package pguserstorage

import (
	"context"
	"regexp"
	"testing"

	"tech-tsarka/internal/storage/user/pguserstorage/fixture"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_Get(t *testing.T) {
	// arrange
	f := setUp(t)
	defer f.tearDown()

	// act
	v := fixture.User().ID("d14d0dd6-0ee5-451c-8c1f-2b97658822f2").FirstName("Damir").LastName("Sharip").V()
	queryStore := regexp.QuoteMeta(`SELECT first_name, last_name FROM public.users WHERE id = $1`)

	row := pgxmock.NewRows([]string{"first_name", "last_name"}).
		AddRow(v.FirstName, v.LastName)

	f.pgxPoolMock.ExpectQuery(queryStore).
		WithArgs(v.ID).
		WillReturnRows(row)

	ctx := context.Background()

	result, err := f.storage.Get(ctx, v.ID)

	expectErr := f.pgxPoolMock.ExpectationsWereMet()

	// assert
	require.NoError(t, err)
	require.NoError(t, expectErr)

	assert.NotEmpty(t, result)
	assert.Equal(t, result.FirstName, v.FirstName)
	assert.Equal(t, result.LastName, v.LastName)
}
