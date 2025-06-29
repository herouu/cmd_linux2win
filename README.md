
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
