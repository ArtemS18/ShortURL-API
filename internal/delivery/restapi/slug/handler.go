package slug

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemS18/ShortURL-API/internal/delivery"
	"github.com/ArtemS18/ShortURL-API/internal/delivery/restapi/utils"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/ArtemS18/ShortURL-API/pkg/ctxLogger"
	"github.com/gorilla/mux"
)

type SlugHandler struct {
	uc delivery.SlugUseCase
}

func NewSlugHandler(uc delivery.SlugUseCase) *SlugHandler {
	return &SlugHandler{
		uc: uc,
	}
}

// @Summary Создание короткой ссылки
// @Description Возвращает короткую ссылку для заданного URL
// @Tags slug
// @Produce json
// @Param request body dto.CreateSlugRequest true "URL для сокращения"
// @Success 200 {object} dto.CreateSlugResponse "Успешный ответ при нахождении соответстующей короткой ссылки в бд"
// @Success 201 {object} dto.CreateSlugResponse "Успешный ответ при создании новой короткой ссылки"
// @Failure 500 {object} utils.BaseErrorResponse "Ошибка сервера"
// @Router /slugs [post]
func (h *SlugHandler) CreateSlugHandler(w http.ResponseWriter, r *http.Request) {
	op := "SlugHandler.CreateSlugHandler"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)
	var req dto.CreateSlugRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("json.NewDecoder: %v", err)
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		log.Error("url is empty")
		utils.WriteError(w, "url is required", http.StatusBadRequest)
		return
	}
	resp, err := h.uc.CreateSlug(r.Context(), &req)
	if err != nil {
		log.Errorf("h.uc.CreateSlugDB: %v", err)
		utils.HandelError(w, err)
		return
	}
	if resp.IsCreated {
		utils.JSONResponse(w, http.StatusCreated, resp)
		return
	}
	utils.JSONResponse(w, http.StatusOK, resp)

}

// @Summary Получение URL по короткой ссылке
// @Description Возвращает оригинальный URL для заданной короткой ссылки
// @Tags slug
// @Produce json
// @Param slug path string true "Короткая ссылка""
// @Param redirect query bool true "Если указан true, то будет выполнин редирект на URL соответсвующий вводимой короткой ссылки"" default(true)
// @Success 200 {object} dto.GetURLResponse "Успешный ответ, если выбран redirect=false"
// @Failure 404 {object} utils.BaseErrorResponse "URL не найден"
// @Failure 500 {object} utils.BaseErrorResponse "Ошибка сервера"
// @Router /{slug} [get]
func (h *SlugHandler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	op := "SlugHandler.GetURLHandler"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	redirect := r.URL.Query().Get("redirect")

	vars := mux.Vars(r)

	slug := vars["slug"]
	if slug == "" {
		log.Errorf("Slug not found in path")
		utils.WriteError(w, "slug is requered", http.StatusBadRequest)
		return
	}
	req := dto.GetURLRequest{
		SlugURL: slug,
	}
	resp, err := h.uc.GetURL(r.Context(), &req)
	if err != nil {
		log.Errorf("h.uc.GetURL: %v", err)
		utils.HandelError(w, err)
		return
	}
	if redirect != "false" {
		http.Redirect(w, r, resp.URL, http.StatusMovedPermanently)
		return
	}
	utils.JSONResponse(w, http.StatusOK, resp)
}
