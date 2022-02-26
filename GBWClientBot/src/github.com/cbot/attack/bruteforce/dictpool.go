package bruteforce

import "sync"

type DictPool struct {
	lock  sync.Mutex
	users map[string]bool

	passwds map[string]bool

	entries []*DictEntry

	isUpdate bool
}

type DictEntry struct {
	user string
	pass string
}

func NewDictPool() *DictPool {

	dp := &DictPool{
		lock:     sync.Mutex{},
		users:    make(map[string]bool),
		passwds:  make(map[string]bool),
		isUpdate: false,
		entries:  nil,
	}

	dp.Add([]string{"root"}, []string{"root", "admin", "123456", "admin123456", "password", "passwd", "passwd123456"})

	return dp
}

func (d *DictPool) Add(users []string, passwds []string) {

	d.lock.Lock()
	defer d.lock.Unlock()

	for _, user := range users {

		d.users[user] = true
	}

	for _, pass := range passwds {

		d.passwds[pass] = true
	}

	d.isUpdate = true

}

func (d *DictPool) Dicts() []*DictEntry {

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.isUpdate || d.entries == nil {

		entries := make([]*DictEntry, 0)

		for user, _ := range d.users {

			for pass, _ := range d.passwds {

				entries = append(entries, &DictEntry{
					user: user,
					pass: pass,
				})
			}
		}

		d.entries = entries
	}

	d.isUpdate = false

	return d.entries
}

func (d *DictPool) ISUPdate() bool {

	return d.isUpdate
}
