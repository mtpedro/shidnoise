package main

import (
	"fmt"
	"log"
	"strings"

	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/dgvoice"         // for playing sound in vc
	"github.com/bwmarrin/discordgo"       // for using the discord API
	youtube "github.com/kkdai/youtube/v2" // for downloading the youtube videos
	//"github.com/mattn/go-shellwords"      // for parsing the arguments given in the commands
)

var (
	cfg       Config
	que       = make(chan string, 1000) // song que of length 1000.
	backtrace string
)

func failErrNow(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	discord, err := discordgo.New("Bot " + cfg.Token)

	if err != nil {
		log.Fatal(err)
	}

	// Register messageCreate as a callback for the messageCreate events.
	discord.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	//discord.AddHandler(guildCreate)

	// Register ready as a callback for the ready events.
	//discord.AddHandler(ready);

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	if err = discord.Open(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Shidnoise started! Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

// this function is passed to the AddHandler handler above every time
// a new message is created on any channel that the bot has acess to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Debugging message
	// Not Necessary in production
	if m.Content == cfg.Prefix+"test" {
		s.ChannelMessageSend(m.ChannelID, "Shidnoise is online!")
	}

	// If the message is "!horn" join VC and play the song.
	if strings.HasPrefix(m.Content, cfg.Prefix+"play") {

		// find out if it STARTS with "!play" instead of being equal to "!play";

		// write message to channell.
		var sendmsg = func(str string) {
			s.ChannelMessageSend(m.ChannelID, str)
		}

		// two types of possible input:
		// !play <song>

		var args = strings.Split(m.Content, " ")

		if len(args) > 1 {
			// therefor, at least one argument is passed,
			//song = args[1];
			que <- args[1]
		} else {
			sendmsg("please pass a valid argument to !play")
			return
		}

		if use != 0 {
			return
		} // see use.go for explination

		occupy()     // begin use of the channel.
		defer free() // end use.

		// TODO:
		// the plan is to have an function that can play all the songs in the cue.

		backtrace := <-que; 
	play:
		
		videoID := backtrace; 
		client := youtube.Client{}

		video, err := client.GetVideo(videoID)
		if err != nil {
			panic(err)
		}

		formats := video.Formats.WithAudioChannels() // only get videos with audio
		stream, _, err := client.GetStream(video, &formats[0])
		if err != nil {
			log.Fatal(err)
		}

		var fname string = "audio/" + videoID + ".mp3"

		file, err := os.Create(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = io.Copy(file, stream)
		if err != nil {
			panic(err)
		}

		c, err := s.State.Channel(m.ChannelID)
		failErrNow(err)

		g, err := s.State.Guild(c.GuildID)
		failErrNow(err)

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println(FailedJoinVC)
					return
				}
				dgvoice.PlayAudioFile(vc, fname, make(<-chan bool))

				select {
				case c := <-que:
					{
						backtrace = c
						goto play
					}
				default:
					{
						return
					}
				}
			}
		}
		//dgvoice.PlayAudioFile(vc, fname, make(chan bool))
	}
}
