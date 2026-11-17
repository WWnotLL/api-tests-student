package favorites

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoritesClient "api-tests-template/internal/client/http/favorites"
)

func readResponse(t *testing.T, response *http.Response) (int, string) {
	t.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать HTTP-ответ")

	return response.StatusCode, string(body)
}

func ListFavorites(t *testing.T, token string) (int, string) {
	return readResponse(
		t,
		favoritesClient.HttpGetFavorites(t, token),
	)
}

func AddFavorite(t *testing.T, token, adID string) (int, string) {
	return readResponse(
		t,
		favoritesClient.HttpAddFavorite(t, token, adID),
	)
}

func DeleteFavorite(t *testing.T, token, adID string) (int, string) {
	return readResponse(
		t,
		favoritesClient.HttpDeleteFavorite(t, token, adID),
	)
}
