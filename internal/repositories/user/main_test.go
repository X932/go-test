package user_repository

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	database "test-go/internal/db"
	"test-go/pkg/config"
	"testing"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

const TEST_ENV_PATH = "../../../.env.test"

var testUserRepo Repo
var testDB *sql.DB
var testConfig *config.Config
var testApp *fx.App

func TestMain(m *testing.M) {
	if err := godotenv.Load(TEST_ENV_PATH); err != nil {
		panic("======= Env loading failed =======")
	}

	testApp = fx.New(
		database.Module,
		config.Module,
		Module,
		fx.NopLogger,
		fx.Populate(&testUserRepo, &testDB, &testConfig),
	)

	if err := runTestMigrations(); err != nil {
		panic("======= Migration failed: " + err.Error() + " =======")
	}

	fmt.Println("======= Migration successful  =======")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func runTestMigrations() error {

	if err := os.Chdir("../../../"); err != nil {
		return err
	}

	cmd := exec.Command("make", "migrate", "DB_URL="+testConfig.DB_URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
