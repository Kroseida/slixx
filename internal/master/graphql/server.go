package graphql

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"kroseida.org/slixx/internal/master/application"
	"kroseida.org/slixx/internal/master/datasource"
	"net/http"
	"time"
)

func schema() *graphql.Schema {
	builder := schemabuilder.NewSchema()
	registerQuery(builder)
	registerMutation(builder)
	registerObject(builder)
	return builder.MustBuild()
}

type simpleLogger struct {
}

func (s *simpleLogger) StartExecution(ctx context.Context, tags map[string]string, initial bool) {

}
func (s *simpleLogger) FinishExecution(ctx context.Context, tags map[string]string, delay time.Duration) {

}
func (s *simpleLogger) Error(ctx context.Context, err error, tags map[string]string) {
	application.Logger.Error(err)
}

func handler(schema *graphql.Schema) http.Handler {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := http.Header{}

		header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		header.Set("Access-Control-Allow-Origin", application.CurrentSettings.Http.AllowedOrigin)

		socket, err := upgrader.Upgrade(w, r, header)
		if err != nil {
			application.Logger.Debug("upgrader.Upgrade: %v", err)
			return
		}
		defer socket.Close()

		context := func(ctx context.Context) context.Context {
			userId, err := datasource.UserProvider.GetUserBySession(r.URL.Query().Get("authorization"))
			if err != nil {
				return context.WithValue(ctx, "user", nil)
			}
			user, err := datasource.UserProvider.GetUser(userId)
			if err != nil {
				return context.WithValue(ctx, "user", nil)
			}
			return context.WithValue(ctx, "user", user)
		}

		graphql.ServeJSONSocket(r.Context(), socket, schema, context, &simpleLogger{})
	})
}

func Start() error {
	schema := schema()
	introspection.AddIntrospectionToSchema(schema)

	// Expose schema and graphiql.
	http.Handle("/graphql", handler(schema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	return http.ListenAndServe(":3030", nil)
}
