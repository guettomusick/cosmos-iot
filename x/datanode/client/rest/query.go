package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/datanode/{address}/records/{channelid}/{date}", queryRecordsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/datanode/{address}", queryDataNodeHandler(cliCtx)).Methods("GET")
}

func queryDataNodeHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/datanode/datanode/%s", address), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRecordsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]
		channelID := vars["channelid"]
		date := vars["date"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/datanode/records/%s/%s/%s", address, channelID, date), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
