package session

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

type Provider interface {
	Create(sid string) (Session, error)
	Read(sid string) (Session, error)
	Destroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type Manager struct {
	cookieName string     //private cookie name
	lock       sync.Mutex // protects session
	provider   Provider
	maxAge     int64
}

var provides = make(map[string]Provider)

func NewManager(provideName, cookieName string, maxAge int64) *Manager {
	provider, ok := provides[provideName]
	if !ok {
		panic(fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName))
	}
	return &Manager{provider: provider, cookieName: cookieName, maxAge: maxAge}
}

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

func (manager *Manager) GetSession(r *http.Request) (Session, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err == nil {
		if cookie.Value != "" {
			sid, _ := url.QueryUnescape(cookie.Value)
			session, error := manager.provider.Read(sid)
			return session, error
		} else {
			err = http.ErrNoCookie
		}
	}
	return nil, err
}

func (manager *Manager) CreateSession(w http.ResponseWriter) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	sid := uuid.NewV4().String()
	session, _ = manager.provider.Create(sid)
	cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxAge)}
	http.SetCookie(w, &cookie)
	return
}

func (manager *Manager) DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.Destroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxAge)
	time.AfterFunc(time.Duration(manager.maxAge)*time.Second, func() { manager.GC() })
}
