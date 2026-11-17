package advertisement

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	base "api-tests-template/tests"
)

func TestCreateAdvertisementWithInvalidToken(t *testing.T) {
	base.SetupSuite()

	photoPath := filepath.Join(
		"..", "..", "testdata", "думайте.jpg",
	)

	_, err := os.Stat(photoPath)
	require.NoError(t, err, "Тестовая фотография не найдена")

	t.Log("Шаг 1: создаём объявление с недействительным токеном")

	body := advertisementManager.CreateAdvertisement(
		t,
		"invalid-token",
		"QA Invalid Token Test",
		"Проверка недействительного Bearer-токена",
		15000,
		2,
		photoPath,
		http.StatusUnauthorized,
	)

	t.Log("Шаг 2: проверяем отказ в авторизации")

	require.Equal(
		t,
		"unauthorized",
		gjson.Get(body, "error").String(),
	)

	require.Equal(
		t,
		"Invalid or expired token",
		gjson.Get(body, "message").String(),
	)
}
