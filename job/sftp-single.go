package main

import (
	"log"
	"os"
	"time"
	"path"
	"io"
	"fmt"
	"sync"
	"strconv"
	"path/filepath"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
)

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*16644525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

func main() {

	host := "202.38.220.15"
	port := "22"
	user := "unkcpz"
	pass := "sunshine"

	// get host public key
	// hostKey := getHostKey(host)

	conn, err := connect(user, pass, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	dir := "/scratch/unkcpz"
	prefix := "tmp"
	pathname, err := TempDir(client, dir, prefix)
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// send files
	err = SendFiles(client, wd, pathname)
	if err != nil {
		log.Fatal(err)
	}

	// create a new file mimic the job execution
	newFile, err := client.Create(path.Join(pathname, "newfile"))
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	_, err = newFile.Write([]byte("writes\n"))
	if err != nil {
		log.Fatal(err)
	}

	// recive files
	err = ReciveFiles(client, pathname, wd)
	if err != nil {
		log.Fatal(err)
	}
}

func connect(user, password, host, port string) (*ssh.Client, error) {
	// ssh client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		// verify host public key
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		// optional tcp connect timeout
		Timeout: 5 * time.Second,
	}

	// connect
	conn, err := ssh.Dial("tcp", host+":"+port, config)

	return conn, err
}

func TempDir(client *sftp.Client, dir, prefix string) (name string, err error) {
	if dir == "" {
		dir = os.TempDir()
	}
	nconflict := 0

	sshFxFailure := uint32(4)

	for i := 0; i < 10000; i++ {
		try := filepath.Join(dir, prefix+nextRandom())
		err = client.Mkdir(try)
		if status, ok := err.(*sftp.StatusError); ok {
			if status.Code == sshFxFailure {
				if nconflict++; nconflict > 10 {
					randmu.Lock()
					rand = reseed()
					randmu.Unlock()
				}
				continue
			}
			return "", err
		}
		name = try
		break
	}
	return name, nil
}

func SendFiles(client *sftp.Client, fromDir, toDir string) error {
	// loop over pwd files
	files, err := ioutil.ReadDir(fromDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		// create source file
		srcFile, err := os.Open(path.Join(fromDir, f.Name()))
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// create destination file
		dstFile, err := client.Create(path.Join(toDir, f.Name()))
		if err != nil {
			return err
		}
		defer dstFile.Close()
		// copy source file to destination file
		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %d bytes copied\n", f.Name(), bytes)
	}
	return nil
}

func ReciveFiles(client *sftp.Client, fromDir, toDir string) error {
	// retrieve all files from the remote directory
	files, err := client.ReadDir(fromDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		// create destination file
		var dstFile *os.File
		filename := path.Join(toDir, f.Name())
		dstFile, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// open source file
		srcFile, err := client.Open(path.Join(fromDir, f.Name()))
		if err != nil {
			return err
		}

		// copy source file to destination file
		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
		fmt.Printf("%d bytes copied\n", bytes)

		// flush in-memory copy
		err = dstFile.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}
