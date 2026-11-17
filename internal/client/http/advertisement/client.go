package advertisement

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// запрос на создание объявления с фотографией.
func HttpCreateAdvertisement(
	t *testing.T,
	token string,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", strconv.Itoa(price)).
		MultipartFormData("quantity", strconv.Itoa(quantity)).
		MultipartFile("photos", photoPath).
		Expect(t).
		End().Response
}

// запрос на удаление объявления по ID.
func HttpDeleteAdvertisement(
	t *testing.T,
	token string,
	id string,
) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// получение объявления по его ID.
func HttpGetAdvertisement(
	t *testing.T,
	id string,
) *http.Response {
	return apiRunner.GetRunner().
		Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
