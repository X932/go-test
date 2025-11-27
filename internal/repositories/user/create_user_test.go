package user_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	startErr := testApp.Start(ctx)
	require.NoError(t, startErr, "Failed to start fx test ====")

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		testApp.Stop(ctx)
	}()

	hashedPassword := []byte("$2a$10$h0smPzTUtamCxzIljsQ7Nugpsn.jvJvmg6huQSh5rXWfIryAwqsZa")
	newUser := CreateUser{
		FirstName: "Yusuf",
		LastName:  "g",
		Email:     "esc@sdc.cs",
		Password:  hashedPassword,
	}

	rowsAffected, repoErr := testUserRepo.CreateUser(newUser)

	require.NoError(t, repoErr)
	assert.Equal(t, int64(1), rowsAffected, "Result is wrong")

	_, _ = testDB.Exec("DELETE FROM \"user\";")
}
