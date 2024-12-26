package tableConfig

import (
	"Common/Table"
	"log/slog"
	"sync/atomic"
)

var value atomic.Value

func Instance() *Table.Tables {
	if val := value.Load(); val != nil {
		return val.(*Table.Tables)
	}
	return nil
}

func reload() (err error) {
	var (
		tables Table.Tables
	)

	err = tables.Init()
	if err != nil {
		return err
	}

	value.Store(&tables)
	slog.Info("tableConfig init ok")

	return
}

func Init() error {
	return reload()
}

func Reload() error {
	return reload()
}
