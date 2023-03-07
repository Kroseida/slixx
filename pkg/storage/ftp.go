package storage

import (
	"bytes"
	"encoding/json"
	"github.com/jlaffaye/ftp"
	"path/filepath"
	"strings"
	"time"
)

type FtpKind struct {
	Client        *ftp.ServerConn
	Configuration *FtpKindConfiguration
}

type FtpKindConfiguration struct {
	Host     string `json:"host" slixx:"HOST"`
	Timeout  int64  `json:"timeout" slixx:"LONG"`
	File     string `json:"file" slixx:"PATH"`
	Username string `json:"username" slixx:"STRING"`
	Password string `json:"password" slixx:"PASSWORD"`
}

func (kind *FtpKind) GetName() string {
	return "FTP"
}

func (kind *FtpKind) Initialize(rawConfiguration any) error {
	configuration := rawConfiguration.(*FtpKindConfiguration)

	client, err := ftp.Dial(configuration.Host, ftp.DialWithTimeout(time.Duration(configuration.Timeout)*time.Millisecond))
	if err != nil {
		return err
	}

	err = client.Login(configuration.Username, configuration.Password)
	if err != nil {
		return err
	}

	kind.Client = client
	kind.Configuration = configuration

	return nil
}

func (kind *FtpKind) Size(name string) (uint64, error) {
	size, err := kind.Client.FileSize(name)
	if err != nil {
		return 0, err
	}
	return uint64(size), nil
}

func (kind *FtpKind) Store(name string, data []byte, offset uint64) error {
	err := kind.createParentDirectory(kind.Configuration.File + name)
	if err != nil {
		return err
	}

	err = kind.Client.StorFrom(kind.Configuration.File+name, bytes.NewBuffer(data), offset)

	if err != nil {
		return err
	}
	return nil
}

func (kind *FtpKind) ListFiles() ([]string, error) {
	entries, err := kind.Client.List(kind.Configuration.File)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFolder {
			kind.listFiles(kind.Configuration.File+"/"+entry.Name, &files)
		} else {
			files = append(files, kind.Configuration.File+"/"+entry.Name)
		}
	}
	return files, nil
}

func (kind *FtpKind) listFiles(path string, files *[]string) error {
	entries, err := kind.Client.List(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFolder {
			kind.listFiles(path+"/"+entry.Name, files)
		} else {
			*files = append(*files, path+"/"+entry.Name)
		}
	}
	return nil
}

func (kind *FtpKind) Read(file string, offset uint64, size uint64) ([]byte, error) {
	reader, err := kind.Client.RetrFrom(file, offset)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	bytes := make([]byte, size)

	_, err = reader.Read(bytes)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (kind *FtpKind) Parse(configurationJson string) (interface{}, error) {
	var configuration FtpKindConfiguration
	err := json.Unmarshal([]byte(configurationJson), &configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (kind *FtpKind) Delete(file string) error {
	return kind.Client.Delete(file)
}

func (kind *FtpKind) DefaultConfiguration() interface{} {
	return &FtpKindConfiguration{}
}

func (kind *FtpKind) Close() error {
	err := kind.Client.Quit()
	if err != nil {
		return err
	}
	return nil
}

func (kind *FtpKind) createParentDirectory(file string) error {
	parent := strings.ReplaceAll(filepath.Dir(file), "\\", "/")
	err := kind.Client.MakeDir(parent)
	if err != nil {
		if err.Error() == "550 Could not create directory. Raw error: 3" {
			kind.createParentDirectory(parent)
		} else {
			return err
		}
	}
	err = kind.Client.MakeDir(parent)
	if err != nil {
		return err
	}
	return nil
}
