package storage

type Kind interface {
	GetName() string
	Initialize(configuration any) error
	Store(name string, data []byte, offset uint64) error
	ListFiles() ([]string, error)
	Size(file string) (uint64, error)
	Read(file string, offset uint64, size uint64) ([]byte, error)
	Delete(file string) error
	Parse(configurationJson string) (interface{}, error)
	DefaultConfiguration() interface{}
	Close() error
}

var kinds = map[string]Kind{
	"FTP": &FtpKind{},
}

func ValueOf(name string) Kind {
	return kinds[name]
}

func Values() []Kind {
	values := make([]Kind, 0, len(kinds))
	for _, value := range kinds {
		values = append(values, value)
	}
	return values
}
