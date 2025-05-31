package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type MinHeap []int

var _ heap.Interface = &MinHeap{}

func (p MinHeap) Len() int {
	return len(p)
}

func (p MinHeap) Less(i, j int) bool {
	return p[i] < p[j]
}
func (p *MinHeap) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}
func (p *MinHeap) Push(x any) {
	*p = append(*p, x.(int))
}
func (p *MinHeap) Pop() any {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

type ParkingLot struct {
	capacity  int
	available *MinHeap
	slotToCar map[int]string
	carToSlot map[string]int
}

func New(size int) *ParkingLot {
	h := &MinHeap{}

	for i := 1; i <= size; i++ {
		*h = append(*h, i)
	}
	heap.Init(h)
	return &ParkingLot{
		capacity:  size,
		available: h,
		slotToCar: make(map[int]string),
		carToSlot: make(map[string]int),
	}
}

func (p *ParkingLot) Park(carNumber string) {
	if p.available.Len() == 0 {
		fmt.Println("Sorry, parking lot is full")
		return
	}
	if _, exists := p.carToSlot[carNumber]; exists {
		fmt.Printf("Registration number %s already exists\n", carNumber)
		return
	}
	slot := heap.Pop(p.available).(int)
	p.slotToCar[slot] = carNumber
	p.carToSlot[carNumber] = slot

	fmt.Print(fmt.Sprintf("Allocated slot number: %d\n", slot))
}

func (p *ParkingLot) Leave(carNumber string, hours int) {
	slot, exists := p.carToSlot[carNumber]
	if !exists {
		fmt.Printf("Registration number %s not found\n", carNumber)
		return
	}
	delete(p.carToSlot, carNumber)
	delete(p.slotToCar, slot)
	heap.Push(p.available, slot)

	charge := 10
	if hours > 2 {
		charge += (hours - 2) * 10
	}
	fmt.Print(fmt.Sprintf("Registration number %s with Slot Number %d is free with Charge: $%d\n", carNumber, slot, charge))
}

func (p *ParkingLot) Status() {
	if p.available.Len() == p.capacity {
		log.Println("Parking lot is empty")
		return
	}
	fmt.Println("Slot No.\tRegistration No.")
	slots := make([]int, 0, len(p.slotToCar))
	for slot := range p.slotToCar {
		slots = append(slots, slot)
	}
	sort.Ints(slots)
	for _, slot := range slots {
		fmt.Printf("%d\t%s\n", slot, p.slotToCar[slot])
	}
}

func main() {
	var lot *ParkingLot
	if len(os.Args) < 2 {
		log.Fatal("Please provide the input file input.txt as a command line argument.")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "exist" {
			break
		}
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "create_parking_lot":
			n, _ := strconv.Atoi(parts[1])
			lot = New(n)
		case "park":
			lot.Park(parts[1])
		case "leave":
			hours, _ := strconv.Atoi(parts[2])
			lot.Leave(parts[1], hours)
		case "status":
			lot.Status()
		default:
			log.Printf("Command '%s' not recognized. Please use 'create_parking_lot', 'park', 'leave', or 'status'.\n", parts[0])
		}

	}
}
