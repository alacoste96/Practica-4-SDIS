package main

import (
	"math/rand"
	"testing"
	"time"
)

type TimeConfig struct {
	name       string
	a, b, c    int // coches A/B/C
	slots      int // NumPlazas
	mechanics  int // NumMecanicos
	mode       int32
	iterations int
}

type Stats struct {
	n    int
	sum  time.Duration
	min  time.Duration
	max  time.Duration
	mean float64
}

func (s *Stats) add(x time.Duration) {
	s.n++
	s.sum += x
	if s.n == 1 || x < s.min {
		s.min = x
	}
	if s.n == 1 || x > s.max {
		s.max = x
	}
	xf := float64(x)
	delta := xf - s.mean
	s.mean += delta / float64(s.n)
}

func (s Stats) getMean() time.Duration {
	if s.n == 0 {
		return 0
	}
	return time.Duration(s.mean)
}

// 1 "segundo" simulado = scale (ej: 100ms)
func genCarsByCounts(a, b, c int, scale time.Duration, seed int64) map[int]*Car {
	rng := rand.New(rand.NewSource(seed))
	total := a + b + c
	cars := make(map[int]*Car, total)

	id := 0
	for i := 0; i < a; i++ {
		cars[id] = newCarWithLag(id, MECH, scale, rng) // A = mecánica
		id++
	}
	for i := 0; i < b; i++ {
		cars[id] = newCarWithLag(id, ELECTRIC, scale, rng) // B = eléctrica
		id++
	}
	for i := 0; i < c; i++ {
		cars[id] = newCarWithLag(id, BODY, scale, rng) // C = carrocería
		id++
	}
	return cars
}

// base 5/3/1 + variación aleatoria 0..2 (en décimas).
func newCarWithLag(id int, issue IssueType, scale time.Duration, rng *rand.Rand) *Car {
	lagSteps := rng.Intn(21) // 0..20 décimas
	lag := time.Duration(lagSteps) * scale / 10

	var baseSeconds int
	switch issue {
	case MECH:
		baseSeconds = 5
	case ELECTRIC:
		baseSeconds = 3
	default:
		baseSeconds = 1
	}

	return &Car{
		id:       id,
		issue:    issue,
		duration: time.Duration(baseSeconds)*scale + lag,
		curphase: 0,
	}
}

func noEvents(events chan Event) {
	for range events {
	}
}

func runSimulationOnce(cfg TimeConfig, scale time.Duration, seed int64) time.Duration {
	// Setup garage y estado
	g := newGarage(cfg.slots)
	g.sts.Store(cfg.mode)

	numCars := cfg.a + cfg.b + cfg.c
	carspool := genCarsByCounts(cfg.a, cfg.b, cfg.c, scale, seed)

	// Canales entre fases
	docChans := initPhaseChans()
	repChans := initPhaseChans()
	cleanChans := initPhaseChans()
	deliverChans := initPhaseChans()
	var noExits [3]chan *Car

	// Canal de eventos: lo silenciamos
	events := make(chan Event)
	go noEvents(events)

	// Canales stop
	stopDoc := make(chan struct{})
	stopRep := make(chan struct{})
	stopClean := make(chan struct{})
	stopDeliver := make(chan struct{})

	finish := make(chan struct{}, 1)

	// Workers por fase:
	startPhase(g, cfg.slots, docChans, repChans, events, DOCPHASE, stopDoc)
	startPhase(g, cfg.mechanics, repChans, cleanChans, events, REPAIRPHASE, stopRep)
	startPhase(g, cfg.slots, cleanChans, deliverChans, events, CLEANPHASE, stopClean)
	startPhase(g, cfg.slots, deliverChans, noExits, events, DELIVERYPHASE, stopDeliver)

	start := time.Now()
	go productor(g, carspool, numCars, docChans, finish)

	select {
	case <-finish:
	}

	elapsed := time.Since(start)

	close(stopDoc)
	close(stopRep)
	close(stopClean)
	close(stopDeliver)

	closeChans(docChans)
	closeChans(repChans)
	closeChans(cleanChans)
	closeChans(deliverChans)

	close(events)
	return elapsed
}

func TestComparativasTiempo(t *testing.T) {
	const scale = 100 * time.Millisecond

	// Modo fijo para hacer comparable (4 = prioridad categoría A)
	const mode int32 = 4

	configs := []TimeConfig{
		{name: "T1_A10_B10_C10_P6_M3", a: 10, b: 10, c: 10, slots: 6, mechanics: 3, mode: mode, iterations: 10},
		{name: "T2_A20_B5_C5_P6_M3", a: 20, b: 5, c: 5, slots: 6, mechanics: 3, mode: mode, iterations: 10},
		{name: "T3_A5_B5_C20_P6_M3", a: 5, b: 5, c: 20, slots: 6, mechanics: 3, mode: mode, iterations: 10},

		{name: "T4_A10_B10_C10_P4_M4", a: 10, b: 10, c: 10, slots: 4, mechanics: 4, mode: mode, iterations: 10},
		{name: "T5_A20_B5_C5_P4_M4", a: 20, b: 5, c: 5, slots: 4, mechanics: 4, mode: mode, iterations: 10},
		{name: "T6_A5_B5_C20_P4_M4", a: 5, b: 5, c: 20, slots: 4, mechanics: 4, mode: mode, iterations: 10},
	}

	for _, cfg := range configs {
		cfg := cfg
		t.Run(cfg.name, func(t *testing.T) {
			var st Stats
			totalCars := cfg.a + cfg.b + cfg.c

			for i := 0; i < cfg.iterations; i++ {
				seed := int64(1000 + i) // cambia la variación de tiempos por iteración
				d := runSimulationOnce(cfg, scale, seed)
				st.add(d)
			}

			mean := st.getMean()
			perCar := time.Duration(0)
			if totalCars > 0 {
				perCar = mean / time.Duration(totalCars)
			}

			t.Logf("Iteraciones: %d | Coches: %d (A=%d,B=%d,C=%d) | Plazas=%d Mecánicos=%d | Modo=%d",
				st.n, totalCars, cfg.a, cfg.b, cfg.c, cfg.slots, cfg.mechanics, cfg.mode,
			)
			t.Logf("Total medio: %v | min=%v max=%v | medio por coche (aprox): %v",
				mean, st.min, st.max, perCar,
			)
		})
	}
}
