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
	base "api-tests-template/tests"
)

func TestCreateAdvertisementWithoutAuth(t *testing.T) {
	base.SetupSuite()

	title := fmt.Sprintf("QA Negative Test %d", time.Now().UnixNano())
	description := "Проверка создания объявления без авторизации"

	photoPath := filepath.Join(
		"..", "..", "testdata", "думайте.jpg",
	)

	_, err := os.Stat(photoPath)
	require.NoError(t, err, "Тестовая фотография не найдена")

	t.Log("Шаг 1: отправляем запрос на создание без токена")

	body := advertisementManager.CreateAdvertisement(
		t,
		"", // Пустой токен: Authorization не передаётся
		title,
		description,
		15000,
		2,
		photoPath,
		http.StatusUnauthorized,
	)

	t.Log("Шаг 2: проверяем ответ об отсутствии авторизации")

	require.NotEmpty(
		t,
		gjson.Get(body, "error").String(),
		"В ответе отсутствует описание ошибки",
	)

	t.Log("Создание объявления без авторизации отклонено")
}
