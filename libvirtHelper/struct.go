package libvirtHelper

import (
	"Instancer-worker-go/config"
	"fmt"
	"sync"

	"github.com/libvirt/libvirt-go"
)

type LibvirtManager struct {
	Lock sync.RWMutex
	Conn *libvirt.Connect
}

func (m *LibvirtManager) GetConnection() (conn *libvirt.Connect, err error) {
	var alive bool

	m.Lock.RLock()
	if m.Conn != nil { // check if conn is exists
		if alive, err = m.Conn.IsAlive(); err == nil && alive { // check if conn is alive
			// alive
			conn = m.Conn
			m.Lock.RUnlock()
			return conn, nil
		}
		// error or not alive -> reconnect
	}
	m.Lock.RUnlock()

	// failed to get a health connection -> reconnect
	m.Lock.Lock()
	defer m.Lock.Unlock()

	// if other goroutine repaired connection
	if m.Conn != nil {
		if alive, err = m.Conn.IsAlive(); err == nil && alive {
			// healthy
			return conn, nil
		}
		// still unhealthy
		m.Conn.Close()
	}

	var newConn *libvirt.Connect
	if newConn, err = libvirt.NewConnect(config.QemuUrl); err != nil {
		// failed to reconnect
		return nil, fmt.Errorf("Failed to reconnect libvirt : %v", err)
	}

	m.Conn = newConn

	// register callbacks

	return m.Conn, nil
}

func (m *LibvirtManager) Close() (err error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if m.Conn != nil {
		_, err = m.Conn.Close()
		m.Conn = nil
		return err
	}

	return nil
}

// 輔助函式：註冊連線關閉的回調
//func RegisterCloseCallback(conn *libvirt.Connect) {
//	// 定義回調函數
//	callback := func(c *libvirt.Connect, reason int, opaque interface{}) {
//		fmt.Printf("警告: Libvirt 連線已斷開 (Reason Code: %d)\n", reason)
//	}
//	// 註冊
//	if err := conn.RegisterCloseCallback(callback, nil); err != nil {
//		fmt.Printf("無法註冊 CloseCallback: %v\n", err)
//	}
//}
