package user_repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	database "test-go/internal/db"
	"test-go/pkg/config"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

const TEST_ENV_PATH = "../../../.env.test"
const PATH_TO_MAIN = "../../../"

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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := testApp.Start(ctx); err != nil {
		panic(err.Error())
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		testApp.Stop(ctx)
	}()

	if err := runTestMigrations(); err != nil {
		panic("======= Migration failed: " + err.Error() + " =======")
	}

	fmt.Println("======= Migration successful =======")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func runTestMigrations() error {
	if err := os.Chdir(PATH_TO_MAIN); err != nil {
		return err
	}

	cmd := exec.Command("make", "migrate", "DB_URL="+testConfig.DB_URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
