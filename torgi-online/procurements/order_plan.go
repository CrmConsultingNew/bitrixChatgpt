package procurements

import (
	"fmt"
	"time"
)

func RunOrderPlanWithExistingStruct(fileName string) {
	start := time.Now()

	logTraceID := fmt.Sprintf("op%d-%d", time.Now().Unix(), time.Now().Nanosecond()%1000)
	logName := "order_plan"
	logMessage(logTraceID, "order_plan start", "i", logName)

	//orderPlan := NewOrderPlan(fileName)
	//orderPlan.Init()

	delta := time.Since(start)
	if delta.Seconds() > 30 {
		logMessage(logTraceID, fmt.Sprintf("order_plan: done with overload, execution time: %.2fs", delta.Seconds()), "w", logName)
	} else {
		logMessage(logTraceID, fmt.Sprintf("order_plan: done, execution time: %.2fs", delta.Seconds()), "i", logName)
	}
}

func logMessage(traceID, message, level, logName string) {
	log := fmt.Sprintf("[%s] %s: %s\n", level, traceID, message)
	fmt.Printf("%s", log) // Log is printed to the console
}
