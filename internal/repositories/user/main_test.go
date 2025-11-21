package user_repository

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const TEST_ENV_PATH = "../../../.env.test"

func TestMain(m *testing.M) {
	if err := godotenv.Load(TEST_ENV_PATH); err != nil {
		panic("======= Env loading failed =======")
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}
