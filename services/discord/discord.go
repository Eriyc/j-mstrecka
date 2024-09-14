package discord

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type DiscordConfig struct {
	Token string
}

type DiscordService struct {
	config  *DiscordConfig
	session *discordgo.Session
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func New(config *DiscordConfig) *DiscordService {
	return &DiscordService{
		config: config,
	}
}

func (s *DiscordService) Name() string {
	return "gostrecka/services/discord"
}

func (s *DiscordService) OnStartup(ctx context.Context, options interface{}) error {
	fmt.Printf("Discord service is starting\n")
	if s.config.Token == "" {
		return errors.New("discord token is not set")
	}

	s.ctx, s.cancel = context.WithCancel(ctx)

	err := s.Connect()
	if err != nil {
		return err
	}

	s.wg.Add(1)
	go s.run()

	return nil
}

func (s *DiscordService) Connect() (err error) {
	session, err := discordgo.New(s.config.Token)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	s.session = session

	os := application.Get().Environment().OS
	fmt.Printf("OS: %s\n", os)

	return nil
}

func (s *DiscordService) run() {
	defer s.wg.Done()

	log.Println("Discord service is running")

	// Add your Discord event handlers here
	// For example:
	// s.session.AddHandler(messageCreate)

	// Keep the service running until the context is cancelled
	<-s.ctx.Done()
	log.Println("Discord service is shutting down")
}

func (s *DiscordService) OnShutdown() {
	log.Println("Discord service is shutting down")
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
	if s.session != nil {
		err := s.session.Close()
		if err != nil {
			log.Printf("Error closing Discord session: %v", err)
		}
	}
}

// You can keep your existing Close method if needed
func (s *DiscordService) Close() (string, error) {
	if s.session != nil {
		err := s.session.Close()
		if err != nil {
			return "", errors.New("failed to close discord session")
		}
	}
	return "Successfully closed discord session", nil
}
