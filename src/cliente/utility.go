package main

import (
	"math/rand"
	"time"
)

// ===================== Utilidades =====================

// funcion que inicializa un garage
func newGarage(numSlots int) *Garage {
	g := &Garage{
		cars:      make(map[int]*Car),
		freeSlots: make(chan struct{}, numSlots),
	}
	for i := 0; i < numSlots; i++ {
		g.freeSlots <- struct{}{}
	}
	return g
}

func getCarFromM(cars map[int]*Car, issue IssueType) *Car {

	for id, c := range cars {
		if c.issue == issue {
			delete(cars, id)
			return c
		}
	}
	return nil
}

func setPrio(mode int32) (IssueType, IssueType, IssueType) {
	var high, medium, low IssueType

	// Valores por defecto:
	high = MECH       // A
	medium = ELECTRIC // B
	low = BODY        // C
	switch mode {
	case 5, 2:
		high = ELECTRIC
		medium = MECH
	case 6, 3:
		high = BODY
		medium = MECH
		low = ELECTRIC
	}
	return high, medium, low
}

func getCarFromQ(cars map[int]*Car, g *Garage) *Car {
	if len(cars) == 0 {
		return nil
	}

	for {
		mode := g.sts.Load()
		switch mode {
		case 1, 2, 3: // restrictivos
			high, _, _ := setPrio(mode)
			if c := getCarFromM(cars, high); c != nil {
				return c
			}
			time.Sleep(300 * time.Millisecond)
		case 4, 5, 6: // prioridades
			high, medium, low := setPrio(mode)
			if c := getCarFromM(cars, high); c != nil {
				return c
			}
			if c := getCarFromM(cars, medium); c != nil {
				return c
			}
			return getCarFromM(cars, low)
		// Otros modos (0, 9, lo que sea): espero a que se ponga en algo útil
		default:
			time.Sleep(300 * time.Millisecond)
		}
		if len(cars) == 0 {
			return nil
		}
	}
}

// genera ncars coches en un mapa con prios aleatorias
func genCars(ncars int) map[int]*Car {
	var i int

	cars := make(map[int]*Car, ncars)
	for i = 0; i < ncars; i++ {
		cars[i] = genCar(i)
	}
	return cars
}

// saca un número float aleatorio entre 0 y 2 con 1 decimal
func randDecimal() float64 {
	n := rand.Intn(21) // 0..20
	return float64(n) / 10.0
}

// genera un coche con prio aleatoria
func genCar(id int) *Car {
	var issue IssueType
	var duration time.Duration
	var lag float64
	var prio int = rand.Intn(3)

	lag = randDecimal()
	switch prio {
	case 0: // mecánica
		issue = MECH
		duration = time.Duration((5 + lag) * float64(time.Second))
	case 1: // eléctrica
		issue = ELECTRIC
		duration = time.Duration((3 + lag) * float64(time.Second))
	case 2: // carroceria
		issue = BODY
		duration = time.Duration((1 + lag) * float64(time.Second))
	}

	return &Car{
		id:       id,
		issue:    issue,
		duration: duration,
		curphase: 0,
	}
}

func freeThreats(docChan, repChan, cleanChan, deliverChan chan struct{}) {
	close(docChan)
	close(repChan)
	close(cleanChan)
	close(deliverChan)
}

// cierra un slice de 3 canales
func closeChans(chans [3]chan *Car) {
	for i := range chans {
		close(chans[i])
	}
}

// obtiene un coche a traves de un canal de comunicacion entre hilos, con prioridades
func getCar(g *Garage, chans [3]chan *Car, stop <-chan struct{}) *Car {
	var c *Car
	var mode int32
	var high, medium, low int32

	mode = g.sts.Load()
	high = 0
	medium = 1
	low = 2
	if mode >= 4 && mode <= 6 {
		high = mode - 4
		if mode != 4 {
			medium = 0
		}
		if mode == 6 {
			low = 1
		}
	}
	select {
	case c = <-chans[high]: // pillo el de alta prio
	default:
		select {
		case c = <-chans[medium]: // pillo el de prio media
		default:
			select { // pillo el de baja prio o espero por el primero que llegue
			case c = <-chans[high]:
			case c = <-chans[medium]:
			case c = <-chans[low]:
			case <-stop:
				return nil
			}
		}
	}
	return c
}

// envía un coche a traves de un canal en funcion de su prioridad
func sendCar(chans [3]chan *Car, c *Car) {
	switch c.issue {
	case MECH:
		chans[0] <- c // prio alta
	case ELECTRIC:
		chans[1] <- c // prio media
	case BODY:
		chans[2] <- c // prio baja
	}
}

// inicializa los canales de una fase (1 por cada prioridad)
// indice 0 = prio alta. indice 1 = prio media. indice 2 = prio baja
func initPhaseChans() [3]chan *Car {
	var chans [3]chan *Car

	for i := range chans {
		chans[i] = make(chan *Car)
	}
	return chans
}

// genera un evento que se envía al logger para que este escriba por stdout
func genEvent(events chan<- Event, c *Car, sts string) {
	events <- Event{
		elapsed: time.Since(c.start),
		car:     c.id,
		phase:   c.curphase,
		status:  sts,
		issue:   c.issue,
	}
}
