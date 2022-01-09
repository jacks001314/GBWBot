package ssh

import (
	"github.com/pkg/sftp"
	"io"
	"os"
	"path/filepath"
)

type SftpClient struct {

	sshClient *SSHClient
	sftClient  *sftp.Client
}


func NewSftpClient(sshClient *SSHClient) (*SftpClient,error) {

	sftpClient,err := sftp.NewClient(sshClient.con)

	if err!=nil {

		return nil,err
	}

	return &SftpClient{
		sshClient: sshClient,
		sftClient: sftpClient,
	},nil
}

func (c *SftpClient) Close() {

	c.sftClient.Close()
}

func (c *SftpClient)UPloadFile(fpath string,remoteDir string) error {


	fname := filepath.Base(fpath)
	remoteFile,err := c.sftClient.Create(sftp.Join(remoteDir,fname))

	if err !=nil {

		return err
	}
	defer remoteFile.Close()

	localFile,err := os.Open(fpath)
	if err!=nil {

		return err
	}

	defer localFile.Close()

	//本地文件流拷贝到上传文件流
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *SftpClient) DownloadFile(rfile string,localDir string) error{

	fname := filepath.Base(rfile)
	lfile := filepath.Join(localDir,fname)

	srcFile,err := c.sftClient.Open(rfile)

	if err!=nil {
		return err
	}
	defer srcFile.Close()

	dstFile,err := os.Create(lfile)

	if err!=nil {

		return err
	}

	defer dstFile.Close()

	_,err = srcFile.WriteTo(dstFile)

	return err
}
