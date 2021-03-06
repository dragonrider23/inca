package grabber

import (
	"github.com/lfkeitel/inca/src/common"
)

type connGroup struct {
	numOfConnections int
	goChan           chan bool
	conf             *common.Config
}

func newConnGroup(conf *common.Config) connGroup {
	return connGroup{
		conf: conf,
	}
}

func (c *connGroup) add(delta int) {
	if c.goChan == nil {
		c.goChan = make(chan bool)
	}
	c.numOfConnections += delta
}

func (c *connGroup) done() {
	c.add(-1)
	finishedDevices++
	if c.numOfConnections < c.conf.MaxSimultaneousConn {
		c.goChan <- true
	}
}

func (c *connGroup) wait() {
	if c.numOfConnections < c.conf.MaxSimultaneousConn {
		return
	}
	<-c.goChan
}
