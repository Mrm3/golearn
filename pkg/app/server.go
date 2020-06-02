package app

import (
	"context"
	"immortality-demo/config"
	"immortality-demo/routers"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ImmortalityAgent struct {
	context    context.Context
	cancelFunc context.CancelFunc
	waitGroup  sync.WaitGroup
}

func NewImmortalityAgent() *ImmortalityAgent {
	ctxt, cancel := context.WithCancel(context.Background())
	return &ImmortalityAgent{
		context:    ctxt,
		cancelFunc: cancel,
	}
}

func (agent *ImmortalityAgent) Start() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	errinfo := make(chan error, 1)

	ulog.Info("Starting immortality agent service...")
	apiCtxt, _ := context.WithCancel(agent.context)
	go startAPIServer(apiCtxt, agent, errinfo)
	select {
	case s := <-signals:
		ulog.Infof("Received system singal %s to abort service...", s)
		agent.Stop()
	case err := <-errinfo:
		ulog.Infof("starting service failed:%v", err)
		agent.Stop()
	}
}

func startAPIServer(ctxt context.Context, agent *ImmortalityAgent, errinfo chan error) {
	server := &http.Server{Addr: config.Config.Listen, Handler: routers.Router()}

	go func() {
		agent.waitGroup.Add(1)
		defer agent.waitGroup.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errinfo <- err
		}
	}()

	select {
	case <-ctxt.Done():
		ulog.Info("immortality agent service exiting:", ctxt.Err())
		if err := server.Shutdown(ctxt); err != nil {
			ulog.Error("immortality agent shutdown failed:", err)
		}
	}
}

func (agent *ImmortalityAgent) Stop() {
	agent.cancelFunc()
	agent.waitGroup.Wait()
}
