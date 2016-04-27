package main

type dbManager struct {
    dbs []*db
}

func newDBManager() *dbManager {
    return & dbManager {
        dbs: []*db{newDB()},
    }
}

func (dm *dbManager) getDefaultDB() *db {
    return dm.getDB(0);
}

func (dm *dbManager) getDB(i int) *db {
    return dm.dbs[i]
}

