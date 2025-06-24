//go:build integration

// В этом пакете мы пишем интеграционные тесты для слоя repository или же для слоя usecase
// В будущем нам понадобитмя использовать фикстуры(скрипты для наполнения базы тестовыми данными),
// fixtureManager нужно будет реализовать самостоятельно по аналогии с migrator
package integration

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/timurzdev/mentorship-test-task/cmd"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/repository"
	"github.com/timurzdev/mentorship-test-task/internal/service/helpers"
)

const (
	envFilePath = "../../deploy/local/.test.env"
)

func Test_CreateHouse(t *testing.T) {
	//загружаем в окружение переменные из .env файла
	godotenv.Load(envFilePath)

	type testcase struct {
		name        string
		house       entity.House
		wantErr     bool
		expectedErr error
	}

	testcases := []testcase{
		{
			name: "successful house creation with developer",
			house: entity.House{
				Address:   "Лесная улица, 7, Москва, 125196",
				Year:      2020,
				Developer: helpers.ToPtr("ПИК"),
			},
			wantErr: false,
		},
		{
			name: "successful house creation without developer",
			house: entity.House{
				Address:   "Тверская улица, 1, Москва, 125009",
				Year:      1950,
				Developer: nil,
			},
			wantErr: false,
		},
		{
			name: "fail on duplicate address",
			house: entity.House{
				Address:   "Duplicate Address Test, 123",
				Year:      2000,
				Developer: helpers.ToPtr("Test Developer"),
			},
			wantErr: false, // первое создание должно пройти успешно
		},
		{
			name: "fail with empty address",
			house: entity.House{
				Address:   "",
				Year:      2000,
				Developer: helpers.ToPtr("Developer"),
			},
			wantErr:     false, // TODO: должно быть true когда валидация будет реализована
			expectedErr: entity.ErrorCreatingHouse,
		},
		{
			name: "fail with short address",
			house: entity.House{
				Address:   "abc",
				Year:      2000,
				Developer: nil,
			},
			wantErr:     false, // TODO: должно быть true когда валидация будет реализована
			expectedErr: entity.ErrorCreatingHouse,
		},
		{
			name: "fail with invalid year (too old)",
			house: entity.House{
				Address:   "Старинная улица, 99",
				Year:      1800,
				Developer: nil,
			},
			wantErr:     false, // TODO: должно быть true когда валидация будет реализована
			expectedErr: entity.ErrorCreatingHouse,
		},
		{
			name: "boundary case - minimum valid year",
			house: entity.House{
				Address:   "Историческая площадь, 1",
				Year:      1900,
				Developer: helpers.ToPtr("Старая компания"),
			},
			wantErr: false,
		},
		{
			name: "boundary case - minimum address length",
			house: entity.House{
				Address:   "ул.50",
				Year:      2000,
				Developer: nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем новый контейнер для каждого теста
			container, closer, err := cmd.InitForTesting()
			if err != nil {
				t.Fatal(err)
			}
			defer closer()

			db := container.GetEmbeddedPostgres()

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

			// Создаем подключение к embedded postgres
			conn, err := sqlx.Connect("postgres", container.GetPostgresConnectionString())
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			// Создаем репозиторий с подключением к embedded postgres
			repo := repository.NewRepository(conn)

			_, err = repo.CreateHouse(container.GetContext(), tc.house)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}

}

func Test_CreateHouse_DuplicateAddress(t *testing.T) {
	//загружаем в окружение переменные из .env файла
	godotenv.Load(envFilePath)

	container, closer, err := cmd.InitForTesting()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	db := container.GetEmbeddedPostgres()

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

	// Создаем подключение к embedded postgres
	conn, err := sqlx.Connect("postgres", container.GetPostgresConnectionString())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Создаем репозиторий с подключением к embedded postgres
	repo := repository.NewRepository(conn)

	// Создаем первый дом
	firstHouse := entity.House{
		Address:   "Уникальный адрес для теста",
		Year:      2010,
		Developer: helpers.ToPtr("Застройщик"),
	}

	createdHouse, err := repo.CreateHouse(container.GetContext(), firstHouse)
	assert.NoError(t, err)
	assert.NotNil(t, createdHouse)
	assert.Greater(t, createdHouse.ID, 0)

	// Пытаемся создать второй дом с тем же адресом
	duplicateHouse := entity.House{
		Address:   "Уникальный адрес для теста", // тот же адрес
		Year:      2015,                         // другой год
		Developer: helpers.ToPtr("Другой застройщик"),
	}

	_, err = repo.CreateHouse(container.GetContext(), duplicateHouse)
	assert.Error(t, err)
	assert.ErrorIs(t, err, entity.ErrorCreatingHouse)
}
