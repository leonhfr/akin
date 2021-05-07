# Roadmap

Those are all the features I'm thinking of implementing for akin. Just because something is on this list doesn't mean it'll get done. I'm building this tool both for my own use and as a side project to learn the Go language, so the implementation of future features depends on whether they serve those aims.

## Before first release

- [ ] Unit tests.
- [ ] Documentation.

## Next features

- [ ] Local database with existing notes (boltdb), check for preexistence when creating new notes
- [ ] Store the files' last modification time in the local database and check it on the next run, parse only the files that have been created or updated since the last run
- [ ] If the card has been updated in Anki but not in the repository, don't update. If the card has been updated in the repository, this change takes priority (but if the card has also been modified in anki, it would be treated as a new card instead of a replacement)
- [ ] Tags management. See the [notes](https://github.com/FooSoft/anki-connect/blob/master/actions/notes.md) actions
- [ ] Make my own anki repository with my cards public on Github
- [ ] Release a pre-compiled binary instead of relying on `go get`
- [ ] Provide an option so that if a card is deleted from the markdown repository, it will be deleted from anki
- [ ] Report count of changes at the end of the CLI execution
- [ ] Update the cards in place, to preserve the learning metadata

## Maybe one day

- [ ] More granular configuration for clean up.
- [ ] Deck configs. See the [decks](https://github.com/FooSoft/anki-connect/blob/master/actions/decks.md) actions
- [ ] Configuration as code (look for `config.yaml` in root directory). Let's try to do this without `viper`.
- [ ] Stdin input.
- [ ] CSS: dark and light themes
- [ ] Template/Models management. See the [models](https://github.com/FooSoft/anki-connect/blob/master/actions/models.md) actions
- [ ] Statistics: the capability for akin to extract and store the cards statistics from Anki, as well as the capability to restore them. See the [statistics](https://github.com/FooSoft/anki-connect/blob/master/actions/statistics.md) actions
- [ ] User profile management. See the [miscellaneous](https://github.com/FooSoft/anki-connect/blob/master/actions/miscellaneous.md) actions
- [ ] Import from Anki: instead of exporting a local state of Markdown files, import cards from Anki and build the local state. This would be great to easily onboard users on the tool.
- [ ] Implement go releases on tag push via Github Actions and [Goreleaser](https://goreleaser.com/ci/actions/)
