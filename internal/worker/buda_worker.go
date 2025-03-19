package worker

import (
	"cryptodashboard/internal/pubsub"
	"cryptodashboard/internal/services"
	"log"
	"time"
)

type BudaWorker struct {
	buda   *services.Buda
	pubsub *pubsub.PubSub
	stopCh chan struct{}
	ticker *time.Ticker
}

func NewBudaWorker(buda *services.Buda, pubsub *pubsub.PubSub, interval time.Duration) *BudaWorker {
	return &BudaWorker{
		buda:   buda,
		pubsub: pubsub,
		stopCh: make(chan struct{}),
		ticker: time.NewTicker(interval),
	}
}

func (w *BudaWorker) Start() {
	log.Println("Iniciando BudaWorker...")
	go w.run()
}

func (w *BudaWorker) Stop() {
	log.Println("Deteniendo BudaWorker...")
	w.ticker.Stop()
	close(w.stopCh)
}

func (w *BudaWorker) run() {
	for {
		select {
		case <-w.stopCh:
			return
		case <-w.ticker.C:
			balances, err := w.buda.GetBalance()
			if err != nil {
				log.Printf("Error obteniendo balances: %v", err)
				continue
			}

			log.Printf("Publicando %d balances", len(balances))
			w.pubsub.Publish("balances", balances)
		}
	}
}
