package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/vanya-egorov/url_shortener.git/internal/lib/api/response"
	"github.com/vanya-egorov/url_shortener.git/internal/lib/logger/sl"
	"github.com/vanya-egorov/url_shortener.git/internal/lib/random"
	"github.com/vanya-egorov/url_shortener.git/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required, url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

// TODO move to config
const aliasLength = 10

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("failed to validate request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, resp.Error("url already exist"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}
		log.Info("url saved", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
