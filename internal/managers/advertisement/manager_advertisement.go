package advertisement

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementClient "api-tests-template/internal/client/http/advertisement"
)

// создание объявления и проверка HTTP-статуса ответа.
func CreateAdvertisement(
	t *testing.T,
	token string,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpCreateAdvertisement(
		t,
		token,
		title,
		description,
		price,
		quantity,
		photoPath,
	)
	defer response.Body.Close()

	require.Equal(
		t,
		expectedStatusCode,
		response.StatusCode,
		"Некорректный HTTP-статус при создании объявления",
	)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	return string(body)
}

// удаление объявления и проверка HTTP-статус ответа.
func DeleteAdvertisement(
	t *testing.T,
	token string,
	id string,
	expectedStatusCode int,
) {
	response := advertisementClient.HttpDeleteAdvertisement(t, token, id)
	defer response.Body.Close()

	require.Equal(
		t,
		expectedStatusCode,
		response.StatusCode,
		"Некорректный HTTP-статус при удалении объявления",
	)
}

// получение объявления и проверка HTTP-статус ответа.
func GetAdvertisement(
	t *testing.T,
	id string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpGetAdvertisement(t, id)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	require.Equalf(
		t,
		expectedStatusCode,
		response.StatusCode,
		"Некорректный статус GET объявления %s. Ответ: %s",
		id,
		string(body),
	)

	return string(body)
}

// CreateAdvertisementWithResult возвращает фактический статус и тело ответа.
func CreateAdvertisementWithResult(
	t *testing.T,
	token, title, description string,
	price, quantity int,
	photoPath string,
) (int, string) {
	response := advertisementClient.HttpCreateAdvertisement(
		t,
		token,
		title,
		description,
		price,
		quantity,
		photoPath,
	)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	return response.StatusCode, string(body)
}
