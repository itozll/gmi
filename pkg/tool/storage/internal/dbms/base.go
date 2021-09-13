package dbms

type DBMS struct{}

func (DBMS) IsMongo() bool { return false }
func (DBMS) IsMysql() bool { return false }
