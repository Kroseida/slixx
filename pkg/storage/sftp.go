package storage

import (
	"encoding/json"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"os"
	"reflect"
	"strings"
	"time"
)

type SFtpKind struct {
	Client        *sftp.Client
	Configuration *SFtpKindConfiguration
}

type SFtpKindConfiguration struct {
	Host       string `json:"host" slixx:"HOST" default:"sftp.slixx.app:22"`
	Username   string `json:"username" slixx:"STRING" default:"root"`
	Password   string `json:"password" slixx:"PASSWORD" default:""`
	File       string `json:"file" slixx:"PATH" default:"/"`
	PrivateKey string `json:"privateKey" slixx:"STRING"` // Optional: Use for key-based authentication
	Timeout    int64  `json:"timeout" slixx:"LONG" default:"1000"`
}

func (kind *SFtpKind) GetName() string {
	return "SFTP"
}

func (kind *SFtpKind) CanStore() bool {
	return true
}

func (kind *SFtpKind) CanRead() bool {
	return true
}

func (kind *SFtpKind) GetConfiguration() any {
	return kind.Configuration
}

func (kind *SFtpKind) Initialize(rawConfiguration any) error {
	configuration := rawConfiguration.(*SFtpKindConfiguration)

	var auth []ssh.AuthMethod
	if configuration.PrivateKey != "" {
		key, err := ssh.ParsePrivateKey([]byte(configuration.PrivateKey))
		if err != nil {
			return err
		}
		auth = append(auth, ssh.PublicKeys(key))
	} else if configuration.Password != "" {
		auth = append(auth, ssh.Password(configuration.Password))
	}

	config := &ssh.ClientConfig{
		User:            configuration.Username,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(configuration.Timeout) * time.Millisecond,
	}

	conn, err := ssh.Dial("tcp", configuration.Host, config)
	if err != nil {
		return err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}

	kind.Client = client
	kind.Configuration = configuration

	return nil
}

func (kind *SFtpKind) Size(name string) (uint64, error) {
	stat, err := kind.Client.Stat(fileutils.FixedPathName(kind.Configuration.File + "/" + name))
	if err != nil {
		return 0, err
	}
	return uint64(stat.Size()), nil
}

func (kind *SFtpKind) Store(name string, data []byte, offset uint64) error {
	// Open the file for writing, create it if it does not exist
	file, err := kind.Client.OpenFile(fileutils.FixedPathName(kind.Configuration.File+name), os.O_WRONLY|os.O_CREATE)
	if err != nil {
		return err
	}
	defer file.Close()
	// If an offset is specified, seek to it
	if offset > 0 {
		_, err = file.Seek(int64(offset), io.SeekStart)
		if err != nil {
			return err
		}
	}
	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (kind *SFtpKind) FileInfo(name string) (fileutils.FileInfo, error) {
	stat, err := kind.Client.Stat(fileutils.FixedPathName(kind.Configuration.File + name))
	if err != nil {
		return fileutils.FileInfo{}, err
	}
	return fileutils.FileInfo{
		FullDirectory: fileutils.FixedPathName(kind.Configuration.File + name),
		RelativePath:  fileutils.FixedPathName(name),
		CreatedAt:     stat.ModTime().Unix(),
		Directory:     stat.IsDir(),
		Size:          uint64(stat.Size()),
	}, nil
}

func (kind *SFtpKind) FileInfoWithoutConfiguration(name string) (fileutils.FileInfo, error) {
	stat, err := kind.Client.Stat(fileutils.FixedPathName(name))
	if err != nil {
		return fileutils.FileInfo{}, err
	}
	return fileutils.FileInfo{
		FullDirectory: fileutils.FixedPathName(name),
		RelativePath:  fileutils.FixedPathName(name),
		CreatedAt:     stat.ModTime().Unix(),
		Directory:     stat.IsDir(),
		Size:          uint64(stat.Size()),
	}, nil
}

func (kind *SFtpKind) CreateDirectory(name string) error {
	_, err := kind.FileInfoWithoutConfiguration(fileutils.FixedPathName(kind.Configuration.File + "/" + name))
	if err != nil && err.Error() == "file does not exist" {
		err = kind.Client.Mkdir(fileutils.FixedPathName(kind.Configuration.File + "/" + name))
		if err != nil {
			err = kind.CreateDirectory(fileutils.ParentDirectory(name))
			if err != nil {
				return err
			}
			return kind.Client.Mkdir(kind.Configuration.File + name)
		}
	}
	return nil
}

func (kind *SFtpKind) ListFiles(directory string) ([]fileutils.FileInfo, error) {
	baseDir := fileutils.FixedPathName(kind.Configuration.File + "/" + directory)
	entries, err := kind.Client.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}
	var files []fileutils.FileInfo
	for _, entry := range entries {
		if entry.Name() == "." || entry.Name() == ".." || entry.Name() == "?" || entry.Name() == "*" {
			continue
		}
		if entry.IsDir() {
			files = append(files, fileutils.FileInfo{
				FullDirectory: fileutils.FixedPathName(baseDir + "/" + entry.Name()),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(baseDir+"/"+entry.Name(), kind.Configuration.File), directory)),
				CreatedAt:     entry.ModTime().Unix(),
				Directory:     true,
				Size:          uint64(entry.Size()),
			})
			err := kind.listFiles(baseDir+"/"+entry.Name(), &files, directory)
			if err != nil {
				return nil, err
			}
		} else {
			files = append(files, fileutils.FileInfo{
				FullDirectory: fileutils.FixedPathName(baseDir + "/" + entry.Name()),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(baseDir+"/"+entry.Name(), kind.Configuration.File), directory)),
				CreatedAt:     entry.ModTime().Unix(),
				Directory:     false,
				Size:          uint64(entry.Size()),
			})
		}
	}
	return files, nil
}

