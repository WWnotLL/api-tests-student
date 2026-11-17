package favorites

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
	favoritesManager "api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
)

func TestFavoriteDisappearsAfterAdvertisementDeletion(t *testing.T) {
	base.SetupSuite()

	if os.Getenv("TEST_LOGIN_2") == "" ||
		os.Getenv("TEST_PASSWORD_2") == "" {
		t.Skip("Для этого теста нужен второй аккаунт: TEST_LOGIN_2 и TEST_PASSWORD_2")
	}

	owner := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	buyer := auth.Login(
		t,
		os.Getenv("TEST_LOGIN_2"),
		os.Getenv("TEST_PASSWORD_2"),
	)

	require.NotEqual(t, owner.User.Id, buyer.User.Id)

	photoPath := filepath.Join("..", "..", "testdata", "думайте.jpg")
	_, err := os.Stat(photoPath)
	require.NoError(t, err)

	title := fmt.Sprintf("QA Favorite Cascade %d", time.Now().UnixNano())

	t.Log("Шаг 1: первый пользователь создаёт объявление")

	createBody := advertisementManager.CreateAdvertisement(
		t,
		owner.Token,
		title,
		"Проверка удаления объявления из избранного",
		1000,
		1,
		photoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(t, adID)

	deleted := false

	t.Cleanup(func() {
		if !deleted {
			advertisementManager.DeleteAdvertisement(
				t, owner.Token, adID, http.StatusNoContent,
			)
		}
	})

	t.Log("Шаг 2: второй пользователь добавляет объявление в избранное")

	_, _ = addWithCleanup(t, buyer.Token, adID)

	status, listBody := favoritesManager.ListFavorites(t, buyer.Token)
	require.Equalf(t, http.StatusOK, status, "%s", listBody)
	require.True(t, containsFavorite(listBody, adID))

	t.Log("Шаг 3: владелец удаляет объявление")

	advertisementManager.DeleteAdvertisement(
		t,
		owner.Token,
		adID,
		http.StatusNoContent,
	)
	deleted = true

	t.Log("Шаг 4: проверяем исчезновение из избранного")

	// Несколько попыток на случай небольшой задержки обновления.
	var lastBody string
	for attempt := 0; attempt < 5; attempt++ {
		status, lastBody = favoritesManager.ListFavorites(t, buyer.Token)
		require.Equalf(t, http.StatusOK, status, "%s", lastBody)

		if !containsFavorite(lastBody, adID) {
			return
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf(
		"Удалённое объявление %s осталось в избранном: %s",
		adID,
		lastBody,
	)
}
