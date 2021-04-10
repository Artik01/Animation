package main

import(
	"fmt"
	"net"
	"encoding/json"
	"math/rand"
	"time"
)

type Cell struct {
	Old Coordinate	`json:"old"`
	New Coordinate	`json:"new"`
	Color int `json:color`
	vx, vy int
}

type Coordinate struct {
	X int		`json:"x"`
	Y int		`json:"y"`
}

var maxX, maxY int = 80, 25
func (c *Cell) Step() {
	if c.New.X+c.vx < 0 || c.New.X+c.vx >= maxX {
		c.vx = -c.vx
	}
	if c.New.Y+c.vy < 0 || c.New.Y+c.vy >= maxY {
		c.vy = -c.vy
	}
	c.Old.X=c.New.X
	c.Old.Y=c.New.Y
	c.New.X += c.vx
	c.New.Y += c.vy
}

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:10257")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	rand.Seed(int64(time.Now().Nanosecond()))
	var object Cell
	object.New.X = rand.Intn(maxX)
	object.New.Y = rand.Intn(maxY)
	for object.vx == 0 {
		object.vx = rand.Intn(3)-1
	}
	for object.vy == 0 {
		object.vy = rand.Intn(3)-1
	}
	object.Color = rand.Intn(16)+1 // 1..16
	
	for {
		data, err := json.Marshal(&object)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		object.Step()
		time.Sleep(250*time.Millisecond)
	}
	conn.Close()
}
