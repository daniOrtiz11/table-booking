package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniOrtiz11/table-booking/pkg/bill"
	"github.com/daniOrtiz11/table-booking/pkg/tables"

	"github.com/daniOrtiz11/table-booking/pkg/booking"

	"github.com/daniOrtiz11/table-booking/internal/utils"
	"github.com/daniOrtiz11/table-booking/pkg/locate"
	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler

	locate  locate.Service
	booking booking.Service
	bill    bill.Service
	tables  tables.Service
}

/*
Server is a interface to define the methods
*/
type Server interface {
	Router() http.Handler
	Addr() string

	locateRequest(w http.ResponseWriter, r *http.Request)
	healthcheckRequest(w http.ResponseWriter, r *http.Request)
	bookingRequest(w http.ResponseWriter, r *http.Request)
	tablesRequest(w http.ResponseWriter, r *http.Request)
	billRequest(w http.ResponseWriter, r *http.Request)
}

func (a *api) Router() http.Handler {
	return a.router
}

func (a *api) Addr() string {
	return fmt.Sprintf("%s:%s", utils.GetEnv("SERVER_HOST", "0.0.0.0"), utils.GetEnv("SERVER_PORT", "9091"))
}

/*
New will retrieve a Service interface and define its attributes
*/
func New() Server {
	a := &api{}
	r := mux.NewRouter()
	/*
		r := mux.NewRouter()
		api := r.PathPrefix("/api/v1").Subrouter()
	*/
	r.HandleFunc("/healthcheck", a.healthcheckRequest).Methods(http.MethodGet)
	r.HandleFunc("/booking", a.bookingRequest).Methods(http.MethodPost)
	r.HandleFunc("/bill", a.billRequest).Methods(http.MethodPost)
	r.HandleFunc("/locate", a.locateRequest).Methods(http.MethodPost)
	r.HandleFunc("/tables", a.tablesRequest).Methods(http.MethodPut)
	a.router = r

	return a
}

func (a *api) locateRequest(w http.ResponseWriter, r *http.Request) {

	contentType := utils.GetContentType(r)
	accept := utils.GetAccept(r)
	if (contentType != "application/x-www-form-urlencoded") || (accept != "application/json") {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		r.ParseForm()
		id, errArg := strconv.Atoi(r.FormValue("ID"))
		if errArg != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status, response := locate.ServiceImpl(id)
		w.WriteHeader(status)
		if response != 0 {
			json.NewEncoder(w).Encode(response)
		}
	}
}

func (a *api) healthcheckRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *api) bookingRequest(w http.ResponseWriter, r *http.Request) {
	contentType := utils.GetContentType(r)
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		defer r.Body.Close()
		body, errBody := ioutil.ReadAll(r.Body)
		if errBody != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status := booking.ServiceImpl(body)
		w.WriteHeader(status)
		return
	}

}

func (a *api) billRequest(w http.ResponseWriter, r *http.Request) {
	contentType := utils.GetContentType(r)
	firstContentType := strings.Split(contentType, ";")
	if firstContentType[0] != "multipart/form-data" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		err := r.ParseMultipartForm(1024 * 1024 * 16)
		if err != nil {

		}
		mapsValue := r.MultipartForm.Value
		idArg := mapsValue["ID"]
		if len(idArg) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, errArg := strconv.Atoi(idArg[0])
		if errArg != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status := bill.ServiceImpl(id)
		w.WriteHeader(status)
	}
}

func (a *api) tablesRequest(w http.ResponseWriter, r *http.Request) {
	contentType := utils.GetContentType(r)
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		defer r.Body.Close()
		body, errBody := ioutil.ReadAll(r.Body)
		if errBody != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status := tables.ServiceImpl(body)
		w.WriteHeader(status)
		return
	}
}
