package main

import (
	"fmt"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var (
	cfg Config
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
	if m.Content == cfg.Prefix+"horn" {

		c, err := s.State.Channel(m.ChannelID)
		failErrNow(err);

		g, err := s.State.Guild(c.GuildID)
		failErrNow(err);

		var fname string = "audio/song.mkv"

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true);
				if err != nil {
					fmt.Println(FailedJoinVC);
					return; 
				}
				dgvoice.PlayAudioFile(vc, fname, make(<-chan bool)); 

				return
			}
		}
		//dgvoice.PlayAudioFile(vc, fname, make(chan bool))
	}
}
