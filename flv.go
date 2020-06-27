package ffmpegTest

import (
	"sync"
)

type Pusher struct {
	ID                uint
	Path              string
	Source            string
	gopCacheEnable    bool
	gopCacheLock      sync.RWMutex
	spsppsInSTAPaPack bool
	cond              *sync.Cond
	Stoped            bool
}

type Server struct {
	Stoped         bool
	pushers        map[string]*Pusher // Path <-> Pusher
	pushersLock    sync.RWMutex
	addPusherCh    chan *Pusher
	removePusherCh chan *Pusher
}

var Instance *Server = &Server{
	Stoped:         true,
	pushers:        make(map[string]*Pusher),
	addPusherCh:    make(chan *Pusher),
	removePusherCh: make(chan *Pusher),
}

func GetServer() *Server {
	return Instance
}

func (server *Server) Start() (err error) {
	go func() { // 保持到本地
		//inFilename := "rtsp://183.59.168.27/PLTV/88888905/224/3221227272/10000100000000060000000001030757_0.smil?icip=88888888"
		inFilename := "rtmp://58.200.131.2:1935/livetv/hunantv"
		ToHls(inFilename, "", "tcp")
	}()

	server.Stoped = false

	return
}

// Stop 停止
func (server *Server) Stop() {
	server.Stoped = true
	server.pushersLock.Lock()
	server.pushers = make(map[string]*Pusher)
	server.pushersLock.Unlock()

	close(server.addPusherCh)
	close(server.removePusherCh)
}

// AddPusher 添加推流进程
func (server *Server) AddPusher(pusher *Pusher) bool {

	added := false
	server.pushersLock.Lock()
	_, ok := server.pushers[pusher.Path]
	if !ok {
		server.pushers[pusher.Path] = pusher
		added = true
	} else {
		added = false
	}
	server.pushersLock.Unlock()
	if added {
		server.addPusherCh <- pusher
	}

	return added
}

// RemovePusher 移除推流
func (server *Server) RemovePusher(pusher *Pusher) {
	removed := false
	server.pushersLock.Lock()
	if _pusher, ok := server.pushers[pusher.Path]; ok && pusher.ID == _pusher.ID {
		delete(server.pushers, pusher.Path)
		removed = true
	}
	server.pushersLock.Unlock()
	if removed {
		server.removePusherCh <- pusher
	}
}

func (server *Server) GetPusher(path string) (pusher *Pusher) {
	server.pushersLock.RLock()
	pusher = server.pushers[path]
	server.pushersLock.RUnlock()
	return
}

func (server *Server) GetPushers() (pushers map[string]*Pusher) {
	pushers = make(map[string]*Pusher)
	server.pushersLock.RLock()
	for k, v := range server.pushers {
		pushers[k] = v
	}
	server.pushersLock.RUnlock()
	return
}

func (server *Server) GetPusherSize() (size int) {
	server.pushersLock.RLock()
	size = len(server.pushers)
	server.pushersLock.RUnlock()
	return
}
