package akin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/leonhfr/akin/cfg"
	"github.com/leonhfr/anki-connect-go"
)

type Akin struct {
	config   *cfg.Config                // Akin configuration
	client   *anki.Client               // Anki client
	mutex    *sync.Mutex                // Anki client mutex
	notes    chan *anki.NoteInput       // Notes input channel
	medias   chan *anki.StoreMediaInput // Medias input channel
	errors   chan error                 // Error channel
	deckMap  map[string]bool            // Existing decks
	mediaMap map[string]bool            // Existing medias
}

func New(config *cfg.Config) *Akin {
	akin := &Akin{
		config: config,
		client: anki.NewClient(config.URL),
		notes:  make(chan *anki.NoteInput),
		medias: make(chan *anki.StoreMediaInput),
		errors: make(chan error),
		mutex:  &sync.Mutex{},
	}
	go akin.notesWorker()
	go akin.mediasWorker()
	go akin.logger()
	return akin
}

func (akin *Akin) Dispose() {
	close(akin.notes)
	close(akin.medias)
	close(akin.errors)
}

func (akin *Akin) logger() {
	for err := range akin.errors {
		fmt.Println(err)
	}
}

func (akin *Akin) notesWorker() {
	notes := []anki.NoteInput{}
	timer := time.NewTimer(time.Second)
	for {
		select {
		case <-timer.C:
			// Partial flush due to time
			akin.createNotes(notes)
			notes = nil
			timer.Reset(time.Second)
		case note := <-akin.notes:
			notes = append(notes, *note)
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time.Second)
			if len(notes) < 5 {
				continue
			}
			akin.createNotes(notes)
			notes = nil
		}
	}
}

func (akin *Akin) createNotes(notes []anki.NoteInput) error {
	akin.mutex.Lock()
	defer akin.mutex.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := akin.client.AddNotes(ctx, notes)
	return err
}

func (akin *Akin) mediasWorker() {
	for media := range akin.medias {
		err := akin.storeMedia(media)
		if err != nil {
			akin.errors <- err
		}
	}
}

func (akin *Akin) storeMedia(media *anki.StoreMediaInput) error {
	_, ok := akin.mediaMap[media.Filename]
	if ok {
		return nil
	}
	akin.mutex.Lock()
	defer akin.mutex.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := akin.client.StoreMediaByData(ctx, *media)
	return err
}
