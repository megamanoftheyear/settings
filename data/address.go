package data

type AddressType uint8

const (
	InternalType AddressType = iota
	ExternalType
	DestinationType
	ServiceType // служебные адреса
)

type Status uint8

func (s Status) IsBlocked() bool {
	return s != ActiveStatus
}

const (
	ActiveStatus Status = iota
	BlockedStatus
	BlockedByAdminStatus
)

type Address struct {
	ID          uint64
	Address     string
	ClientID    string
	Status      Status
	Type        AddressType
	IsDefault   bool
	DelegatedTo string
}

type AddressDict interface {
	Get(address string) (*Address, bool)
	FindDefault(clientID string, addressType AddressType) (*Address, bool)
}
