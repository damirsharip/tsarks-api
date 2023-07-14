package pguserstorage

import (
	"context"
	"regexp"
	"testing"

	"tech-tsarka/internal/storage/user/entity"
	"tech-tsarka/internal/storage/user/pguserstorage/fixture"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_Create(t *testing.T) {
	// arrange
	f := setUp(t)
	defer f.tearDown()

	// act
	v := fixture.User().ID("d14d0dd6-0ee5-451c-8c1f-2b97658822f2").FirstName("Damir").LastName("Sharip").V()
	queryStore := regexp.QuoteMeta(`INSERT INTO public.users (first_name,last_name) VALUES ($1,$2) RETURNING id`)

	row := pgxmock.NewRows([]string{"id"}).
		AddRow(v.ID)

	f.pgxPoolMock.ExpectQuery(queryStore).
		WithArgs(v.FirstName, v.LastName).
		WillReturnRows(row)

	ctx := context.Background()
	arg := entity.UserCreateInput{
		FirstName: v.FirstName,
		LastName:  v.LastName,
	}

	result, err := f.storage.Create(ctx, arg)

	expectErr := f.pgxPoolMock.ExpectationsWereMet()

	// assert
	require.NoError(t, err)
	require.NoError(t, expectErr)

	assert.NotEmpty(t, result)
	assert.Equal(t, result.ID, v.ID)
}
