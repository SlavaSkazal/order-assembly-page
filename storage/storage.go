package storage

type Storage interface {
	CreateTables() error
	CreateRecords() error
	PrintAssemblyPage(orderNumbers []int) error
}
