package http

import (
	"context"
	"github.com/oskov/dictionary-service/internal/api/http/oapi"
	"github.com/oskov/dictionary-service/internal/application"
)

var _ oapi.StrictServerInterface = &API{}

type API struct {
	app application.App
}

func NewAPI(app application.App) *API {
	return &API{
		app: app,
	}
}

func (api *API) GetWordWord(
	ctx context.Context,
	request oapi.GetWordWordRequestObject,
) (oapi.GetWordWordResponseObject, error) {
	if len(request.Word) < 1 {
		return oapi.GetWordWord400Response{}, nil
	}

	word, err := api.app.Core.GetWord(request.Word)
	if err != nil {
		return oapi.GetWordWord500Response{}, err
	}

	resp := oapi.GetWordWord200JSONResponse{
		Word: word.Word,
	}

	for _, definition := range word.Definitions {
		resp.Definitions = append(resp.Definitions, oapi.GetWordResultDefinition{
			Definition: definition.Definition,
			Examples:   definition.Examples,
		})
	}

	return resp, nil
}
