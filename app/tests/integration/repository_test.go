//go:build integration

// В этом пакете мы пишем интеграционные тесты для слоя repository или же для слоя usecase
// В будущем нам понадобитмя использовать фикстуры(скрипты для наполнения базы тестовыми данными),
// fixtureManager нужно будет реализовать самостоятельно по аналогии с migrator
package integration

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/timurzdev/mentorship-test-task/cmd"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/service/helpers"
)

const (
	envFilePath = "../../.test.env"
)

func Test_CreateHouse(t *testing.T) {
	//загружаем в окружение переменные из .env файла
	godotenv.Load(envFilePath)

	container := cmd.NewInternal(cmd.NewContainer())
	db := container.GetEmbeddedPostgres()

	type testcase struct {
		name        string
		house       entity.House
		wantErr     bool
		expectedErr error
	}

	testcases := []testcase{
		{
			name: "successful house creation",
			house: entity.House{
				Address:   "some address",
				Year:      1999,
				Developer: helpers.ToPtr("samolet"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			//такое поведение нужно, если мы хотим на каждый новый тесткейс иметь чистый постгрес,
			//т.к при вызове db.Stop() мы убиваем процесс постгреса со всеми данными внутри него
			if err := db.Start(); err != nil {
				t.Fatal(err)
			}

			defer func() {
				if err := db.Stop(); err != nil {
					t.Fatal(err)
				}
			}()

			// накатываем наши миграции
			if err := container.GetMigrator().MigrateUp(); err != nil {
				t.Fatal(err)
			}

			// перед инициализацией репозиторию нужно запустить embeded postgres: db.Start()
			repo := container.GetRepository()
			_, err := repo.CreateHouse(container.GetGlobalContext(), tc.house)
			if tc.wantErr {
				assert.Equal(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}

}
