```shell
go build -ldflags="-s -w" -trimpath -o true.exe .\src\coreutils\true\true.go
upx -9 true.exe
```

```text
                       Ultimate Packer for eXecutables
UPX 5.0.1       Markus Oberhumer, Laszlo Molnar & John Reiser    May 6th 2025

        File size         Ratio      Format      Name
   --------------------   ------   -----------   -----------
1628160 ->    867840   53.30%    win64/pe     false.exe

```

### coreutils命令列表

```
-- 基于(GNU coreutils) 8.30
dpkg -L coreutils | grep -E '^(/usr/bin/|/bin/)' | wc -l
```

* [ ] /bin/cat
* [ ] /bin/chgrp
* [ ] /bin/chmod
* [ ] /bin/chown
* [ ] /bin/cp
* [ ] /bin/date
* [ ] /bin/dd
* [ ] /bin/df
* [ ] /bin/dir
* [x] /bin/echo
* [x] /bin/false
* [ ] /bin/ln
* [ ] /bin/ls
* [ ] /bin/mkdir
* [ ] /bin/mknod
* [ ] /bin/mktemp
* [ ] /bin/mv
* [ ] /bin/pwd
* [ ] /bin/readlink
* [ ] /bin/rm
* [ ] /bin/rmdir
* [x] /bin/sleep
* [ ] /bin/stty
* [ ] /bin/sync
* [ ] /bin/touch
* [x] /bin/true
* [x] /bin/uname
* [ ] /bin/vdir
* [ ] /usr/bin/[
* [ ] /usr/bin/arch
* [ ] /usr/bin/b2sum
* [ ] /usr/bin/base32
* [ ] /usr/bin/base64
* [ ] /usr/bin/basename
* [ ] /usr/bin/chcon
* [ ] /usr/bin/cksum
* [ ] /usr/bin/comm
* [ ] /usr/bin/csplit
* [ ] /usr/bin/cut
* [ ] /usr/bin/dircolors
* [ ] /usr/bin/dirname
* [ ] /usr/bin/du
* [ ] /usr/bin/env
* [ ] /usr/bin/expand
* [ ] /usr/bin/expr
* [ ] /usr/bin/factor
* [ ] /usr/bin/fmt
* [ ] /usr/bin/fold
* [ ] /usr/bin/groups
* [ ] /usr/bin/head
* [ ] /usr/bin/hostid
* [ ] /usr/bin/id
* [ ] /usr/bin/install
* [ ] /usr/bin/join
* [ ] /usr/bin/link
* [ ] /usr/bin/logname
* [ ] /usr/bin/md5sum
* [ ] /usr/bin/mkfifo
* [ ] /usr/bin/nice
* [ ] /usr/bin/nl
* [ ] /usr/bin/nohup
* [ ] /usr/bin/nproc
* [ ] /usr/bin/numfmt
* [ ] /usr/bin/od
* [ ] /usr/bin/paste
* [ ] /usr/bin/pathchk
* [ ] /usr/bin/pinky
* [ ] /usr/bin/pr
* [ ] /usr/bin/printenv
* [ ] /usr/bin/printf
* [ ] /usr/bin/ptx
* [ ] /usr/bin/realpath
* [ ] /usr/bin/runcon
* [ ] /usr/bin/seq
* [ ] /usr/bin/sha1sum
* [ ] /usr/bin/sha224sum
* [ ] /usr/bin/sha256sum
* [ ] /usr/bin/sha384sum
* [ ] /usr/bin/sha512sum
* [ ] /usr/bin/shred
* [ ] /usr/bin/shuf
* [ ] /usr/bin/sort
* [ ] /usr/bin/split
* [ ] /usr/bin/stat
* [ ] /usr/bin/stdbuf
* [ ] /usr/bin/sum
* [ ] /usr/bin/tac
* [ ] /usr/bin/tail
* [ ] /usr/bin/tee
* [ ] /usr/bin/test
* [ ] /usr/bin/timeout
* [ ] /usr/bin/tr
* [ ] /usr/bin/truncate
* [ ] /usr/bin/tsort
* [ ] /usr/bin/tty
* [ ] /usr/bin/unexpand
* [ ] /usr/bin/uniq
* [ ] /usr/bin/unlink
* [ ] /usr/bin/users
* [ ] /usr/bin/wc
* [ ] /usr/bin/who
* [x] /usr/bin/whoami
* [x] /usr/bin/yes
* [ ] /usr/bin/md5sum.textutils

### grep包替代

```bash
scoop install grep
```

### sudo

```bash
scoop install sudo
```

### 测试

* -v
* -vv
* -help
* -help_1
* --help
* --help_1
* --other
* --h
* --

###  

```shell
for cmd in /usr/bin/* /bin/*; do
    if [ -x "$cmd" ]; then
        output=$(dpkg -S "$cmd" 2>/dev/null)
        [ $? -eq 0 ] && echo "$output"
    fi
done | sort -u > out.txt
```

```shell
cat 文件名out.txt | cut -d: -f1 | sort | uniq -c | sort -nr
```