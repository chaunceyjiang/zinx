package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func (c *ConnManager) Add(connection ziface.IConnection) {
	//panic("implement me")
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[connection.GetConnID()] = connection
	fmt.Println("connection add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Remove(connection ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, connection.GetConnID())
	fmt.Println("connection Remove ConnID=", connection.GetConnID(), " successfully: conn num = ", connMgr.Len())
}
func (c *ConnManager) Len() int {
	return len(c.connections)
}
func (c *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnManager) ClearConn() {
	for _, conn := range c.connections {
		conn.Stop()
	}
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
}

