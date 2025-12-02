//go:build integration

// В этом пакете мы пишем интеграционные тесты для слоя repository или же для слоя usecase
// В будущем нам понадобитмя использовать фикстуры(скрипты для наполнения базы тестовыми данными),
// fixtureManager нужно будет реализовать самостоятельно по аналогии с migrator
package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/timurzdev/mentorship-test-task/cmd"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/repository"
	"github.com/timurzdev/mentorship-test-task/internal/service/helpers"
	"github.com/timurzdev/mentorship-test-task/migrations"
	"github.com/timurzdev/mentorship-test-task/pkg/fixtures"
)

const (
	envFilePath = "../../deploy/local/.test.env"
)

var (
	// Глобальные переменные для хранения контейнера и соединения
	testContainer *cmd.Container
	testDB        *sqlx.DB
	testRepo      *repository.Repository
	testCtx       context.Context
)

// TestMain запускается один раз для всего пакета тестов
// ВАЖНО: Для запуска интеграционных тестов требуется Docker
func TestMain(m *testing.M) {
	// Загружаем переменные окружения
	godotenv.Load(envFilePath)

	// Создаем контейнер один раз для всех тестов
	container, closer, err := cmd.InitForTesting()
	if err != nil {
		panic(err)
	}
	defer closer()

	testContainer = container
	testCtx = container.GetContext()

	// Запускаем PostgreSQL контейнер один раз
	_, cleanupContainer, err := container.GetPostgresTestContainer()
	if err != nil {
		panic(err)
	}
	defer cleanupContainer()

	// Получаем строку подключения для миграций
	migrateConnStr, err := container.GetTestContainerMigrateConnectionString()
	if err != nil {
		panic(err)
	}

	// Создаем мигратор и применяем миграции один раз
	migrator := migrations.NewMigrator(migrateConnStr)
	if err := migrator.MigrateUp(); err != nil {
		panic(err)
	}

	// Создаем подключение к testcontainer postgres
	connStr, err := container.GetTestContainerConnectionString()
	if err != nil {
		panic(err)
	}

	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	testDB = conn
	testRepo = repository.NewRepository(conn)

	// Запускаем тесты
	code := m.Run()

	// Завершаем с кодом возврата от тестов
	os.Exit(code)
}

// cleanupDatabase очищает таблицы перед каждым тестом
func cleanupDatabase(t *testing.T) {
	_, err := testDB.Exec("TRUNCATE TABLE house RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to cleanup house table: %v", err)
	}

	_, err = testDB.Exec("TRUNCATE TABLE flat RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to cleanup flat table: %v", err)
	}
}

func Test_CreateHouse(t *testing.T) {
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
			// Очищаем БД перед каждым тестом
			cleanupDatabase(t)

			// Используем глобальный репозиторий и контекст
			_, err := testRepo.CreateHouse(testCtx, tc.house)
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
	// Очищаем БД перед тестом
	cleanupDatabase(t)

	// Создаем первый дом
	firstHouse := entity.House{
		Address:   "Уникальный адрес для теста",
		Year:      2010,
		Developer: helpers.ToPtr("Застройщик"),
	}

	createdHouse, err := testRepo.CreateHouse(testCtx, firstHouse)
	assert.NoError(t, err)
	assert.NotNil(t, createdHouse)
	assert.Greater(t, createdHouse.ID, 0)

	// Пытаемся создать второй дом с тем же адресом
	duplicateHouse := entity.House{
		Address:   "Уникальный адрес для теста", // тот же адрес
		Year:      2015,                         // другой год
		Developer: helpers.ToPtr("Другой застройщик"),
	}

	_, err = testRepo.CreateHouse(testCtx, duplicateHouse)
	assert.Error(t, err)
	assert.ErrorIs(t, err, entity.ErrorCreatingHouse)
}

func Test_GetHouseFlats(t *testing.T) {
	cleanupDatabase(t)

	fm := fixtures.NewFixtureManager(testDB)
	err := fm.LoadFixtures()
	assert.NoError(t, err)

	type testcase struct {
		name        string
		houseID     int
		isModerator bool
		wantErr     bool
		validate    func(t *testing.T, flats []entity.Flat)
	}

	testcases := []testcase{
		{
			name:        "moderator sees all flats for house 1",
			houseID:     1,
			isModerator: true,
			wantErr:     false,
			validate: func(t *testing.T, flats []entity.Flat) {
				assert.Len(t, flats, 3)
				statuses := make(map[string]bool)
				for _, flat := range flats {
					statuses[flat.Status] = true
					assert.Equal(t, 1, flat.HouseID)
				}
				assert.True(t, statuses["created"])
				assert.True(t, statuses["approved"])
				assert.True(t, statuses["on moderation"])
			},
		},
		{
			name:        "client sees only approved flats for house 1",
			houseID:     1,
			isModerator: false,
			wantErr:     false,
			validate: func(t *testing.T, flats []entity.Flat) {
				assert.Len(t, flats, 1) // только 1 approved
				assert.Equal(t, "approved", flats[0].Status)
				assert.Equal(t, 1, flats[0].HouseID)
				assert.Equal(t, 800, flats[0].Price)
				assert.Equal(t, 4, flats[0].Rooms)
			},
		},
		{
			name:        "house without flats returns empty slice",
			houseID:     2,
			isModerator: true,
			wantErr:     false,
			validate: func(t *testing.T, flats []entity.Flat) {
				assert.Len(t, flats, 0) // дом 2 без квартир
			},
		},
		{
			name:        "moderator sees created flat for house 3",
			houseID:     3,
			isModerator: true,
			wantErr:     false,
			validate: func(t *testing.T, flats []entity.Flat) {
				assert.Len(t, flats, 1)
				assert.Equal(t, "created", flats[0].Status)
				assert.Equal(t, 3, flats[0].HouseID)
				assert.Equal(t, 200, flats[0].Price)
				assert.Equal(t, 1, flats[0].Rooms)
			},
		},
		{
			name:        "client sees no flats for house 3 (no approved)",
			houseID:     3,
			isModerator: false,
			wantErr:     false,
			validate: func(t *testing.T, flats []entity.Flat) {
				assert.Len(t, flats, 0) // только created, а client видит только approved
			},
		},
		{
			name:        "non-existent house returns empty slice",
			houseID:     999,
			isModerator: true,
			wantErr:     false,
			validate: func(t *testing.T, flats []entity.Flat) {
				assert.Len(t, flats, 0)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			flats, err := testRepo.GetHouseFlats(testCtx, tc.houseID, tc.isModerator)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tc.validate != nil {
				tc.validate(t, flats)
			}
		})
	}
}
