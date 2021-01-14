# password-manager

The Password manager is a full-featured command-line interface (CLI) tool for helping you store all of your login credentials. Password managers helps you to have easy access to all of your password - all you need is to remember your email address and master password. 

The CLI is written with Golang and can be run on Windows, macOS, and Linux distributions.


>Note: This password manager was made as a project and is NOT intended for actual use. Please use more sophisticated and well-tested/trusted password managers to store sensitive data.


## Download/Install
If you already have the golang installed on your system, you can install the CLI using :
```
    go get github.com/nnachevv/passmag
```


## Documentation
The Password manager CLI is self-documented with --help content and examples for every command. You should start exploring the CLI by using the global --help option:

`passmag --help`

This option will list all available commands that you can use with the CLI.

Additionally, you can run the --help option on a specific command to learn more about it:

`passmag list --help`
`passmag init --help`

## Detailed documentation

`passmag init` - allows any user to register for the Password manager vault service.

`passmag login` - login in your registred user.  After successfully logging into the CLI a session key will be returned. This session key is necessary to perform any commands that require your vault to be unlocked (list, get, edit, etc.).

`passmag add` - allows you to add a entry and password in your vault. Password can be randomly generated.

`passmag edit` - allows to edit your added password

`passmag get` - get password from already added passwords

`passmag cp` - copy password to clipboard.  

> **WARNING**: For some operation system maybe this could not work.** 

`passmag list` - lists all current passwords from vault

`passmag logout` - delete vault from filesystem.



