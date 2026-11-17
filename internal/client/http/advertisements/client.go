package advertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// ищет объявления по названию.
func HttpSearchAdvertisements(
	t *testing.T,
	title string,
) *http.Response {
	return apiRunner.GetRunner().
		Create().
		Get(path.AdvertisementsPath).
		Query("search", title).
		Query("limit", "50").
		Expect(t).
		End().Response
}

// получает фотографии объявления.
func HttpGetAdvertisementPhotos(
	t *testing.T,
	token string,
	id string,
) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.AdvertisementsPath + "/" + id + "/photos").
		Expect(t).
		End().Response
}

func HttpListAdvertisements(t *testing.T) *http.Response {
	return apiRunner.GetRunner().
		Create().
		Get(path.AdvertisementsPath).
		Query("limit", "50").
		Expect(t).
		End().Response
}
