package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers poll-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
  // this line is used by starport scaffolding
	r.HandleFunc("/poll/poll", listPollHandler(cliCtx, "poll")).Methods("GET")
	r.HandleFunc("/poll/poll", createPollHandler(cliCtx)).Methods("POST")
}
