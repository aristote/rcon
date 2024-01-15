/*
Copyright 2024 JME

Redistribution and use in source and binary forms, with or without modification,
 are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this l
ist of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, thi
s list of conditions and the following disclaimer in the documentation and/or ot
her materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS” AND 
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WA
RRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
 IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT
, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
 BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, D
ATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF L
IABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE O
R OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 OF THE POSSIBILITY OF SUCH DAMAGE.-
*/

package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/gorcon/websocket"
	"github.com/guptarohit/asciigraph"
	"golang.org/x/term"
)

func main() {
	var fps int = 0
	var data []float64

	//Getting terminal height and width for configuring graph options
	width, height, err := term.GetSize(0)
	if err != nil {
		log.Fatal(err)
	}
	dataBuffer := width - 10 //Getting a bit of space free for text on the left

	conn, err := websocket.Dial("127.0.0.1:5305", "yourPassword")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		response, err := conn.Execute("fps")
		if err != nil {
			log.Fatal(err)
		}

		//Extracting the fps number from the response
		re := regexp.MustCompile("[0-9]+")
		fps, err = strconv.Atoi(re.FindAllString(response, -1)[0])
		if err != nil {
			log.Fatal(err)
		}

		//Adding the current fps to the data array and then triming the array if len > terminal width-10
		//so we have a nice graph display
		data = append(data, float64(fps))
		if dataBuffer > 0 && len(data) > dataBuffer {
			data = data[len(data)-dataBuffer:]
		}

		graph := asciigraph.Plot(data, asciigraph.Height(height-5))
		asciigraph.Clear()
		fmt.Println("Current FPS: " + strconv.Itoa(fps))
		fmt.Println(graph)

		time.Sleep(10 * time.Second)
	}
}
