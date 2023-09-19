# terraform-version-inspect
A CLI application to determine the required terraform version

Does the shallow terraform project parsing to provide the required version.
The list of available versions is taken from https://releases.hashicorp.com/.

### Usage
        terraform-version-inspect [flags]

### Flags
            --debug             turn on debug logging
            --dir string        path that contains terraform configuration files (default ".")
            --registry string   ensure the terraform image being available in the specified registry
        -h, --help              help for terraform-version-inspect

### Example

```bash
$ cat test.tf
terraform {
    required_version = "~> 1.3, < 1.4"
}
terraform {
    required_version = ">= 1.2.31"
}
$ terraform-version-inspect --dir '.'
1.3.9
```
