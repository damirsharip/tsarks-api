package pguserstorage

import (
	"context"
	"regexp"
	"testing"

	"tech-tsarka/internal/storage/user/pguserstorage/fixture"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func TestStorage_Delete(t *testing.T) {
	// arrange
	f := setUp(t)
	defer f.tearDown()

	// act
	v := fixture.User().ID("d14d0dd6-0ee5-451c-8c1f-2b97658822f2").FirstName("Damir").LastName("Sharip").V()
	queryStore := regexp.QuoteMeta(`DELETE FROM public.users WHERE id = $1`)

	f.pgxPoolMock.ExpectExec(queryStore).
		WithArgs(v.ID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	ctx := context.Background()

	err := f.storage.Delete(ctx, v.ID)

	expectErr := f.pgxPoolMock.ExpectationsWereMet()

	// assert
	require.NoError(t, err)
	require.NoError(t, expectErr)
}
