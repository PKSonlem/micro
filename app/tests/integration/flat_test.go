//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/PKSonlem/micro/internal/entity"
	"github.com/PKSonlem/micro/internal/service/helpers"
)

// Проверяет успешное создание квартиры с автогенерацией flat_number
func Test_CreateFlat_Success(t *testing.T) {
	// Очищаем БД перед тестом
	cleanupDatabase(t)

	// Создаем дом
	house := entity.House{
		Address:   "Тестовая улица, 1",
		Year:      2020,
		Developer: helpers.ToPtr("Тестовый застройщик"),
	}
	createdHouse, err := testRepo.CreateHouse(testCtx, house)
	assert.NoError(t, err)
	assert.NotNil(t, createdHouse)

	// Создаем первую квартиру
	firstFlat := entity.Flat{
		HouseID: createdHouse.ID,
		Price:   5000000,
		Rooms:   3,
	}
	createdFlat, err := testRepo.CreateFlat(testCtx, firstFlat)
	assert.NoError(t, err)
	assert.NotNil(t, createdFlat)
	assert.Greater(t, createdFlat.ID, 0)
	assert.Equal(t, "created", createdFlat.Status)
	assert.Equal(t, 1, createdFlat.FlatNumber) // Первая квартира должна иметь номер 1

	// Создаем вторую квартиру в том же доме
	secondFlat := entity.Flat{
		HouseID: createdHouse.ID,
		Price:   6000000,
		Rooms:   2,
	}
	createdSecondFlat, err := testRepo.CreateFlat(testCtx, secondFlat)
	assert.NoError(t, err)
	assert.NotNil(t, createdSecondFlat)
	assert.Equal(t, 2, createdSecondFlat.FlatNumber) // Вторая квартира должна иметь номер 2
}

// Проверяет, что нельзя создать квартиру в несуществующем доме
func Test_CreateFlat_NonExistentHouse(t *testing.T) {
	// Очищаем БД перед тестом
	cleanupDatabase(t)

	// Пытаемся создать квартиру с house_id, который не существует
	flat := entity.Flat{
		HouseID: 99999,
		Price:   5000000,
		Rooms:   3,
	}

	_, err := testRepo.CreateFlat(testCtx, flat)
	assert.Error(t, err)
	assert.ErrorIs(t, err, entity.ErrorCreatingFlat)
}
