package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/samsarahq/thunder/graphql"
	userService "kroseida.org/slixx/internal/supervisor/service/user"
	"net/http"
	"sync"

	"github.com/samsarahq/thunder/batch"
	"github.com/samsarahq/thunder/reactive"
)

func HTTPHandler(schema *graphql.Schema, middlewares ...graphql.MiddlewareFunc) http.Handler {
	return &httpHandler{
		schema:      schema,
		middlewares: middlewares,
	}
}

type httpHandler struct {
	schema      *graphql.Schema
	middlewares []graphql.MiddlewareFunc
}

type httpPostBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type httpResponse struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeResponse := func(value interface{}, err error) {
		response := httpResponse{}
		if err != nil {
			response.Errors = []string{err.Error()}
		} else {
			response.Data = value
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Error(w, string(responseJSON), http.StatusOK)
	}

	if r.Method != "POST" {
		writeResponse(nil, errors.New("request must be a POST"))
		return
	}

	if r.Body == nil {
		writeResponse(nil, errors.New("request must include a query"))
		return
	}

	var params httpPostBody
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeResponse(nil, err)
		return
	}

	query, err := graphql.Parse(params.Query, params.Variables)
	if err != nil {
		writeResponse(nil, err)
		return
	}

	schema := h.schema.Query
	if query.Kind == "mutation" {
		schema = h.schema.Mutation
	}

	if err := graphql.PrepareQuery(schema, query.SelectionSet); err != nil {
		writeResponse(nil, err)
		return
	}

	var wg sync.WaitGroup
	e := graphql.Executor{}

	wg.Add(1)
	runner := reactive.NewRerunner(r.Context(), func(ctx context.Context) (interface{}, error) {
		defer wg.Done()

		ctx = batch.WithBatching(ctx)

		if r.Header.Get("authorization") != "" {
			userID, err := userService.GetUserBySession(r.Header.Get("authorization"))
			if err != nil {
				ctx = context.WithValue(ctx, "user", nil)
			}
			user, err := userService.Get(userID)
			if err != nil {
				ctx = context.WithValue(ctx, "user", nil)
			}
			ctx = context.WithValue(ctx, "user", user)
		}

		var middlewares []graphql.MiddlewareFunc
		middlewares = append(middlewares, h.middlewares...)
		middlewares = append(middlewares, func(input *graphql.ComputationInput, next graphql.MiddlewareNextFunc) *graphql.ComputationOutput {
			output := next(input)
			output.Current, output.Error = e.Execute(input.Ctx, schema, nil, input.ParsedQuery)
			return output
		})

		output := graphql.RunMiddlewares(middlewares, &graphql.ComputationInput{
			Ctx:         ctx,
			ParsedQuery: query,
			Query:       params.Query,
			Variables:   params.Variables,
		})

		current, err := output.Current, output.Error

		if err != nil {
			if graphql.ErrorCause(err) == context.Canceled {
				return nil, err
			}

			writeResponse(nil, err)
			return nil, err
		}

		writeResponse(current, nil)
		return nil, nil
	}, graphql.DefaultMinRerunInterval)

	wg.Wait()
	runner.Stop()
}
