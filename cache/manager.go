package cache

type Manager struct {
	defaultName string
	drivers     map[string]SimpleCacher
}

func NewManager() *Manager {
	return &Manager{
		drivers: make(map[string]SimpleCacher, 8),
	}
}

func (m *Manager) SetDefaultName(driverName string) {
	exist := m.CheckedDriverExist(driverName)
	if !exist {
		panic("cache driver: " + driverName + " is not registered")
	}
	m.defaultName = driverName
}

func (m *Manager) CheckedDriverExist(driverName string) bool {
	if _, ok := m.drivers[driverName]; !ok {
		return false
	} else {
		return true
	}
	return false
}
func (m *Manager) Register(driverName string, driver SimpleCacher) *Manager {
	m.defaultName = driverName
	m.drivers[driverName] = driver
	return m
}

func (m *Manager) UnRegister(driverName string) bool {
	exist := m.CheckedDriverExist(driverName)
	if !exist {
		return true
	}
	delete(m.drivers, driverName)
	return true
}
func (m *Manager) DefaultDriver() SimpleCacher {

	return m.GetDriver(m.defaultName)
}
func (m *Manager) Use(driverName string) SimpleCacher {
	m.SetDefaultName(driverName)
	return m.DefaultDriver()
}

func (m *Manager) GetDriver(driverName string) SimpleCacher {
	if driver, ok := m.drivers[driverName]; !ok {
		panic("cache driver: DefaultDriver is not registered")
	} else {
		return driver
	}
}

func (m *Manager) CloseAll() (err error) {
	for _, cache := range m.drivers {
		err = cache.Close()
	}
	return err
}
