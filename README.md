# LIME - Simple Lightweight editor

(Still under development)

![Initial Screen](./assets/images/start.png)

### Installation
```
$ git clone github.com/DeeStarks/lime && \
$ cd lime && make
```

### Usage
- `./lime <filename>`: start the editor

### Configuration
Configuration files are located in `configs` directory.
- `notification.go`: configuration file for notifications. ***(NB: the `utils.PlaySound` function have issues playing sounds, so the configurations aren't currently used anywhere. The bug will be fixed later on.)***
- `editor.go`: configuration file for the editor (e.g. Tab size)
- `version.go`