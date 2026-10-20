    
# Shorthand - installation

# Installation

*shorthand* is a command line program run from a shell like Bash.  
You can find compiled version in the 
[releases](https://github.com/rsdoiel/shorthand/releases/latest) 

## Quick install with curl

The following curl command can be used to run the installer on most
POSIX systems. Programs are installed into `$HOME/bin`. `$HOME/bin` will
need to be in your path. From a shell (or terminal session) run the
following.

~~~
curl https://rsdoiel.github.io/shorthand/installer.sh | sh
~~~

### If you are running macOS or Windows

You may get security warnings, because the programs are not signed. See the notes for the system you are using.

- [INSTALL_NOTES_macOS.md](INSTALL_NOTES_macOS.md)
- [INSTALL_NOTES_Windows.md](INSTALL_NOTES_Windows.md)

## Compiled version

This is generalized instructions for a release. 

Compiled versions are available for macOS (Intel, macOS-x86_64, and Apple silicon, macOS-arm64),
Linux (Intel, Linux-x86_64, and 64 bit ARM such as a Raspberry Pi 4 or 5 or a Pine64, Linux-aarch64),
Raspberry Pi OS 32 bit (Linux-armv7l) and Windows (Intel, Windows-x86_64, and ARM, Windows-arm64).

VERSION_NUMBER is a [symantic version number](http://semver.org/) (e.g. v0.3.0)

For all the released version go to the project page on Github and click latest release

>    https://github.com/rsdoiel/shorthand/releases/latest

| Platform             | Zip Filename                                  |
|----------------------|-----------------------------------------------|
| Windows              | shorthand-VERSION_NUMBER-Windows-x86_64.zip   |
| Windows on ARM       | shorthand-VERSION_NUMBER-Windows-arm64.zip    |
| macOS Intel          | shorthand-VERSION_NUMBER-macOS-x86_64.zip     |
| macOS Apple silicon  | shorthand-VERSION_NUMBER-macOS-arm64.zip      |
| Linux Intel          | shorthand-VERSION_NUMBER-Linux-x86_64.zip     |
| Linux 64 bit ARM     | shorthand-VERSION_NUMBER-Linux-aarch64.zip    |
| Raspberry Pi OS 32 bit | shorthand-VERSION_NUMBER-Linux-armv7l.zip   |


## The basic recipe

1. Download the zip file using your web browser
2. Change to your HOME directory in your shell, command shell or terminal application
3. Unzip the files in the bin folder of the zip archive
4. Make sure the bin folder location is in our path
5. Test



### Mac OS X

Here's an example of the commands run in the Terminal App after 
downloading the zip file in to your "Downloads" folder.

```shell
    cd 
    unzip Downloads/shorthand-*-macOS-arm64.zip bin/*   # or macOS-x86_64 on an Intel Mac
    export PATH="$HOME/bin:$PATH"
    shorthand -version
```

### Windows

Here's an example of the commands run in from the Bash shell on 
Windows 10 after downloading the zip file.

```shell
    cd 
    unzip Downloads/shorthand-*-Windows-x86_64.zip bin/*
    export PATH="$HOME/bin:$PATH"
    shorthand -version
```


### Linux 

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd
    unzip Downloads/shorthand-*-Linux-x86_64.zip bin/*
    export PATH="$HOME/bin:$PATH"
    shorthand -version
```


### Raspberry Pi

The 32 bit Raspberry Pi OS version is for a Raspberry Pi 2 or later (i.e. requires ARM 7 support). On a 64 bit Raspberry Pi OS use the Linux-aarch64 zip file.

1. Download the zip file
2. Change to your HOME directory
3. Unzip the files in bin folder of the zip archive
4. Make sure the new location is in our path
5. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd 
    unzip Downloads/shorthand-*-Linux-armv7l.zip bin/*
    export PATH="$HOME/bin:$PATH"
    shorthand -version
```


## Compiling from source

_shorthand_ needs Go 1.26.4 or later and has no other dependencies. To install the latest release into `$HOME/go/bin`

```shell
    go install github.com/rsdoiel/shorthand/cmd/shorthand@latest
```

or clone the repository and use Make, which also builds the man page. The Makefile uses CMTools,
Pandoc and GNU Make, which are only needed to regenerate `version.go`, the man page and the website.

```shell
    git clone https://github.com/rsdoiel/shorthand
    cd shorthand
    make
    make test
    make install
```
