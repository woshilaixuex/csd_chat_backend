package ticker

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: The timer is implemented by using the time wheel,referred to godis`s TimeWheel
 * @Date: 2025-05-16 14:34
 */
type taskInfo struct {
	slotIndex int
	etask     *list.Element
}

type TimeWheel struct {
	interval   time.Duration
	ticker     *time.Ticker
	currentPos int
	slotNum    int
	slots      []*list.List

	taskInfos   map[string]*taskInfo
	addTaskC    chan *task
	removeTaskC chan string
	stopChannel chan struct{}

	sync.RWMutex
	sync.Once
}

type task struct {
	t time.Duration
	// The number of rotations needed
	circle     int
	key        string
	isLongTerm bool
	job        func()
}

func NewTimeWheel(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:    interval,
		currentPos:  0,
		slotNum:     slotNum,
		slots:       make([]*list.List, slotNum),
		taskInfos:   make(map[string]*taskInfo),
		addTaskC:    make(chan *task),
		removeTaskC: make(chan string),
		stopChannel: make(chan struct{}),
	}
	tw.initSlots()
	return tw
}

func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

func (tw *TimeWheel) Start() *TimeWheel {
	start := func() {
		tw.ticker = time.NewTicker(tw.interval)
		go tw.start()
	}
	tw.Once.Do(start)
	return tw
}

func (tw *TimeWheel) Stop() {
	tw.stopChannel <- struct{}{}
}

func (tw *TimeWheel) AddTaskWithKey(key string, job func(), t time.Duration, isLongTerm bool) {
	tw.addTaskC <- &task{
		t:          t,
		key:        key,
		isLongTerm: isLongTerm,
		job:        job,
	}
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskC:
			tw.addTaskHandler(task)
		case key := <-tw.removeTaskC:
			tw.removeTaskHandler(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	tw.Lock()
	list := tw.slots[tw.currentPos]
	tw.currentPos = (tw.currentPos + 1) % tw.slotNum
	tw.Unlock()
	go tw.scanAndRunTask(list)
}

func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	var tasksToRemove []*struct {
		key     string
		element *list.Element
	}
	var tasksToResave []*task
	tw.RLock()

	for e := l.Front(); e != nil; {
		etask := e.Value.(*task)
		if etask.circle > 0 {
			etask.circle--
			e = e.Next()
			continue
		}

		go func(job func()) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			job()
		}(etask.job)
		next := e.Next()
		tasksToRemove = append(tasksToRemove, &struct {
			key     string
			element *list.Element
		}{
			key:     etask.key,
			element: e,
		})
		if etask.isLongTerm {
			tasksToResave = append(tasksToResave, etask)
		}
		e = next
	}
	tw.RUnlock()

	tw.Lock()
	defer tw.Unlock()
	for _, remove := range tasksToRemove {
		l.Remove(remove.element)
		delete(tw.taskInfos, remove.key)
	}
	for _, resave := range tasksToResave {
		circle, pos := tw.getCircleAndPos(resave)
		resave.circle = circle
		etask := tw.slots[pos].PushBack(resave)
		info := &taskInfo{
			slotIndex: pos,
			etask:     etask,
		}
		tw.taskInfos[resave.key] = info
	}
}

// The function is add handler
//
// If the key is empty, the task cannot be deleted
func (tw *TimeWheel) addTaskHandler(task *task) {
	circle, pos := tw.getCircleAndPos(task)
	task.circle = circle

	if task.key == "" {
		fmt.Println("task key cannot null")
		return
	}
	if _, ok := tw.taskInfos[task.key]; ok {
		tw.removeTaskWithKey(task.key)
	}
	tw.Lock()
	defer tw.Unlock()

	etask := tw.slots[pos].PushBack(task)

	info := &taskInfo{
		slotIndex: pos,
		etask:     etask,
	}
	tw.taskInfos[task.key] = info
}

// The function return circle num and task`s position in time wheel
func (tw *TimeWheel) getCircleAndPos(task *task) (circle int, pos int) {
	ticks := int(task.t / tw.interval)
	circle = ticks / tw.slotNum
	pos = (tw.currentPos + ticks) % tw.slotNum
	return circle, pos
}

func (tw *TimeWheel) removeTaskHandler(key string) {
	tw.Lock()
	defer tw.Unlock()
	tw.removeTaskWithKey(key)
}

func (tw *TimeWheel) removeTaskWithKey(key string) {
	info, ok := tw.taskInfos[key]
	if !ok {
		return
	}
	list := tw.slots[info.slotIndex]

	list.Remove(info.etask)
	delete(tw.taskInfos, key)
}
