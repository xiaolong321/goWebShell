package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"os/exec"
	"bufio"
	"io"
	"strings"
)


func main() {
	handle := martini.Classic()
	handle.Use(martini.Static("static"))
	handle.Use(render.Renderer(render.Options{
		Directory: "view",
		Extensions: []string{".html", },
	}))

	handle.Get("/", webShell)
	handle.Post("/", webShell)
	handle.Get("/ws", ws)
	handle.RunOnAddr(":8080")
}


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GoCommand struct {
	Name string
	Args []string
}

var cmd string

func webShell(ren render.Render, r *http.Request, w http.ResponseWriter) {
	host := r.Host
	if r.Method == "POST" {
		r.ParseForm()
		cmd = r.Form["command"][0]
	}

	ren.HTML(200, "webshell", host)

}



func ws (r *http.Request, w http.ResponseWriter){
	ws, err := upgrader.Upgrade(w, r, nil)

	go func() {
		if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(1 * time.Second)); err != nil {
			return
		}
	}()

	if err == nil {
		data := make(chan string, 1)
		cmdList := strings.Split(cmd, " ")
		cmd = ""
		cmd := GoCommand{
			Name:cmdList[0],
			Args:cmdList[1:],
		}
		go cmd.Run(data)
		var str string
		hasMore := true
		for hasMore {
			if str, hasMore = <-data; hasMore {
				ws.WriteMessage(websocket.TextMessage, []byte(str))
			} else {
				ws.Close()
				return
			}
		}
	}
}


func (c *GoCommand) Run(data chan string) (err error) {
	cmd := exec.Command(c.Name, c.Args...)
	out, err := cmd.StdoutPipe()
	cmd.Start()
	if err != nil {
		data <- err.Error()
	}
	buf := bufio.NewReader(out)
	for {
		if line, err := buf.ReadString('\n'); err != nil || io.EOF == err {
			close(data)
			return err

		} else {
			data <- line
		}
	}
	cmd.Wait()
	return nil
}