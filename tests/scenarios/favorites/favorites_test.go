package favorites

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	advertisementsManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	favoritesManager "api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
)

// Авторизация основного тестового пользователя.
func login(t *testing.T) (string, string) {
	t.Helper()
	data := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)
	return data.Token, data.User.Id
}

// Проверяем наличие конкретного объявления в избранном.
func containsFavorite(body, adID string) bool {
	for _, item := range gjson.Get(body, "items").Array() {
		if item.Get("advertisement.id").String() == adID {
			return true
		}
	}
	return false
}

// Находим чужое объявление, которого ещё нет в избранном.
func selectForeignAdvertisement(t *testing.T, token, userID string) string {
	t.Helper()

	status, favoritesBody := favoritesManager.ListFavorites(t, token)
	require.Equalf(t, http.StatusOK, status, "%s", favoritesBody)

	if gjson.Get(favoritesBody, "total").Int() >= 50 {
		t.Skip("Избранное заполнено: достигнут лимит 50 объявлений")
	}

	existing := make(map[string]bool)

	for _, item := range gjson.Get(favoritesBody, "items").Array() {
		existing[item.Get("advertisement.id").String()] = true
	}

	catalogBody := advertisementsManager.ListAdvertisements(t)

	for _, ad := range gjson.Get(catalogBody, "items").Array() {
		id := ad.Get("id").String()
		ownerID := ad.Get("user_id").String()

		if id != "" && ownerID != "" &&
			ownerID != userID && !existing[id] {
			return id
		}
	}

	t.Skip("В первой странице каталога нет подходящего чужого объявления")
	return ""
}

// Добавляем объявление и заранее регистрируем очистку.
func addWithCleanup(
	t *testing.T,
	token, adID string,
) (string, func()) {
	t.Helper()

	status, body := favoritesManager.AddFavorite(t, token, adID)

	active := status == http.StatusCreated

	if active {
		t.Cleanup(func() {
			if !active {
				return
			}

			cleanupStatus, cleanupBody := favoritesManager.DeleteFavorite(
				t, token, adID,
			)

			require.Truef(
				t,
				cleanupStatus == http.StatusNoContent ||
					cleanupStatus == http.StatusNotFound,
				"Ошибка очистки избранного: %d %s",
				cleanupStatus,
				cleanupBody,
			)
		})
	}

	require.Equalf(
		t,
		http.StatusCreated,
		status,
		"Ошибка добавления в избранное: %s",
		body,
	)

	return body, func() { active = false }
}

// Генерируем корректно сформированный UUID для несуществующего объекта.
func missingAdvertisementID(t *testing.T) string {
	t.Helper()

	var b [16]byte
	_, err := rand.Read(b[:])
	require.NoError(t, err)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
	)
}

// Позитивный сценарий: добавить → увидеть → удалить → не увидеть.
func TestFavoritesFullCycle(t *testing.T) {
	base.SetupSuite()
	token, userID := login(t)

	adID := selectForeignAdvertisement(t, token, userID)

	t.Log("Шаг 1: добавление объявления в избранное")
	body, markRemoved := addWithCleanup(t, token, adID)

	require.NotEmpty(t, gjson.Get(body, "favorite_id").String())
	require.NotEmpty(t, gjson.Get(body, "added_at").String())
	require.Equal(t, adID, gjson.Get(body, "advertisement.id").String())

	t.Log("Шаг 2: проверка списка избранного")
	status, listBody := favoritesManager.ListFavorites(t, token)
	require.Equalf(t, http.StatusOK, status, "%s", listBody)
	require.True(t, containsFavorite(listBody, adID))

	t.Log("Шаг 3: удаление из избранного")
	status, deleteBody := favoritesManager.DeleteFavorite(t, token, adID)
	require.Equalf(t, http.StatusNoContent, status, "%s", deleteBody)
	markRemoved()

	t.Log("Шаг 4: проверка отсутствия объявления")
	status, listBody = favoritesManager.ListFavorites(t, token)
	require.Equalf(t, http.StatusOK, status, "%s", listBody)
	require.False(t, containsFavorite(listBody, adID))
}

// Повторное добавление должно возвращать конфликт.
func TestFavoritesDuplicate(t *testing.T) {
	base.SetupSuite()
	token, userID := login(t)
	adID := selectForeignAdvertisement(t, token, userID)

	_, _ = addWithCleanup(t, token, adID)

	status, body := favoritesManager.AddFavorite(t, token, adID)

	require.Equalf(t, http.StatusConflict, status, "%s", body)
	require.NotEmpty(t, gjson.Get(body, "error").String())
}

// Собственное объявление нельзя добавлять в избранное.
func TestFavoritesCannotAddOwnAdvertisement(t *testing.T) {
	base.SetupSuite()
	token, _ := login(t)

	photoPath := filepath.Join("..", "..", "testdata", "думайте.jpg")
	_, err := os.Stat(photoPath)
	require.NoError(t, err)

	title := fmt.Sprintf("QA Own Ad %d", time.Now().UnixNano())

	createBody := advertisementManager.CreateAdvertisement(
		t,
		token,
		title,
		"Проверка запрета на добавление собственного объявления",
		1000,
		1,
		photoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(t, adID)

	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(
			t, token, adID, http.StatusNoContent,
		)
	})

	status, body := favoritesManager.AddFavorite(t, token, adID)

	require.Equalf(t, http.StatusForbidden, status, "%s", body)
	require.NotEmpty(t, gjson.Get(body, "error").String())
}

// Несуществующее объявление нельзя добавить или удалить.
func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	base.SetupSuite()
	token, _ := login(t)
	adID := missingAdvertisementID(t)

	t.Run("add", func(t *testing.T) {
		status, body := favoritesManager.AddFavorite(t, token, adID)
		require.Equalf(t, http.StatusNotFound, status, "%s", body)
	})

	t.Run("delete", func(t *testing.T) {
		status, body := favoritesManager.DeleteFavorite(t, token, adID)
		require.Equalf(t, http.StatusNotFound, status, "%s", body)
	})
}

// Все три операции должны отклоняться без авторизации.
func TestFavoritesWithoutAuthorization(t *testing.T) {
	base.SetupSuite()
	adID := missingAdvertisementID(t)

	t.Run("list", func(t *testing.T) {
		status, body := favoritesManager.ListFavorites(t, "")
		require.Equalf(t, http.StatusUnauthorized, status, "%s", body)
	})

	t.Run("add", func(t *testing.T) {
		status, body := favoritesManager.AddFavorite(t, "", adID)
		require.Equalf(t, http.StatusUnauthorized, status, "%s", body)
	})

	t.Run("delete", func(t *testing.T) {
		status, body := favoritesManager.DeleteFavorite(t, "", adID)
		require.Equalf(t, http.StatusUnauthorized, status, "%s", body)
	})
}
