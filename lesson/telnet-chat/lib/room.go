package lib

import (
	"bytes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type Room struct {
	listener net.Listener
	peoples  map[string]*People
	timeout  time.Duration
	closed   bool
	mu       sync.Mutex
	cap      int
}

func NewRoom(cap int, timeout time.Duration, listener net.Listener) (*Room, error) {
	if cap < 1 {
		return nil, errors.New("cap can not less than 1")
	}
	if timeout < 1 {
		return nil, errors.New("timeout can not less than 1")
	}

	room := &Room{
		listener: listener,
		peoples:  make(map[string]*People),
		timeout:  timeout,
		closed:   false,
		mu:       sync.Mutex{},
		cap:      cap,
	}
	return room, nil
}

func (r *Room) UpdatePeople(people *People) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for name, _ := range r.peoples {
		if name == people.name {
			r.peoples[name] = people
			return nil
		}
	}
	return errors.New("can not found the people")
}

func (r *Room) AddPeople(people *People) error {
	err := r.UpdatePeople(people)
	if err == nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.peoples)+1 > r.cap {
		people.Close()
		return errors.New("can not add into chat room.")
	}
	r.peoples[people.name] = people
	return nil
}

func (r *Room) RemovePeople(people *People) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for name, p := range r.peoples {
		if name != people.name {
			continue
		}
		p.Close()
		delete(r.peoples, name)
	}
}

func (r *Room) Dispatch(message *Message) {
	for name, people := range r.peoples {
		if name == message.src.name {
			continue
		}
		select {
		case people.message <- message:
			break
		default:
			log.Errorf("send msg to %s error", name)
			break
		}
	}
}

func (r *Room) Receive(people *People) {
	defer r.RemovePeople(people)
	for {
		receive := make([]byte, 2048)
		_, err := people.conn.Read(receive)
		if err != nil {
			log.Errorf("server receive msg error :", err)
			return
		}

		receive = bytes.TrimRight(receive, "\x00")
		receive = bytes.TrimRight(receive, "\x0a")
		receive = bytes.TrimRight(receive, "\x0d")
		if len(receive) == 0 {
			continue
		}

		message := &Message{
			src:     people,
			target:  nil,
			content: string(receive),
		}

		r.Dispatch(message)
	}
}

func (r *Room) Send(people *People) {
	for {
		select {
		case message, ok := <-people.message:
			if !ok {
				log.Errorf("people %s get message failed", people.name)
				continue
			}
			content := message.src.name + "说:" + message.content
			_, err := people.conn.Write([]byte(content))
			if err != nil {
				log.Errorf("people %s send message failed(%s)", people.name, message.content)
				continue
			}
			break
		}
	}
}

func (r *Room) Start() {
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			log.Errorf("server accept connect error :%v", err)
			continue
		}

		name := conn.RemoteAddr().String()
		people := &People{
			conn:    conn,
			name:    name,
			created: time.Now(),
			message: make(chan *Message, 200),
		}

		err = r.AddPeople(people)
		if err != nil {
			log.Error(err)
			continue
		}
		go r.Receive(people)
		go r.Send(people)
	}
}

func (r *Room) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true
	for _, people := range r.peoples {
		people.Close()
	}
}
