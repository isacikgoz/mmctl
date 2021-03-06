## mmctl command create

Create a custom slash command

### Synopsis

Create a custom slash command for the specified team.

```
mmctl command create [team] [flags]
```

### Examples

```
  command create myteam --title MyCommand --description "My Command Description" --trigger-word mycommand --url http://localhost:8000/my-slash-handler --creator myusername --response-username my-bot-username --icon http://localhost:8000/my-slash-handler-bot-icon.png --autocomplete --post
```

### Options

```
      --autocomplete               Show Command in autocomplete list
      --autocompleteDesc string    Short Command Description for autocomplete list
      --autocompleteHint string    Command Arguments displayed as help in autocomplete list
      --creator string             Command Creator's Username (required)
      --description string         Command Description
  -h, --help                       help for create
      --icon string                Command Icon URL
      --post                       Use POST method for Callback URL
      --response-username string   Command Response Username
      --title string               Command Title
      --trigger-word string        Command Trigger Word (required)
      --url string                 Command Callback URL (required)
```

### Options inherited from parent commands

```
      --format string   the format of the command output [plain, json] (default "plain")
```

### SEE ALSO

* [mmctl command](mmctl_command.md)	 - Management of slash commands

