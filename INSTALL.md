    
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

## Compiled version

This is generalized instructions for a release. 

Compiled versions are available for Mac OS X (Intel processor, macOS-x86_64), 
Linux (Intel process, linux-x86_64), Windows (Intel processor, Windows-x86_64), 
Rapsberry Pi (arm7 processor, RaspberryPiOS-arm7) and Pine64 (arm64 processor, Linux-aarch64)


VERSION_NUMBER is a [symantic version number](http://semver.org/) (e.g. v0.1.2)


For all the released version go to the project page on Github and click latest release

>    https://github.com/rsdoiel/shorthand/releases/latest


| Platform    | Zip Filename                           |
|-------------|----------------------------------------|
| Windows     | shorthand-VERSION_NUMBER-Windows-x86_64.zip |
| Mac OS X    | shorthand-VERSION_NUMBER-macOS-x86_64.zip  |
| Linux/Intel | shorthand-VERSION_NUMBER-Linux-x86_64.zip   |
| Raspbery Pi | shorthand-VERSION_NUMBER-RaspberryPiOS-arm7.zip |
| Pine64      | shorthand-VERSION_NUMBER-Linux-aarch64.zip   |


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
    unzip Downloads/shorthand-*-macOS-x86_64.zip bin/*
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

Released version is for a Raspberry Pi 2 or later use (i.e. requires ARM 7 support).

1. Download the zip file
2. Change to your HOME directory
3. Unzip the files in bin folder of the zip archive
4. Make sure the new location is in our path
5. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd 
    unzip Downloads/shorthand-*-RaspberryPiOS-arm7.zip bin/*
    export PATH="$HOME/bin:$PATH"
    shorthand -version
```


## Compiling from source

_shorthand_ is "go gettable".  Use the "go get" command to download the dependant packages
as well as _shorthand_'s source code.

```shell
    go get -u github.com/rsdoiel/shorthand/...
```

Or clone the repository and then compile

```shell
    cd
    git clone https://github.com/rsdoiel/shorthand src/github.com/rsdoiel/shorthand
    cd src/github.com/rsdoiel/shorthand
    make
    make test
    make install
```

