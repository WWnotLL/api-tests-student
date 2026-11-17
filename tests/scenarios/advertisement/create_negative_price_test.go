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

func TestCreateAdvertisementWithNegativePrice(t *testing.T) {
	base.SetupSuite()

	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	title := fmt.Sprintf(
		"QA Negative Price %d",
		time.Now().UnixNano(),
	)

	photoPath := filepath.Join(
		"..", "..", "testdata", "думайте.jpg",
	)

	_, err := os.Stat(photoPath)
	require.NoError(t, err, "Тестовая фотография не найдена")

	t.Log("Шаг 1: пытаемся создать объявление с отрицательной ценой")

	body := advertisementManager.CreateAdvertisement(
		t,
		loginData.Token,
		title,
		"Проверка отрицательной цены",
		-100,
		2,
		photoPath,
		http.StatusBadRequest,
	)

	t.Log("Шаг 2: проверяем ошибку валидации цены")

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
		"price must be between 0 and 1000000 inclusive",
		gjson.Get(body, "message").String(),
	)
}
