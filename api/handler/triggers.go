package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/moira-alert/moira-alert/api/controller"
	"github.com/moira-alert/moira-alert/api/dto"
	"net/http"
	"strings"
)

func triggers(router chi.Router) {
	router.Get("/", getAllTriggers)
	router.Put("/", saveTrigger)
	router.Get("/page", getTriggersPage)
	router.Route("/{triggerId}", trigger)
}

func getAllTriggers(writer http.ResponseWriter, request *http.Request) {
	//todo очень странная параша, отдает все триггреры
}

func saveTrigger(writer http.ResponseWriter, request *http.Request) {
	//Сохранение триггреа при его создании
}

func getTriggersPage(writer http.ResponseWriter, request *http.Request) {
	var onlyErrors bool
	var filterTags []string

	filterOkCookie, err := request.Cookie("moira_filter_ok")
	if err == http.ErrNoCookie {
		onlyErrors = false
	} else {
		onlyErrors = filterOkCookie.Value == "true"
	}

	fillerTagsCookie, err := request.Cookie("moira_filter_tags")
	if err == http.ErrNoCookie {
		filterTags = make([]string, 0)
	} else {
		str := strings.Split(fillerTagsCookie.Value, ",")
		for _, tag := range str {
			if tag != "" {
				filterTags = append(filterTags, tag)
			}
		}
	}

	page, size := getPageAndSize(request, 0, 10)

	triggersList, errorResponse := controller.GetTriggerPage(database, page, size, onlyErrors, filterTags)
	if errorResponse != nil {
		logger.Error(errorResponse.Err)
		render.Render(writer, request, errorResponse)
		return
	}

	if err := render.Render(writer, request, triggersList); err != nil {
		render.Render(writer, request, dto.ErrorRender(err))
		return
	}
}