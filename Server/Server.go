package main

import(
	"fmt"
	"net"
	"encoding/json"
	"github.com/nsf/termbox-go"
	"time"
)

type object struct {
	Old Coordinate	`json:"old"`
	New Coordinate	`json:"new"`
	Color int `json:color`
}

type Coordinate struct {
	X int		`json:"x"`
	Y int		`json:"y"`
}

type Data struct {
	x int
	y int
	color termbox.Attribute
}

var DB []Data

func main() {
	adr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10257")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	termbox.Init()
	termbox.Clear(termbox.ColorWhite,termbox.ColorBlack)
	
	go Draw()
	ln, err := net.ListenUDP("udp", adr)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	handleConnection(ln)
	termbox.Close()
}

func handleConnection(src net.Conn) {
	for {
		buf := make([]byte, 2000)
		n, err := src.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		var obj object
		err = json.Unmarshal(buf[0:n], &obj)
		if err != nil {
			fmt.Println(err)
			return
		}
		var found bool
		for i, v := range DB {
			if (v.x == obj.Old.X) && (v.y == obj.Old.Y) {
				DB[i].x = obj.New.X
				DB[i].y = obj.New.Y
				found = true
				break
			}
		}
		if !found {
			DB=append(DB, Data{obj.New.X, obj.New.Y, termbox.Attribute(obj.Color)})
		}
	}
}

func Draw() {
	for {
		termbox.Clear(termbox.ColorWhite,termbox.ColorBlack)
		for _, val := range DB {
			termbox.SetCell(val.x, val.y, ' ', termbox.ColorDefault, val.color)
		}
		termbox.Flush()
		time.Sleep(1*time.Millisecond)
	}
}
