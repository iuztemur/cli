## tkn completion

Prints shell completion scripts

### Usage

```
tkn completion [SHELL]
```

### Synopsis


This command prints shell completion code which must be evaluated to provide
interactive completion

Supported Shells:
	- bash
	- zsh


### Examples


  # generate completion code for bash
  source <(tkn completion bash)

  # generate completion code for zsh
  source <(tkn completion zsh)


### Options

```
  -h, --help   help for completion
```

### SEE ALSO

* [tkn](tkn.md)	 - CLI for tekton pipelines

