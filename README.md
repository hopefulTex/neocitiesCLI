# neocities-cli

A command line interface for Neocities

## Usage

### TUI

```bash
neocities
```

### CLI

```bash
neocities [command] [options]
```

### Commands

#### `push`

Uploads files to Neocities.

```bash
neocities push [directory]
```

#### `upload`

Uploads files to Neocities.

```bash
neocities upload [path/to/first.file, path/to/second.file]
```

#### `delete`

Deletes files from Neocities.

```bash
neocities delete [file.txt, file2.txt]
```

#### `list`

Lists files on Neocities.

```bash
neocities list [directory]
```

#### `info`

Displays information about a site.

```bash
neocities info
```
```bash
neocities info [sitename]
```


#### `config`

Set default account, List accounts, Login, and Reset the configuration file.

```bash
neocities config set domain example.com
```
```bash
neocities config get domain
```
```bash
neocities config reset
```

get the path to the configuration file

```bash
neocities config path
```

## Configuration

The configuration file is located at `~/.config/neocities/config.json`.

The configuration file contains the following fields:

- `api_key`: The API key for the Neocities API.
- `domain`: The domain for the Neocities site.
- `pw`: The password for the Neocities site.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 

## Acknowledgments

- [Neocities](https://neocities.org/)
- [Neocities API](https://neocities.org/api)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/main/examples)
- [Bubble Tea Documentation](https://pkg.go.dev/github.com/charmbracelet/bubbletea)
