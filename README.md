# Shell Application

## How to run it 

```
./mockshell
```

## List of Commands 

  - ls - Lists the files in the current directory
  - mkdir - Creates one or more directories passed as arguments in the current directory
    
    Example: `mkdir dir1 dir2`
    
  - rm - Removes one or more directories passed as arguments in the current directory
    
    Example: `rm dir1`
  - pwd - Prints the current path
  - cd - Switches context to the directory passed as an argument
  
    Example: `cd dir2`
    
    Note: `cd` does not support nested folders like `/dir1/dir2` at the moment
    
    
### What I'd do as next iteration 
 
  1. Add support for nested directory traversal in `cd`
  2. Add support for nested directory creation in `mkdir`
  3. Add support for nested directory removal in `rm`
  4. Add help option support for all the commands
  5. Add command for creating a file for something like the `touch` command
  6. Add options support for `ls` to display file sizes 
  7. Setup https://github.com/spf13/cobra for cli
  8. Setup https://github.com/c-bata/go-prompt to support autocompletion/suggestion
