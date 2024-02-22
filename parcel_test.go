package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare

	store := NewParcelStore(s.db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и
	//наличии идентификатора

	idAddParc, err := Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, idAddParc)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают
	//со значениями полей в переменной parcel

	getLine, err := Get(idAddParc)
	require.NoError(t, err)
	assert.Equal(t, getLine, parsel)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД

	err := Delete(idAddParc)
	require.NoError(t, err)
	getLine, err := Get(idAddParc)
	require.NoError(t, err)
	assert.Empty(getLine)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии
	// идентификатора

	idAddParc, err := Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, idAddParc)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = SetAdress(idAddParc, newAddress)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	getLine, err := Get(idAddParc)
	require.NoError(t, err)
	assert.NotEmpty(t, getLine)
	assert.Equal(t, getLine.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	idAddParc, err := Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, idAddParc)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки

	err := SetStatus(idAddParc, ParcelStatusSent)
	require.NoError(err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	getLine, err := Get(idAddParc)
	require.NoError(t, err)
	assert.NotEmpty(t, getLine)
	assert.Equal(t, getLine.Status, ParcelStatusSent)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД, убедитесь в отсутствии
		// ошибки и наличии идентификатора
		id, err := Add(parcel)
		require.NoError(t, err)
		assert.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно
		// было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := GetByClient(client)

	// получите список посылок по идентификатору клиента,
	// сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством
	// добавленных

	require.NoError(t, err)
	assert.Equal(t, len(storedParcels), len(parcelMap))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки,
		// значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		valMap, ok := parcelMap[parcel.Number]
		assert.True(t, ok)
		assert.Equal(t, valMap, parcel)
	}
}