func (kind *SFtpKind) listFiles(path string, files *[]fileutils.FileInfo, directory string) error {
	entries, err := kind.Client.ReadDir(fileutils.FixedPathName(path))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == "." || entry.Name() == ".." || entry.Name() == "?" || entry.Name() == "*" {
			continue
		}
		if entry.IsDir() {
			*files = append(*files, fileutils.FileInfo{
				FullDirectory: fileutils.FixedPathName(path + "/" + entry.Name()),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(path+"/"+entry.Name(), kind.Configuration.File), directory)),
				CreatedAt:     entry.ModTime().Unix(),
				Directory:     true,
				Size:          uint64(entry.Size()),
			})
			err := kind.listFiles(path+"/"+entry.Name(), files, directory)
			if err != nil {
				return err
			}
		} else {
			*files = append(*files, fileutils.FileInfo{
				FullDirectory: fileutils.FixedPathName(path + "/" + entry.Name()),
				RelativePath:  fileutils.FixedPathName(strings.TrimPrefix(strings.TrimPrefix(path+"/"+entry.Name(), kind.Configuration.File), directory)),
				CreatedAt:     entry.ModTime().Unix(),
				Directory:     false,
				Size:          uint64(entry.Size()),
			})
		}
	}
	return nil
}

func (kind *SFtpKind) Read(file string, offset uint64, size uint64) ([]byte, error) {
	if size == 0 && offset == 0 {
		// If no size or offset is specified, read the entire file
		f, err := kind.Client.Open(fileutils.FixedPathName(kind.Configuration.File + "/" + file))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// Read the entire file
		buf, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}

	fullPath := fileutils.FixedPathName(kind.Configuration.File + "/" + file)

	// Open the file for reading
	f, err := kind.Client.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Seek to the specified offset
	_, err = f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read the specified amount of data
	buf := make([]byte, size)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf[:n], nil // Return the read data, up to n bytes
}

func (kind *SFtpKind) Parse(configurationJson string) (interface{}, error) {
	var configuration SFtpKindConfiguration
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

func (kind *SFtpKind) Delete(file string) error {
	return kind.Client.Remove(fileutils.FixedPathName(kind.Configuration.File + "/" + file))
}

func (kind *SFtpKind) DeleteDirectory(directory string) error {
	return kind.Client.RemoveDirectory(fileutils.FixedPathName(kind.Configuration.File + "/" + directory))
}

func (kind *SFtpKind) DefaultConfiguration() interface{} {
	return reflect.New(reflect.TypeOf(SFtpKindConfiguration{})).Interface()
}

func (kind *SFtpKind) Close() error {
	err := kind.Client.Close()
	if err != nil {
		return err
	}
	return nil
}
