package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.FavoritesPath).
		Query("page", "1").
		Query("page_size", "50").
		Expect(t).
		End().Response
}

func HttpAddFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.FavoritesPath).
		JSON(map[string]string{
			"advertisement_id": adID,
		}).
		Expect(t).
		End().Response
}

func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(path.FavoritesPath + "/" + adID).
		Expect(t).
		End().Response
}
