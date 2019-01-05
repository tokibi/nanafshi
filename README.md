# nanafshi

## Description

WIP

## Usage

```console
$ nanafshi -c config.yml /mnt/nanafshi

$ tree /mnt/nanafshi
/mnt/nanafshi
└── dir
    ├── now
    ├── example
    └── writable

$ umount /mnt/nanafshi
```

### Configure example

```yaml
shell: /bin/bash
services:
  - name: dir
    files:
      - name: example
        read:
          command: curl http://www.example.com/

      - name: now
        read:
          command: date

      - name: writable
        read: 
          command: |
            touch /tmp/file
            cat /tmp/file
        write: 
          async: true
          command: echo $FUSE_STDIN >> /tmp/file
```

```console
$ cat /mnt/nanafshi/dir/now
Sun Jan  6 03:54:44 JST 2019
```

### Environment variables

The following environment variables can be used in the command.

|env|description|
|---|---|
|FUSE_FILENAME|Filename|
|FUSE_STDIN|Input data written to file (write only)|
|FUSE_OPENPID|PID of the process that opened this file|
|FUSE_OPENUID|UID of the process that opened this file|
|FUSE_OPENGID|GID of the process that opened this file|

## TODO

- [ ] Configure authority (mode, owner)
- [ ] Logger
- [ ] Read cache
- [ ] Support interactive command
- [ ] Injection measures
