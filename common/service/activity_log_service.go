package service

import (
	"context"
	"errors"
	"fmt"
	"pi-inventory/common/models"
	"sync/atomic"
)

func NewActivityLogService(handler ActivityLogHandlerInterface) ActivityLogServiceInterface {
	if activityLogSvcObj == nil {
		// initiate(handler)
		x := int32(0)
		numberOfWorker := 2
		bufferSize := 100
		activityLogSvcObj = &ActivityLogService{
			config: ActivityLogConfig{
				BufferSize:                  bufferSize,
				NumberOfWorker:              numberOfWorker,
				ShutdownGracePeriodInSecond: 2,
			},
			activityLogBuffer:    make(chan models.ActivityLog, bufferSize),
			handler:              handler,
			workerStopperChannel: make(chan bool),
			shuttingDownFlag:     &x,
		}
		for i := 0; i < numberOfWorker; i++ {
			go activityLogSvcObj.processActivityLog(i)
		}
		//go activityLogSvcObj.gracefulShutter()
	}
	return activityLogSvcObj
}

type ActivityLogConfig struct {
	BufferSize                  int
	NumberOfWorker              int
	ShutdownGracePeriodInSecond int
}

func (s ActivityLogService) Shutdown(ctx context.Context) {
	s.getReadyForShutdown()
	go func() {
		<-ctx.Done()
		s.shutdownWorkers()
	}()
}

func (s ActivityLogService) shutdownWorkers() {
	for i := 0; i < s.config.NumberOfWorker; i++ {
		s.workerStopperChannel <- true
	}
	fmt.Println("activity log stopped")
}

func (s ActivityLogService) processActivityLog(workerId int) {
Loop:
	for {
		select {
		case <-s.workerStopperChannel:
			break Loop

		case activityLogMod := <-s.activityLogBuffer:
			err := s.handler.Create(activityLogMod)
			_ = err
		}
	}
}

var activityLogSvcObj *ActivityLogService

type ActivityLogService struct {
	activityLogBuffer    chan models.ActivityLog
	config               ActivityLogConfig
	handler              ActivityLogHandlerInterface
	workerStopperChannel chan bool //don't use any buffer here.
	shuttingDownFlag     *int32
}

func (s ActivityLogService) getReadyForShutdown() {
	atomic.StoreInt32(s.shuttingDownFlag, 1)
	fmt.Println("activity log not accepting new elements anymore")
}

func (s ActivityLogService) isShuttingDown() bool {
	val := atomic.LoadInt32(s.shuttingDownFlag)
	return val == 1
}

func (s ActivityLogService) CreateActivityLog(activityLogModel models.ActivityLog) error {
	if s.isShuttingDown() {
		return errors.New("cannot create activity")
	}

	s.activityLogBuffer <- activityLogModel
	return nil
}

type ActivityLogServiceInterface interface {
	CreateActivityLog(activityLogModel models.ActivityLog) error
	Shutdown(ctx context.Context)
}
