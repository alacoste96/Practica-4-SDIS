package main

import (
	"fmt"
	"time"
)

// ====================== Funciones que generan o son hilos ===========================

func productor(g *Garage, carspool map[int]*Car, numCars int, docChans [3]chan *Car, finish chan struct{}) {
	var mode int32
	var i int = 1

	// Generación de Cars (productor)
	for i = 0; i < numCars; i++ {
		// Esperar a tener una plaza libre
		mode = g.sts.Load()
		if mode == 0 || mode == 9 {
			time.Sleep(300 * time.Millisecond)
			i--
			continue
		}
		<-g.freeSlots
		c := getCarFromQ(carspool, g)
		if c == nil {
			g.freeSlots <- struct{}{}
			break
		}
		c.start = time.Now()
		// Registrar el Car en el mapa global
		g.signInCar(c)
		// Este Car hará todo el ciclo -> lo contamos en el WaitGroup
		g.wg.Add(1)

		// Entra en la primera fase (documentación)
		sendCar(docChans, c)
	}

	// Esperar a que TODOS los Cars terminen la fase de entrega
	g.wg.Wait()
	finish <- struct{}{}
}

// cada uno de estos worker gestiona 1 coche en su fase correspondiente
func worker(g *Garage, entrys [3]chan *Car, exits [3]chan *Car, events chan<- Event, phase int, stop <-chan struct{}) {
	var c *Car
	for {
		c = getCar(g, entrys, stop)
		if c == nil {
			return
		}
		g.updatePhase(c.id, phase)

		genEvent(events, c, "entra")
		time.Sleep(c.duration)
		genEvent(events, c, "sale")
		if phase == DELIVERYPHASE {
			time.Sleep(c.duration)
			g.delCar(c.id)
			// Liberar una plaza (podrá entrar otro Car en fase 1)
			g.freeSlots <- struct{}{}

			// Este Car ha terminado TODO el ciclo
			g.wg.Done()
		} else {
			sendCar(exits, c)
		}

	}
}

// funcion que genera workers para cada fase
func startPhase(g *Garage, nWorkers int, entrys [3]chan *Car, exits [3]chan *Car, events chan<- Event, phase int, stop <-chan struct{}) {
	for i := 0; i < nWorkers; i++ {
		go worker(g, entrys, exits, events, phase, stop)
	}
}

// ===================== Logger =====================

// este hilo es el único que tiene permitido escribir en stdout
func logManager(events <-chan Event) {
	for e := range events {
		secs := e.elapsed.Seconds()
		fmt.Printf("%-9.2f %-9d %-10s %-6d %-8s\n",
			secs,
			e.car,
			e.issue,
			e.phase,
			e.status,
		)

	}
}
