package storage

import (
	"bytes"
	"encoding/json"
	"github.com/jlaffaye/ftp"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"reflect"
	"strings"
	"time"
)

type FtpKind struct {
	Client        *ftp.ServerConn
	Configuration *FtpKindConfiguration
}

type FtpKindConfiguration struct {
	Host     string `json:"host" slixx:"HOST" default:"ftp.slixx.app:21"`
	Timeout  int64  `json:"timeout" slixx:"LONG" default:"1000"`
	File     string `json:"file" slixx:"PATH" default:"/"`
	Username string `json:"username" slixx:"STRING" default:"root"`
	Password string `json:"password" slixx:"PASSWORD" default:""`
}

func (kind *FtpKind) GetName() string {
	return "FTP"
}

func (kind *FtpKind) CanStore() bool {
	return true
}

func (kind *FtpKind) CanRead() bool {
	return true
}

func (kind *FtpKind) GetConfiguration() any {
	return kind.Configuration
}

func (kind *FtpKind) Initialize(rawConfiguration any) error {
	configuration := rawConfiguration.(*FtpKindConfiguration)

	client, err := ftp.Dial(
		configuration.Host,
		ftp.DialWithTimeout(time.Duration(configuration.Timeout)*time.Millisecond),
	)
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
	size, err := kind.Client.FileSize(fileutils.FixedPathName(name))
	if err != nil {
		return 0, err
	}
	return uint64(size), nil
}

func (kind *FtpKind) Store(name string, data []byte, offset uint64) error {
	err := kind.Client.StorFrom(fileutils.FixedPathName(kind.Configuration.File+name), bytes.NewBuffer(data), offset)

	if err != nil {
		return err
	}
	return nil
}

func (kind *FtpKind) FileInfo(name string) (fileutils.FileInfo, error) {
	info, err := kind.Client.GetEntry(fileutils.FixedPathName(kind.Configuration.File + "/" + name))
	if err != nil {
		return fileutils.FileInfo{}, err
	}
	return fileutils.FileInfo{
		Name:          name,
		FullDirectory: fileutils.FixedPathName(kind.Configuration.File + name),
		RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(kind.Configuration.File+name, kind.Configuration.File)),
		CreatedAt:     info.Time.Unix(),
		Directory:     info.Type == ftp.EntryTypeFolder,
		Size:          info.Size,
	}, nil
}

func (kind *FtpKind) CreateDirectory(name string) error {
	err := kind.Client.MakeDir(fileutils.FixedPathName(kind.Configuration.File + name))
	if err != nil {
		err = kind.CreateDirectory(fileutils.ParentDirectory(name))
		if err != nil {
			return err
		}
		return kind.Client.MakeDir(fileutils.FixedPathName(kind.Configuration.File + name))
	}
	return nil
}

func (kind *FtpKind) ListFiles(directory string) ([]fileutils.FileInfo, error) {
	baseDir := fileutils.FixedPathName(kind.Configuration.File + "/" + directory)
	entries, err := kind.Client.List(baseDir)
	if err != nil {
		return nil, err
	}
	var files []fileutils.FileInfo
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFolder {
			files = append(files, fileutils.FileInfo{
				Name:          entry.Name,
				FullDirectory: fileutils.FixedPathName(baseDir + "/" + entry.Name),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(baseDir+"/"+entry.Name, kind.Configuration.File), directory)),
				CreatedAt:     entry.Time.Unix(),
				Directory:     true,
				Size:          entry.Size,
			})
			err := kind.listFiles(baseDir+"/"+entry.Name, &files, directory)
			if err != nil {
				return nil, err
			}
		} else {
			files = append(files, fileutils.FileInfo{
				Name:          entry.Name,
				FullDirectory: fileutils.FixedPathName(baseDir + "/" + entry.Name),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(baseDir+"/"+entry.Name, kind.Configuration.File), directory)),
				CreatedAt:     entry.Time.Unix(),
				Directory:     false,
				Size:          entry.Size,
			})
		}
	}
	return files, nil
}

func (kind *FtpKind) listFiles(path string, files *[]fileutils.FileInfo, directory string) error {
	entries, err := kind.Client.List(fileutils.FixedPathName(path))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFolder {
			*files = append(*files, fileutils.FileInfo{
				Name:          entry.Name,
				FullDirectory: fileutils.FixedPathName(path + "/" + entry.Name),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(path+"/"+entry.Name, kind.Configuration.File), directory)),
				CreatedAt:     entry.Time.Unix(),
				Directory:     true,
				Size:          entry.Size,
			})
			err := kind.listFiles(path+"/"+entry.Name, files, directory)
			if err != nil {
				return err
			}
		} else {
			*files = append(*files, fileutils.FileInfo{
				Name:          entry.Name,
				FullDirectory: fileutils.FixedPathName(path + "/" + entry.Name),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(path+"/"+entry.Name, kind.Configuration.File), directory)),
				CreatedAt:     entry.Time.Unix(),
				Directory:     false,
				Size:          entry.Size,
			})
		}
	}
	return nil
}

func (kind *FtpKind) Read(file string, offset uint64, size uint64) ([]byte, error) {
	reader, err := kind.Client.RetrFrom(fileutils.FixedPathName(file), offset)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	bytes := make([]byte, size)

	_, err = reader.Read(bytes)
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

	// Add trailing slash to file to prevent issues
	if !strings.HasSuffix(configuration.File, "/") {
		configuration.File = configuration.File + "/"
	}

	return &configuration, nil
}

func (kind *FtpKind) Delete(file string) error {
	return kind.Client.Delete(fileutils.FixedPathName(file))
}

func (kind *FtpKind) DeleteDirectory(directory string) error {
	return kind.Client.RemoveDir(fileutils.FixedPathName(directory))
}

func (kind *FtpKind) DefaultConfiguration() interface{} {
	return reflect.New(reflect.TypeOf(FtpKindConfiguration{})).Interface()
}

func (kind *FtpKind) Close() error {
	err := kind.Client.Quit()
	if err != nil {
		return err
	}
	return nil
}
