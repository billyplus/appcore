package appcore

import (
	"github.com/pkg/errors"
)

type IComponent interface {
	// Init接口在服务启动前调用，可以进行一些配置检查，或提前连接远程服务，如果Init发生错误，则服务启动不成功
	Init() error
	Start() error
	Shutdown()
	Name() string

	// 所有服务都需要能够接收其它服务发过去的消息
	// Receive(msg interface{})
}

type Mgr struct {
	logger     Logger
	errChan    chan error
	components []IComponent
}

func NewServiceMgr(o ...Option) *Mgr {
	m := &Mgr{
		errChan: make(chan error),
	}
	for _, option := range o {
		option(m)
	}
	return m
}

func (mgr *Mgr) Add(service ...IComponent) {
	mgr.components = append(mgr.components, service...)
}

func (mgr *Mgr) Serve() error {
	for _, svr := range mgr.components {
		if err := svr.Init(); err != nil {
			return errors.Wrapf(err, "[%s] failed to init", svr.Name())
		}
	}

	for _, svr := range mgr.components {
		go func(name string, fn func() error) {
			mgr.errChan <- func() (err error) {
				defer func() {
					var err error
					e := recover()
					if e != nil {
						if erro, ok := e.(error); ok {
							err = errors.WithStack(erro)
						} else {
							err = errors.Errorf("%v", e)
						}

						if mgr.logger != nil {
							mgr.logger.Errorf(err, "[%s] recover from panic", name)
						}
					}

					if mgr.logger != nil {
						if err != nil {
							mgr.logger.Errorf(err, "[%s] exit with error", name)
						} else {
							mgr.logger.Infof("[%s] exit", name)
						}
					}
				}()

				return fn()
			}()
		}(svr.Name(), svr.Start)
	}

	// Wait for the first actor to stop.
	err := <-mgr.errChan
	// if mgr.logger != nil {
	// 	if err != nil {
	// 		mgr.logger.Errorf(err, "service exit with error")
	// 	} else {
	// 		mgr.logger.Infof("service exit")
	// 	}
	// }

	// shutdown all service，实际关闭所有服务
	for _, svr := range mgr.components {
		svr.Shutdown()
	}

	for i := 1; i < len(mgr.components); i++ {
		<-mgr.errChan
	}
	return err
}
