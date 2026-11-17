package advertisement

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	base "api-tests-template/tests"
)

func TestCreateAdvertisementWithEmptyTitle(t *testing.T) {
	base.SetupSuite()

	t.Log("Предусловие: авторизуем пользователя")

	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	photoPath := filepath.Join(
		"..", "..", "testdata", "думайте.jpg",
	)

	_, err := os.Stat(photoPath)
	require.NoError(t, err, "Тестовая фотография не найдена")

	t.Log("Шаг 1: создаём объявление с пустым названием")

	body := advertisementManager.CreateAdvertisement(
		t,
		loginData.Token,
		"", // Специально оставляем название пустым
		"Проверка валидации обязательного поля title",
		15000,
		2,
		photoPath,
		http.StatusBadRequest,
	)

	t.Log("Шаг 2: проверяем ответ сервера")

	require.NotEmpty(
		t,
		gjson.Get(body, "error").String(),
		"В ответе отсутствует поле error",
	)

	require.Equal(
		t,
		"invalid_request",
		gjson.Get(body, "error").String(),
	)

	require.Equal(
		t,
		"title is required",
		gjson.Get(body, "message").String(),
	)
}
