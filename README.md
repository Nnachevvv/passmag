# password-manager
[![Build Status](https://travis-ci.com/Nnachevvv/passmag.svg?branch=main)](https://travis-ci.com/Nnachevvv/passmag)[![codecov](https://codecov.io/gh/Nnachevvv/passmag/branch/main/graph/badge.svg?token=AO0YBN997F)](https://codecov.io/gh/Nnachevvv/passmag)


The Password manager is a full-featured command-line interface (CLI) tool for helping you store all of your login credentials. Password managers helps you to have easy access to all of your password - all you need is to remember your email address and master password. 

The CLI is written with Golang and can be run on Windows, macOS, and Linux distributions.


>Note: This password manager was made as a project and is NOT intended for actual use. Please use more sophisticated and well-tested/trusted password managers to store sensitive data.


## Download/Install
If you already have the golang installed on your system, you can install the CLI using :
```
    go get github.com/Nnachevvv/passmag
```

MongoDB should be installed as prerequisite.
## Running
If you've installed it:

    $ passmag

Otherwise, from the build directory:

    $ ./passmag



## Documentation
The Password manager CLI is self-documented with --help content and examples for every command. You should start exploring the CLI by using the global --help option:
    
    $ passmag --help

This option will list all available commands that you can use with the CLI.

Additionally, you can run the --help option on a specific command to learn more about it:

    $ passmag init --help
    $ passmag list --help

## Managing Your Vault

#### Init new vault
To init and register new account in Password manager use:
    
    $ passmag init

This command will add encrypted email address with your password. 

#### Login and download your vault
After you init your password you should login to download your vault locally:

    $ passmag login

 After sussecfuly logged session key will be generated and your vault will be encrypted with your password and session_key. You should export SESSION_KEY variable or pass it every time to any command that require your vault to be unlocked (list, get, edit, etc.).

#### Manage your vault
There several options to manage your vault:

To add password , or randomly generate new one:

    $ passmag add

Rename an already added password:

    $ passmag edit

Change an already added password:

    $ passmag edit

Get password from your vault:

    $ passmag get 

Copy password without expose to clipboard:

    $ passmag cp   

List all passwords from your vault:

    $ passmag list 

Logout and delete vault from filesystem

    $ passmag logout 

