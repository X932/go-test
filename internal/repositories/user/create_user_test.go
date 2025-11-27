package user_repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
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
