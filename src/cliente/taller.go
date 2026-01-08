/*
ESTE ES EL ÚNICO ARCHIVO QUE SE PUEDE MODIFICAR

RECOMENDACIÓN: Solo modicar a partir de la parte
				donde se encuentran la explicación
				de las otras variables.

*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	msg    string
)

func main() {
	var g *Garage
	var noExits [3]chan *Car // para fase 4 que no tiene canal de salida
	numCars := 30            // N coches totales
	numSlots := 10           // numPlazas
	numMechs := 4            // numMecanicos

	numWorkersDoc := numSlots
	numWorkersRep := numMechs // reparación limitada por mecánicos
	numWorkersClean := numSlots
	numWorkersDeliber := numSlots

	// Canales entre fases
	docChans := initPhaseChans()
	repChans := initPhaseChans()
	cleanChans := initPhaseChans()
	deliverChans := initPhaseChans()
	// Canal para events hacia el logger
	events := make(chan Event)

	// canales para terminar a las goroutines
	stopDoc := make(chan struct{})
	stopRep := make(chan struct{})
	stopClean := make(chan struct{})
	stopDeliver := make(chan struct{})
	finishChan := make(chan struct{})
	g = newGarage(numSlots)
	carspool := genCars(numCars)

	fmt.Printf("%-8s %-8s %-10s %-6s %-8s\n",
		"Tiempo[s]", "Coche[id]", "Incidencia", "Fase", "Estado")
	fmt.Printf("-------------------------------------------------\n")
	// Lanzar logger (único que imprime por stdout)
	go logManager(events)

	// Lanzar workers de cada fase

	// Fase de documentacion:
	startPhase(g, numWorkersDoc, docChans, repChans, events, DOCPHASE, stopDoc)

	// Fase de reparacion
	startPhase(g, numWorkersRep, repChans, cleanChans, events, REPAIRPHASE, stopRep)

	// Fase de limpieza
	startPhase(g, numWorkersClean, cleanChans, deliverChans, events, CLEANPHASE, stopClean)

	// Fase de entrega
	startPhase(g, numWorkersDeliber, deliverChans, noExits, events, DELIVERYPHASE, stopDeliver)

	go productor(g, carspool, numCars, docChans, finishChan)

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if n > 0 {
			msg = strings.TrimSpace(string(buf[:n]))
			sts, err := strconv.Atoi(msg)
			if err != nil {
				continue
			}
			if sts != 7 && sts != 8 {
				g.sts.Store(int32(sts))
			}
			select {
			case <-finishChan:
				freeThreats(stopDoc, stopRep, stopClean, stopDeliver)
				closeChans(docChans)
				closeChans(repChans)
				closeChans(cleanChans)
				closeChans(deliverChans)
				close(events)
				fmt.Println("Simulación completada. No hay más coches por atender.")
				return
			default:
			}
		}
	}
}
