package main

import (
  "fmt"
  "sync"
  "time"
  "math/rand"
)

type ChopS struct{ sync.Mutex }	// chopstick held by one at a time

type Philo struct {
  badge int		// philosopher's id
  stick []*ChopS	// left and right chopstick
}

const pop = 5		// population
const sit = 3		// sittings
const plates = 2	// simultaneous eaters

func (p Philo) eat(table chan int, wg *sync.WaitGroup) {
  for b:=0; b<sit; b++ {	// each p eats sit times
    side := rand.Intn(2)	// random 0 or 1

    table <- p.badge		// this p wants to eat
    p.stick[side].Lock()	// get left chopstick
    p.stick[1-side].Lock()	// get right chopstick

    fmt.Println("starting to eat", p.badge)
    fmt.Println("finishing eating", p.badge)

    p.stick[1-side].Unlock()	// done
    p.stick[side].Unlock()
  }
  wg.Done()		// everyone has eaten thrice
}

var line int		// ~ status at the table
var eating []int	// currently eating at the table

func host(table chan int) {		// host to run in thread
  for badge := range table {		// ids via channel
    eating = append(eating, badge)	// ids at the table
    if len(eating) > plates {		// only two at a time
      eating = eating[1:]
    }
    line++
    fmt.Println(line, eating, "hosted")
  }
}

func main() {
  rand.Seed(time.Now().UnixNano())

  CSticks := make([]*ChopS, pop)
  for i := 0; i < pop; i++ {
    CSticks[i] = new(ChopS)
  }

  philos := make([]*Philo, pop)
  for i := 0; i < pop; i++ {
    philos[i] = &Philo{i+1,[]*ChopS{CSticks[i],CSticks[(i+1)%5]}}
  }

  table := make(chan int, plates)
  var wg sync.WaitGroup		// setup to wait
  wg.Add(pop)			// for all eaters

  go host(table)		// start thread using channel

  for i := 0; i < pop; i++ {
    go philos[i].eat(table, &wg)
  }
  wg.Wait()
  close(table)
}

