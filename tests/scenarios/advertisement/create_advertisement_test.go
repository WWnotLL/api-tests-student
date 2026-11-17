package advertisement

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"

	advertisementsManager "api-tests-template/internal/managers/advertisements"
)

type CreateAdvertisementSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestCreateAdvertisementSuite(t *testing.T) {
	suite.Run(t, new(CreateAdvertisementSuite))
}

func (s *CreateAdvertisementSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация тестового пользователя")

	s.loginData = auth.Login(
		s.T(),
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)
}

func (s *CreateAdvertisementSuite) TestCreateWithAllFields() {
	title := fmt.Sprintf("QA API Test %d", time.Now().UnixNano())
	description := "Тестовое объявление для проверки API"
	price := 15000
	quantity := 2

	photoPath := filepath.Join(
		"..", "..", "testdata", "думайте.jpg",
	)

	s.T().Log("Шаг 1: создание объявления со всеми полями")

	body := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		description,
		price,
		quantity,
		photoPath,
		http.StatusCreated,
	)

	s.T().Log("Шаг 2: проверка ответа на создание")

	id := gjson.Get(body, "id").String()

	require.NotEmpty(s.T(), id, "ID объявления не получен")

	s.T().Cleanup(func() {
		s.T().Logf("Очистка: удаляем тестовое объявление %s", id)

		advertisementManager.DeleteAdvertisement(
			s.T(),
			s.loginData.Token,
			id,
			http.StatusNoContent,
		)
	})

	require.Equal(s.T(), title, gjson.Get(body, "title").String())
	require.Equal(s.T(), description, gjson.Get(body, "description").String())
	require.Equal(s.T(), int64(price), gjson.Get(body, "price").Int())
	require.Equal(s.T(), int64(quantity), gjson.Get(body, "quantity").Int())

	photos := gjson.Get(body, "photos").Array()
	require.Len(s.T(), photos, 1, "Ожидалась одна фотография")

	s.T().Logf("Объявление создано: %s", id)

	s.T().Log("Шаг 3: получаем созданное объявление по ID")

	getBody := advertisementManager.GetAdvertisement(
		s.T(),
		id,
		http.StatusOK,
	)

	s.T().Log("Шаг 4: проверяем сохранённые данные объявления")

	require.Equal(
		s.T(),
		id,
		gjson.Get(getBody, "id").String(),
		"ID полученного объявления не совпадает",
	)

	require.Equal(
		s.T(),
		title,
		gjson.Get(getBody, "title").String(),
		"Название объявления не совпадает",
	)

	require.Equal(
		s.T(),
		description,
		gjson.Get(getBody, "description").String(),
		"Описание объявления не совпадает",
	)

	require.Equal(
		s.T(),
		int64(price),
		gjson.Get(getBody, "price").Int(),
		"Цена объявления не совпадает",
	)

	require.Equal(
		s.T(),
		int64(quantity),
		gjson.Get(getBody, "quantity").Int(),
		"Количество товара не совпадает",
	)

	require.Equal(
		s.T(),
		s.loginData.User.Id,
		gjson.Get(getBody, "user_id").String(),
		"Объявление принадлежит другому пользователю",
	)

	require.Len(
		s.T(),
		gjson.Get(getBody, "photos").Array(),
		1,
		"Количество фотографий не совпадает",
	)

	require.NotEmpty(
		s.T(),
		gjson.Get(getBody, "photos.0.url").String(),
		"URL фотографии отсутствует",
	)

	s.T().Log("Шаг 5: ищем объявление по полному названию")

	searchBody := advertisementsManager.SearchAdvertisements(
		s.T(),
		title,
		http.StatusOK,
	)

	items := gjson.Get(searchBody, "items").Array()

	found := false

	for _, item := range items {
		if item.Get("id").String() == id {
			require.Equal(s.T(), title, item.Get("title").String())
			found = true
			break
		}
	}

	require.Truef(
		s.T(),
		found,
		"Созданное объявление с ID %s не найдено в поиске",
		id,
	)

	s.T().Log("Шаг 7: получаем фотографии созданного объявления")

	photosBody := advertisementsManager.GetAdvertisementPhotos(
		s.T(),
		s.loginData.Token,
		id,
		http.StatusOK,
	)

	s.T().Log("Шаг 8: проверяем загруженную фотографию")

	photoItems := gjson.Parse(photosBody).Array()

	require.Len(
		s.T(),
		photoItems,
		1,
		"Ожидалась одна фотография в ответе GET",
	)

	require.NotEmpty(
		s.T(),
		photoItems[0].Get("id").String(),
		"ID фотографии отсутствует",
	)

	require.NotEmpty(
		s.T(),
		photoItems[0].Get("url").String(),
		"URL фотографии отсутствует",
	)

	require.Equal(
		s.T(),
		gjson.Get(body, "photos.0.id").String(),
		photoItems[0].Get("id").String(),
		"ID фотографии не совпадает с ответом POST",
	)

	require.Equal(
		s.T(),
		gjson.Get(body, "photos.0.url").String(),
		photoItems[0].Get("url").String(),
		"URL фотографии не совпадает с ответом POST",
	)
}
