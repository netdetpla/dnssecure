package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	randString := ""
	_ = SendUDP("", randString, "start")
	err := os.Mkdir(AppstatusPath, 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err.Error())
		os.Exit(10)
	}
	err = os.Mkdir(LogPath, 0777)
	if err != nil && !os.IsExist(err) {
		WriteError2Appstatus(err.Error(), 9)
	}
	//网络检查
	netCheckFlag, err := NetCheck()
	if err != nil || !netCheckFlag {
		ConnectFail()
		WriteError2Appstatus("Can not connect to the Internet.", 2)
	}
	//任务开始
	TaskStart()
	//读取配置
	GetConf()
	tasks, err := GetTaskConfig()
	if err != nil {
		GetConfFail()
		WriteError2Appstatus(err.Error(), 3)
	}
	GetConfSuccess()
	startTime := time.Now().Unix()
	_ = SendUDP(tasks.taskID, tasks.subID, "run: "+strconv.FormatInt(startTime, 10))
	//任务执行
	TaskRun()
	err = ControlDNSQueryRoutine(tasks)
	if err != nil {
		TaskRunFail()
		WriteError2Appstatus(err.Error(), 1)
	}
	ControlCompareRoutine(tasks)
	TaskRunSuccess()
	//进度
	process := len(tasks.records) / 30
	final_process := len(tasks.records) % 30
	for i := 0; i < process; i++ {
		err = SendProcess(tasks.taskID, tasks.uuid, "DomainInfo", 30, false)
		if err != nil {
			WriteResultFail()
			WriteError2Appstatus(err.Error(), 1)
		}
		//time.Sleep(time.Duration(1 * time.Second))
	}
	err = SendProcess(tasks.taskID, tasks.uuid, "DomainInfo", final_process, true)
	if err != nil {
		WriteResultFail()
		WriteError2Appstatus(err.Error(), 1)
	}
	//写结果
	WriteResult()
	err = ControlWriteResultRoutine(tasks)
	if err != nil {
		WriteResultFail()
		WriteError2Appstatus(err.Error(), 1)
	}
	WriteResultSuccess()
	endTime := time.Now().Unix()
	duration := endTime - startTime
	_ = SendUDP(tasks.taskID, tasks.subID, "len: "+strconv.Itoa(len(tasks.records))+"; duration: "+strconv.FormatInt(duration, 10))
	//写状态文件
	var (
		//err        error
		sftpClient *sftp.Client
	)

	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = connect("root", "root111111", "10.96.129.6", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	// 用来测试的本地文件路径 和 远程机器上的文件夹
	var localFilePath = "/tmp/result/" + tasks.taskID + ".result"
	var remoteDir = "/home/result/"
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	now := strconv.FormatInt(time.Now().Unix(), 10)
	var remoteFileName = randString + "-" + now //path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	for {
		buf := make([]byte, 1024)
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	fmt.Println("copy file to remote server finished!")
	WriteSuccess2Appstatus()
}

func connect(user, password, host string, port int) (*sftp.Client, error) {
 var (
   auth         []ssh.AuthMethod
   addr         string
   clientConfig *ssh.ClientConfig
   sshClient    *ssh.Client
   sftpClient   *sftp.Client
   err          error
 )
 // get auth method
 auth = make([]ssh.AuthMethod, 0)
 auth = append(auth, ssh.Password(password))

 clientConfig = &ssh.ClientConfig{
   User:    user,
   Auth:    auth,
   Timeout: 30 * time.Second,
	HostKeyCallback:func(hostname string,remote net.Addr,key ssh.PublicKey) error {
		return nil
	},
 }

 // connet to ssh
 addr = fmt.Sprintf("%s:%d", host, port)

 if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
   return nil, err
 }

 // create sftp client
 if sftpClient, err = sftp.NewClient(sshClient); err != nil {
   return nil, err
 }

 return sftpClient, nil
}
