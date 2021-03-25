package devices

import (
	"container/heap"
	"github.com/J00LZZ/go6502/pkg/bus"
	"log"
	"sync"
)

// from https://golang.org/pkg/container/heap/
type PriorityQueue []*InterruptItem
func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*InterruptItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *InterruptItem, id uint8, priority uint16) {
	item.id = id
	item.priority = priority
	heap.Fix(pq, item.index)
}


// The InterruptManager is a device to multiplex a number of interrupt sources over
// the single IRQ (or NMI) pin the 6502 has. This device is emulating an imaginary
// support chip on the motherboard which takes in a number of interrupt sources from
// other devices. When any one of them wants to trigger an interrupt, it instead signals
// this chip and it will interrupt the cpu. This device will then put information about
// the interrupt source in a predetermined location so software on the cpu can determine what the
// source of the interrupt was.
//
// # Interrupt manager memory layout:
//
// | Address | Usage |
// | --- | --- |
// | `Start+0` | Control register |
// | `Start+1` | When an interrupt fired, the cpu can query this address to see what's the source interrupt id |
// | `Start+2` | Write any nonzero value to this register to signify the interrupt was handled. This may immediately trigger another interrupt.|
// | `Start+3` | Software interrupt. Write a value n to this address to programmatically fire interrupt n |
// | `Start+4..Start+(n + 4)` | n descriptors representing the priority of each interrupt. Zero means the interrupt is completely inactive. |
//
//
// # Control register layout
// | Bit | Usage |
// | --- | --- |
// | `0`    | 0 to disable the interrupt manager completely |
// | `1`    | 1 to fire NMI instead of IRQ |
// | `2..7` | Unused |
//
type InterruptManager struct {
	RangeStart uint16
	NumInterrupts uint8

	// InterruptDescriptors is an array of numbers that describe that interrupt
	// 0 means the interrupt is inactive and will never actually interrupt the cpu
	// Any positive number represents the priority of this interrupt
	InterruptDescriptors []byte

	Active bool
	NMI bool

	// the last interrupt was handled
	handled bool

	// which interrupt the cpu is currently servicing
	currentInterruptId uint8

	// Queue of interrupts that occurred
	InterruptQueue PriorityQueue

	irq func()
	nmi func()

	sync.Mutex
}

type InterruptItem struct {
	id uint8
	priority uint16
	index int
}

func NewInterruptManager(start uint16, numInterrupts uint8) *InterruptManager {
	InterruptDescriptors := make([]byte, numInterrupts)
	for i := uint16(0); i < uint16(numInterrupts); i++ {
		InterruptDescriptors[i] = 0
	}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &InterruptManager{
		RangeStart:           start,
		NumInterrupts:        numInterrupts,
		InterruptDescriptors: InterruptDescriptors,
	}
}

func (i InterruptManager) Start() uint16 {
	return i.RangeStart
}

func (i InterruptManager) End() uint16 {
	return i.RangeStart + 3 + uint16(i.NumInterrupts)
}

func (i InterruptManager) LoadAddress(address uint16) byte {
	switch address {
	case i.RangeStart + 0:
		res := byte(0)
		if i.Active {
			res |= 0b00000001
		}
		if i.NMI {
			res |= 0b00000010
		}
		return res
	case i.RangeStart + 1:
		// Return the interrupt that has fired last
		return i.currentInterruptId
	case i.RangeStart + 2:
		// unreadable
		return 0
	case i.RangeStart + 3:
		// unreadable
		return 0
	default:
		if address < i.RangeStart + 4 + uint16(i.NumInterrupts) {
			index := address - i.RangeStart - 4
			return i.InterruptDescriptors[index]
		}
	}

	return 0
}

func (i InterruptManager) WriteAddress(address uint16, data byte) {
	switch address {
	case i.RangeStart + 0:
		if data & 0b00000001 != 0 {
			i.Active = true
		} else {
			i.Active = false
		}

		if data & 0b00000010 != 0 {
			i.NMI = true
		} else {
			i.NMI = false
		}

	case i.RangeStart + 1:
		// unwritable

	case i.RangeStart + 2:
		// finish interrupt
		i.handled = true

		i.FireNext()
	case i.RangeStart + 3:
		i.Interrupt(data)
	default:
		if address < i.RangeStart + 4 + uint16(i.NumInterrupts) {
			index := address - i.RangeStart - 4
			i.InterruptDescriptors[index] = data
		}
	}
}

func (i InterruptManager) GetName() string {
	return "Interrupt manager"
}

func (i InterruptManager) GetType() bus.Type {
	return bus.RW
}

func (i *InterruptManager) SetNMIFunc(nmi func()) {
	i.nmi = nmi
}

func (i *InterruptManager) SetIRQFunc(irq func()) {
	i.irq = irq
}

func (i *InterruptManager) Interrupt(id uint8) {
	if id < i.NumInterrupts {
		i.Lock()
		heap.Push(&i.InterruptQueue, &InterruptItem{
			id: id,
			priority: uint16(i.InterruptDescriptors[id]),
		})
		i.Unlock()

		i.FireNext()
	} else {
		log.Printf("interrupt fired with out-of-bounds id %d", id)
	}
}

func (i *InterruptManager) FireNext() {

	i.Lock()
	defer i.Unlock()

	if i.InterruptQueue.Len() > 0 && i.handled {

		next := heap.Pop(&i.InterruptQueue).(*InterruptItem)
		i.currentInterruptId = next.id

		i.handled = false

		if i.NMI {
			i.nmi()
		} else {
			i.irq()
		}
	}
}