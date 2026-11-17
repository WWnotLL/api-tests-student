package advertisements

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementsClient "api-tests-template/internal/client/http/advertisements"
)

// ищет объявления и проверяет HTTP-статус ответа.
func SearchAdvertisements(
	t *testing.T,
	title string,
	expectedStatusCode int,
) string {
	response := advertisementsClient.HttpSearchAdvertisements(t, title)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ поиска")

	require.Equalf(
		t,
		expectedStatusCode,
		response.StatusCode,
		"Некорректный статус поиска. Ответ: %s",
		string(body),
	)

	return string(body)
}

// GetAdvertisementPhotos получает фотографии и проверяет HTTP-статус.
func GetAdvertisementPhotos(
	t *testing.T,
	token string,
	id string,
	expectedStatusCode int,
) string {
	response := advertisementsClient.HttpGetAdvertisementPhotos(
		t,
		token,
		id,
	)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ с фотографиями")

	require.Equalf(
		t,
		expectedStatusCode,
		response.StatusCode,
		"Некорректный статус получения фотографий. Ответ: %s",
		string(body),
	)

	return string(body)
}

func ListAdvertisements(t *testing.T) string {
	response := advertisementsClient.HttpListAdvertisements(t)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	require.Equalf(
		t,
		http.StatusOK,
		response.StatusCode,
		"Ошибка получения каталога: %s",
		string(body),
	)

	return string(body)
}
