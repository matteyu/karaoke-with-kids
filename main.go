package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func main() {
	lyrics := handleFileParse()

	chester := slacker.NewClient(os.Getenv("SLACK_TOKEN_CHESTER"))
	pablo := slacker.NewClient(os.Getenv("SLACK_TOKEN_PABLO"))

	pabloRTM := pablo.RTM()
	pabrtm := pabloRTM.NewRTM()
	go pabrtm.ManageConnection()

	definition := &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Typing()
			rtm := botCtx.RTM()
			channel := botCtx.Event().Channel
			client := botCtx.Client()

			// userTalking := botCtx.Event().User
			// userInfo, _ := chester.GetUserInfo(userTalking)
			// fmt.Printf("%v\n", userInfo.RealName)

			keys := make([]int, 0)
			for k := range lyrics {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			for _, k := range keys {
				lyric := strings.Split(lyrics[k], ":")
				i, err := strconv.Atoi(lyric[0])
				if err != nil {
					// handle error
					fmt.Println(err)
					os.Exit(2)
				}
				if lyric[1] != "image" {
					if k%2 == 0 {
						rtm.SendMessage(rtm.NewOutgoingMessage(lyric[1], channel))
					} else {
						pabrtm.SendMessage(rtm.NewOutgoingMessage(lyric[1], channel))
					}
				} else {
					client.UploadFile(slack.FileUploadParameters{File: "goodtime.gif", Channels: []string{channel}})
				}
				time.Sleep(time.Duration(i) * time.Second)
			}

		},
	}

	chester.Command("what is love?", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := chester.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func handleFileParse() map[int]string {
	file, errFile := os.Open("whatislove.txt")
	if errFile != nil {
		log.Fatal(errFile)
	}

	count := 0
	scanner := bufio.NewScanner(file)
	lyrics := make(map[int]string)

	for scanner.Scan() {
		lyrics[count] = scanner.Text()
		count++
	}

	return lyrics
}
