package akin

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/leonhfr/akin/internal/akin/files"
	"github.com/leonhfr/akin/internal/akin/lexer"
	"github.com/leonhfr/akin/internal/akin/parser"
	"github.com/leonhfr/akin/internal/akin/utils"
)

// TODO, Verbose, DryRun
// Process synchronizes the local Markdown state to Anki.
func (akin *Akin) Process() error {
	var wg sync.WaitGroup

	numJobs := runtime.NumCPU()
	mds := make(chan string)

	go files.Files(akin.config, mds, akin.errors)

	for w := 0; w < numJobs; w++ {
		wg.Add(1)
		go akin.worker(w, &wg, mds)
	}

	wg.Wait()

	return nil
}

func (akin *Akin) worker(id int, wg *sync.WaitGroup, mds <-chan string) {
	defer wg.Done()

	for md := range mds {
		akin.fileWorker(md)
	}
}

func (akin *Akin) fileWorker(md string) {
	l, err := lexer.New(md, akin.errors)
	if err != nil {
		akin.errors <- err
		return
	}

	tokens, err := l.Lex()
	if err != nil {
		akin.errors <- err
		return
	}

	abs := utils.AbsoluteDirPath(md)
	p := parser.New(tokens, abs, akin.notes, akin.medias, akin.errors)
	err = p.Parse()
	if err != nil {
		akin.errors <- err
		return
	}

	err = akin.createDeck(p.Deck)
	if err != nil {
		akin.errors <- err
		return
	}
}

func (akin *Akin) createDeck(deck string) error {
	_, ok := akin.deckMap[deck]
	if ok {
		return nil
	}
	akin.mutex.Lock()
	defer akin.mutex.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	_, err := akin.client.CreateDeck(ctx, deck)
	return err
}
