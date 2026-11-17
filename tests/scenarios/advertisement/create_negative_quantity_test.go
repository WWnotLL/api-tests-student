package advertisement

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	base "api-tests-template/tests"
)

func TestCreateAdvertisementWithNegativeQuantity(t *testing.T) {
	base.SetupSuite()

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

	title := fmt.Sprintf(
		"QA Negative Quantity %d",
		time.Now().UnixNano(),
	)

	t.Log("Шаг 1: создаём объявление с quantity = -1")

	status, body := advertisementManager.CreateAdvertisementWithResult(
		t,
		loginData.Token,
		title,
		"Проверка отрицательного количества",
		15000,
		-1,
		photoPath,
	)

	// Если сервер неожиданно создал объявление,
	// зарегистрируем его удаление.
	id := gjson.Get(body, "id").String()

	if id != "" {
		t.Cleanup(func() {
			t.Logf("Очистка тестового объявления %s", id)

			advertisementManager.DeleteAdvertisement(
				t,
				loginData.Token,
				id,
				http.StatusNoContent,
			)
		})
	}

	t.Log("Шаг 2: проверяем статус ответа")

	require.Equalf(
		t,
		http.StatusBadRequest,
		status,
		"Ожидали 400 Bad Request. Ответ сервера: %s",
		body,
	)

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
		"quantity must be between 1 and 100 inclusive",
		gjson.Get(body, "message").String(),
	)

	t.Logf("Ответ сервера: %s", body)
}
