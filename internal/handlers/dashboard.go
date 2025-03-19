package handlers

import (
	"log"
	"net/http"
	"time"

	"cryptodashboard/internal/pubsub"
	"cryptodashboard/internal/services"
	"cryptodashboard/internal/views"
)

type Context struct {
	BudaServices *services.Buda
	PubSub       *pubsub.PubSub
	Loc          *time.Location
}

func (c *Context) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Iniciando solicitud de dashboard desde %s", r.RemoteAddr)

	// // Suscribirse a las actualizaciones de balances
	// sub := c.PubSub.Subscribe("balances", func(balances []*services.Balance) {
	// 	log.Printf("Recibida actualización de %d balances", len(balances))
	// 	// Aquí puedes procesar los balances actualizados
	// 	// Por ejemplo, actualizar una base de datos o enviar notificaciones
	// })
	// defer c.PubSub.Unsubscribe("balances", sub)

	// balances, err := c.BudaServices.GetBalance()
	// if err != nil {
	// 	log.Printf("Error obteniendo balances: %v", err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// log.Printf("Balances obtenidos exitosamente: %d registros", len(balances))

	// cryptos := make([]models.Crypto, len(balances))
	// for i, balance := range balances {
	// 	cryptos[i] = models.Crypto{
	// 		ID:            1,
	// 		Name:          balance.ID,
	// 		Balance:       balance.Amount[0],
	// 		LastUpdatedAt: time.Now(),
	// 	}
	// }

	// data := models.DashboardData{
	// 	Crypto: cryptos,
	// }

	// Renderiza el template usando templ
	if err := views.Dashboard().Render(r.Context(), w); err != nil {
		log.Printf("Error renderizando el dashboard: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Dashboard renderizado exitosamente")
}

//init the loc

func (c *Context) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Iniciando solicitud de balance desde %s", r.RemoteAddr)

	balances, err := c.BudaServices.GetBalance()
	if err != nil {
		log.Printf("Error obteniendo balances: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Balances obtenidos exitosamente: %d registros", len(balances))
	//set timezone,
	now := time.Now().Add(time.Hour * -3)
	lastUpdatedAt := now.Format("2006-01-02 15:04:05")
	// Renderiza el template usando templ
	if err := views.BalanceTable(balances, lastUpdatedAt).Render(r.Context(), w); err != nil {
		log.Printf("Error renderizando el dashboard: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
