package main

type db struct {
	kv map[string]value
}

func newDB() *db {
	return &db {
		kv: make(map[string]value),
	}
}

func (d *db) get(key string) value {
	v, ok := d.kv[key]	
	if ok {
		return v
	}
	return nil
}
