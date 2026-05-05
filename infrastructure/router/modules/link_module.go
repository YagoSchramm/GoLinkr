package modules

import (
	"encoding/json"
	"net/http"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/domain/usecase/dtos"
	domainutil "github.com/YagoSchramm/Golinkr/domain/util"
	approuter "github.com/YagoSchramm/Golinkr/infrastructure/router"
	"github.com/gorilla/mux"
)

type linkModule struct {
	linkUsecase usecase.LinkUsecase
	name        string
	path        string
}

func NewLinkModule(linkUsecase usecase.LinkUsecase) approuter.Module {
	return linkModule{
		linkUsecase: linkUsecase,
		name:        "Link",
		path:        "/link",
	}
}

func (m linkModule) Name() string {
	return m.name
}

func (m linkModule) Path() string {
	return m.path
}

func (m linkModule) Routes() []approuter.RouteDefinition {
	return []approuter.RouteDefinition{
		{
			Path:        "",
			Description: "Create link",
			Handler:     m.create,
			HttpMethods: []string{http.MethodPost},
			Public:      true,
		},
		{
			Path:        "/{code}",
			Description: "Redirect to original URL",
			Handler:     m.redirect,
			HttpMethods: []string{http.MethodGet},
			Public:      true,
		},
	}
}

func (m linkModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}

func (m linkModule) create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateLinkDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		approuter.HandleError(w, derr.BadRequestError)
		return
	}

	link, err := m.linkUsecase.Create(r.Context(), entity.Link{
		OriginalURL: req.OriginalURL,
	})
	if err != nil {
		approuter.HandleError(w, err)
		return
	}

	if err := approuter.Write(w, domainutil.BuildLinkDTO(r, *link)); err != nil {
		approuter.HandleError(w, err)
	}
}

func (m linkModule) redirect(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	if code == "" {
		approuter.HandleError(w, derr.BadRequestError)
		return
	}

	link, err := m.linkUsecase.FindByCode(r.Context(), code)
	if err != nil {
		approuter.HandleError(w, err)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
