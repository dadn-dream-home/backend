package topic

import "github.com/mattn/go-sqlite3"

type Topic struct {
	Op    int
	Table string
}

func Insert(table string) Topic {
	return Topic{sqlite3.SQLITE_INSERT, table}
}

func Delete(table string) Topic {
	return Topic{sqlite3.SQLITE_DELETE, table}
}
