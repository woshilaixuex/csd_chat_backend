package ticker

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: The timer
 * @Date: 2025-02-28 21:11
 */

type Ticker interface {
	AddTask(job func())
	AddTaskWithKey(key string, job func())
	StopTask(key string)
}
