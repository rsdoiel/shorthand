    
# Shorthand - installation

# Installation

*shorthand* is a command line program run from a shell like Bash.  
You can find compiled version in the 
[releases](https://github.com/rsdoiel/shorthand/releases/latest) 

## Compiled version

This is generalized instructions for a release. 

Compiled versions are available for Mac OS X (amd64 processor, macosx-amd64), 
Linux (amd64 process, linux-amd64), Windows (amd64 processor, windows-amd64), 
Rapsberry Pi (arm7 processor, raspbian-arm7) and Pine64 (arm64 processor, linux-arm64)


VERSION_NUMBER is a [symantic version number](http://semver.org/) (e.g. v0.1.2)


For all the released version go to the project page on Github and click latest release

>    https://github.com/rsdoiel/shorthand/releases/latest


| Platform    | Zip Filename                           |
|-------------|----------------------------------------|
| Windows     | shorthand-VERSION_NUMBER-windows-amd64.zip |
| Mac OS X    | shorthand-VERSION_NUMBER-macosx-amd64.zip  |
| Linux/Intel | shorthand-VERSION_NUMBER-linux-amd64.zip   |
| Raspbery Pi | shorthand-VERSION_NUMBER-raspbian-arm7.zip |
| Pine64      | shorthand-VERSION_NUMBER-linux-arm64.zip   |


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
    unzip Downloads/shorthand-*-macosx-amd64.zip bin/*
    export PATH="$HOME/bin:$PATH"
    shorthand -version
```

### Windows

Here's an example of the commands run in from the Bash shell on 
Windows 10 after downloading the zip file.

```shell
    cd 
    unzip Downloads/shorthand-*-windows-amd64.zip bin/*
    export PATH="$HOME/bin:$PATH"
    shorthand -version
```


### Linux 

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd
    unzip Downloads/shorthand-*-linux-amd64.zip bin/*
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
    unzip Downloads/shorthand-*-raspbian-arm7.zip bin/*
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

