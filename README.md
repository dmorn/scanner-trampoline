# scanner-trampoline
Used to execute a command on each scan. Allows to configure the command
to be executed and how the scanned string should be trimmed.

## Usage
- execute the binary (I'll setup a release soon)
- finish the configuration setup
- 1: focus the terminal window
- scan something
- goto: 1

## Scanner Configuration
- scanner must emit a Carriage Return / Line Feed character after scanning the
  barcode
- the Zebra DS22 can be configured with the relative option found in the manual
  under "Miscellaneous Scanner Parameters -> Enter Key"
