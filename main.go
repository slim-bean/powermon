package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/common/model"
	"github.com/slim-bean/powermon/pkg/sigineer"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grafana/loki-client-go/loki"
	"go.bug.st/serial"

	"github.com/slim-bean/powermon/pkg/eg4"
)

var (
	labels = model.LabelSet{
		"job":  "powermon",
		"type": "batt",
	}
	labels_raw = model.LabelSet{
		"job":  "powermon",
		"type": "batt_raw",
	}
	labels_inverter = model.LabelSet{
		"job":  "powermon",
		"type": "inverter",
	}
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)
	logger = level.NewFilter(logger, level.AllowDebug())

	cfg := loki.Config{}
	// Sets defaults as well as anything from the command line
	cfg.RegisterFlags(flag.CommandLine)
	flag.Parse()

	c, err := loki.NewWithLogger(cfg, logger)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create client", "err", err)
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
	}
	port, err := serial.Open("/dev/ttyUSB1", mode)
	if err != nil {
		level.Error(logger).Log("msg", "failed to open serial port", "err", err)
		os.Exit(1)
	}
	err = port.SetReadTimeout(1 * time.Second)
	if err != nil {
		level.Error(logger).Log("msg", "failed to set port read timeout, continuing anyway", "err", err)
	}

	batt := 0
	// I haven't been able to figure out the checksum, so instead I'm just hard coding the commands here.
	commands := [][]byte{
		{0x7E, 0x01, 0x01, 0x00, 0xFE, 0x0D},
		{0x7E, 0x02, 0x01, 0x00, 0xFC, 0x0D},
	}

	go func() {
		buff := make([]byte, 500)
		for {
			time.Sleep(1000 * time.Millisecond)
			n, err := port.Write(commands[batt])
			if err != nil {
				level.Error(logger).Log("msg", "failed to send command on serial port", "err", err)
				continue
			}

			n, err = port.Read(buff)
			if err != nil {
				level.Error(logger).Log("msg", "failed to read from serial port", "err", err)
				continue
			}

			if n > 0 && len(buff) > 0 && buff[0] == 0x7E && buff[n-1] == 0x0D {
				packet, err := eg4.Parse(buff)
				if err != nil {
					level.Error(logger).Log("msg", "error parsing battery packet", "err", err)
					continue
				}
				ps, err := json.Marshal(packet)
				if err != nil {
					level.Error(logger).Log("msg", "failed to marshal battery packet to json", "err", err)
					continue
				}
				fmt.Println(string(ps))
				err = c.Handle(labels, time.Now(), string(ps))
				if err != nil {
					level.Error(logger).Log("msg", "failed to send logs to loki client", "err", err)
				}
				sb := strings.Builder{}
				for i := 0; i < n; i++ {
					sb.WriteString(fmt.Sprintf("%d:0x%02X ", i, buff[i]))
				}
				err = c.Handle(labels_raw, time.Now(), sb.String())
				if err != nil {
					level.Error(logger).Log("msg", "failed to send logs to loki client", "err", err)
				}

				// Move between batteries, only do this after a successful read or else you can get out of "time" with a battery
				batt = batt + 1
				// Currently only have 2 batteries, so if the address after increment is 3, move it back to 1
				if batt > len(commands)-1 {
					batt = 0
				}
			} else {
				level.Error(logger).Log("msg", "did not receive valid packet from battery, ignoring this poll.")
			}
		}
	}()

	inverterMode := &serial.Mode{
		BaudRate: 2400,
		DataBits: 8,
	}
	inverterPort, err := serial.Open("/dev/ttyUSB0", inverterMode)
	if err != nil {
		level.Error(logger).Log("msg", "failed to open serial port", "err", err)
		os.Exit(1)
	}
	err = inverterPort.SetReadTimeout(1 * time.Second)
	if err != nil {
		level.Error(logger).Log("msg", "failed to set port read timeout, continuing anyway", "err", err)
	}

	commandChan := make(chan string)
	responseChan := make(chan string)

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		buff := make([]byte, 500)
		for {
			select {
			case command := <-commandChan:
				n, err := inverterPort.Write([]byte(command))
				if err != nil {
					level.Error(logger).Log("msg", "failed to send command on serial port", "err", err)
					continue
				}
				n, err = inverterPort.Read(buff)
				if err != nil {
					level.Error(logger).Log("msg", "failed to read from serial port", "err", err)
					continue
				}
				resp := string(buff[:n])
				responseChan <- resp
			case <-t.C:
				// Q1<cr>
				// <cr>
				n, err := inverterPort.Write([]byte{0x0D})
				if err != nil {
					level.Error(logger).Log("msg", "failed to send command on serial port", "err", err)
					continue
				}

				n, err = inverterPort.Read(buff)
				if err != nil {
					level.Error(logger).Log("msg", "failed to read from serial port", "err", err)
					continue
				}
				resp := string(buff[:n])
				packet, err := sigineer.Parse(resp)
				if err != nil {
					level.Error(logger).Log("msg", "error parsing UPS packet", "err", err)
					continue
				}
				ps, err := json.Marshal(packet)
				if err != nil {
					level.Error(logger).Log("msg", "failed to marshal UPS packet to json", "err", err)
					continue
				}
				fmt.Println(resp)
				fmt.Println(string(ps))
				err = c.Handle(labels_inverter, time.Now(), string(ps))
				if err != nil {
					level.Error(logger).Log("msg", "failed to send logs to loki client", "err", err)
				}
			}
		}
	}()

	sh := commandHandler{
		commandChan:  commandChan,
		responseChan: responseChan,
		logger:       logger,
	}

	http.HandleFunc("/command", sh.command)
	go func() { http.ListenAndServe(":8080", nil) }()

	done := make(chan struct{})

	go func() {
		level.Info(logger).Log("msg", "running...")
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		level.Info(logger).Log("msg", "shutting down...")
		close(done)
	}()

	<-done
}

type commandHandler struct {
	commandChan  chan string
	responseChan chan string
	logger       log.Logger
}

func (s *commandHandler) command(w http.ResponseWriter, req *http.Request) {
	c := req.URL.Query().Get("command")
	if c == "" {
		fmt.Fprintf(w, "empty command query param\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c = c + "\r"
	level.Info(s.logger).Log("msg", "sending command to inverter", "command", c)

	s.commandChan <- c
	level.Info(s.logger).Log("msg", "waiting for response")
	r := <-s.responseChan
	level.Info(s.logger).Log("msg", "received response", "response", r)
	fmt.Fprintf(w, "response: %s", r)
	w.WriteHeader(http.StatusOK)
}
