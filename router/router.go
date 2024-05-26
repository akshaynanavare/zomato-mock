package router

import (
	"fmt"

	api "github.com/akshaynanavare/zomato-mock/api"
	"github.com/akshaynanavare/zomato-mock/endpoints/delivery"
	repository "github.com/akshaynanavare/zomato-mock/repository"

	"log"
	"net/http"
	"os"

	"strings"

	"github.com/julienschmidt/httprouter"
)

func NewServer() *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	return &http.Server{Addr: "localhost:" + port, Handler: newHandler()}
}

func addRoutes(r *httprouter.Router) {
	order := repository.NewOrder(defaultOrders())

	deliverPartner := repository.NewDeliveryPartner(defaultDeliveryPartners())

	deliveryHandler := delivery.NewHandler(delivery.NewService(order, deliverPartner))

	r.GET("/delivery/shortest-time/:deliveryPartner", deliveryHandler.GetMinimumTimeToDeliverAll)

}

func newHandler() http.Handler {
	r := httprouter.New()
	addRoutes(r)

	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Printf("panic: %+v", err)
		api.Error(w, r, fmt.Errorf("whoops! My handler has run into a panic"), http.StatusInternalServerError)
	}
	r.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Error(w, r, fmt.Errorf("we have OPTIONS for you but %v is not among them", r.Method), http.StatusMethodNotAllowed)
	})
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			sendBrowserDoc(w, r)
			return
		}

		api.Error(w, r, fmt.Errorf("whatever route you've been looking for, it's not here"), http.StatusNotFound)
	})

	return r
}

func sendBrowserDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusUnsupportedMediaType)

	b, err := os.ReadFile("browser.htm")
	if err != nil {
		api.Error(w, r, fmt.Errorf("read browser.htm failed: %w", err), http.StatusInternalServerError)
	}

	_, err = w.Write(b)
	if err != nil {
		api.Error(w, r, fmt.Errorf("send browser.htm failed: %w", err), http.StatusInternalServerError)
	}
}
