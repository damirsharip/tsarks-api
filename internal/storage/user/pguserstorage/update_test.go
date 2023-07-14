package pguserstorage

import (
	"context"
	"regexp"
	"testing"

	"tech-tsarka/internal/storage/user/entity"
	"tech-tsarka/internal/storage/user/pguserstorage/fixture"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func TestStorage_Update(t *testing.T) {
	// arrange
	f := setUp(t)
	defer f.tearDown()

	// act
	v := fixture.User().ID("d14d0dd6-0ee5-451c-8c1f-2b97658822f2").FirstName("Damir").LastName("Sharip").V()
	queryStore := regexp.QuoteMeta(`UPDATE public.users SET first_name = $1, last_name = $2 WHERE id = $3`)

	f.pgxPoolMock.ExpectExec(queryStore).
		WithArgs(v.FirstName, v.LastName, v.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	ctx := context.Background()
	arg := entity.UserUpdateInput{
		FirstName: &v.FirstName,
		LastName:  &v.LastName,
	}

	err := f.storage.Update(ctx, v.ID, arg)

	expectErr := f.pgxPoolMock.ExpectationsWereMet()

	// assert
	require.NoError(t, err)
	require.NoError(t, expectErr)
}
