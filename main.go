package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)


type Directory struct{

	Path string
	parent *Directory
	contents []Content
	name string

}

type Content struct{
	itemName string
	itemType string
	folder *Directory
	parent *Directory
}

var Homedirectory Directory 

func Init(){
	Homedirectory = Directory{Path:"/", name:""}
}

func main() {
	// Create new reader
	reader := bufio.NewReader(os.Stdin)
	Init() // Initialise HomeDirectory
	// Set CurrentDirectory to HomeDirectory
	var CurrentDirectory *Directory
	CurrentDirectory = &Homedirectory

	// Infinite loop that listen to user input
	for {
		fmt.Print(CurrentDirectory.Path + ":$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Call runCommand which exectues the appropriate behaviour from user command
		CurrentDirectory, err = runCommand(cmdString, CurrentDirectory)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}


func get_dirs_and_files_in_current_path(CurrentDirectory *Directory) (string, string) {

	dirs_string := "DIRS: "
	files_string := "FILES: "

	for _, item := range CurrentDirectory.contents{
		if item.itemType=="dir"{
			dirs_string = dirs_string + "\t" + item.itemName
		}

		if item.itemType=="file"{
			files_string = files_string + "\t" + item.itemName
		}
	}

	return dirs_string, files_string

}

func get_current_path(CurrentDirectory *Directory,current_path string) string {

	if CurrentDirectory.parent == nil{
		return current_path
	}
	current_path = CurrentDirectory.Path + current_path
	return get_current_path(CurrentDirectory.parent, current_path)

}

func check_and_traverse_user_input_path(CurrentDirectory *Directory, user_input string) (*Directory, error){

	dir_found := false
	for _, content := range CurrentDirectory.contents{
		if content.itemType != "dir"{
			continue
		}

		if user_input == content.itemName{
			dir_found = true
			CurrentDirectory = content.folder
		}
	}

	if dir_found == false{
		return CurrentDirectory, errors.New("ERR: Invalid Path")
	}

	return CurrentDirectory, nil
}

func input_contains_dir(name string, user_input []string) bool {
	for _, item_name := range user_input{
		if item_name==name{
			return true
		}
	}
	return false
}

func check_if_directory_exists(name string, contents []Content) bool {

	for _,content := range contents{
		if content.itemName == name{
			return true
		}
	}

	return false
}

func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
			m[item] = true
	}

	for _, item := range a {
			if _, ok := m[item]; !ok {
					diff = append(diff, item)
			}
	}
	return
}

func runCommand(commandStr string, CurrentDirectory *Directory) (*Directory, error) {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)

	// To handle empty input when user hits enter key
	if len(arrCommandStr)==0{
		return CurrentDirectory, nil
	}

	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)

	case "session":
		if (arrCommandStr[1]=="clear"){		
			Init()
			CurrentDirectory = &Homedirectory
			return CurrentDirectory, nil
		}

		return CurrentDirectory, errors.New("ERR: invalid command " + arrCommandStr[1])

	case "ls":
		// Show error if user passes more arguments to this command
		if len(arrCommandStr) > 1 {
			return CurrentDirectory, errors.New("ls does not take any arguments")
		}

		// get list of files and directories as a printable string
		dirs_string, files_string := get_dirs_and_files_in_current_path(CurrentDirectory)
		fmt.Println(dirs_string)
		fmt.Println(files_string)
		return CurrentDirectory, nil
	
	case "cd":
		if len(arrCommandStr) > 2 {
			return CurrentDirectory, errors.New("ERR: cd requires only 1 argument. Got " + strconv.Itoa(len(arrCommandStr[1:])))
		}

		// return user to HomeDirectory when second argument is empty or / or is not given
		if (len(arrCommandStr)==1 || arrCommandStr[1] == "/" || arrCommandStr[1] == ""){
			CurrentDirectory = &Homedirectory
			return CurrentDirectory, nil
		}

		// Return to parent directory if possible otherwise print message
		if (arrCommandStr[1]==".."){
			if CurrentDirectory.parent==nil{
				return CurrentDirectory, errors.New("ERR: This is the top most directory.")
			}

			CurrentDirectory = CurrentDirectory.parent

			return CurrentDirectory, nil
		}

		// Check to see if user input path is valid and go to directory
		// Note: Need to implement nested path traversal
		CurrentDirectory, err := check_and_traverse_user_input_path(CurrentDirectory, arrCommandStr[1])
		return CurrentDirectory, err

	case "mkdir":
		if len(arrCommandStr) < 2 {
			return CurrentDirectory, errors.New("ERR: mkdir requires atleast 1 argument")
		}
		var invalid_dir_names, valid_dir_names [] string
		// Loop through list of directories,check if it can be created and then create 
		// the object and add it to the contents of the current folder
		for _, dirName := range arrCommandStr[1:]{
			dir_exists := check_if_directory_exists(dirName, CurrentDirectory.contents)
			if dir_exists==true{
				invalid_dir_names = append(invalid_dir_names, dirName)
			}else{			
				newDirectory := Directory{
								name: dirName, 
								parent: CurrentDirectory,
								Path: get_current_path(CurrentDirectory, "/") + dirName}
				newContent := Content{
								itemName: dirName,
								itemType: "dir",
								folder: &newDirectory,
								parent: CurrentDirectory}
				CurrentDirectory.contents = append(CurrentDirectory.contents, newContent)
				valid_dir_names = append(valid_dir_names, dirName)
			}
		}

		// Print valid directory names that got created
		if len(valid_dir_names) > 0{
			fmt.Println("SUCC: " + strings.Join(valid_dir_names, " "))
		}

		// Return invalid directory names as a string
		if len(invalid_dir_names) > 0{
			return CurrentDirectory, errors.New("ERR: Dirs already exist - " + strings.Join(invalid_dir_names, " "))
		}

		return CurrentDirectory, nil

	case "rm":
		// Check to see if user has passed atleast one argument
		if len(arrCommandStr) < 2 {
			return CurrentDirectory, errors.New("ERR: rm requires atleast 1 argument")
		}

		var indexes_to_delete []int
		var valid_dir_name []string
		// Loop through contents of this folder to 
		for index, dir := range CurrentDirectory.contents{
			if input_contains_dir(dir.itemName, arrCommandStr[1:])==true{
				indexes_to_delete = append(indexes_to_delete, index)
				valid_dir_name = append(valid_dir_name, dir.itemName)
			}
		}

		// Delete folders that have been deleted from the contents array in the directory
		for _, position := range indexes_to_delete{
			CurrentDirectory.contents = append(CurrentDirectory.contents[:position], CurrentDirectory.contents[position+1:]...)
		}
	
		// Filter out the folder names that don't exist and display error message
		invalid_dir_names := Difference(arrCommandStr[1:], valid_dir_name)
		if len(invalid_dir_names) > 0{
			return CurrentDirectory, errors.New("ERR: Invalid paths " + strings.Join(invalid_dir_names, " "))
		}

	case "pwd":
		// Print error message is user passes arguments to pwd
		if len(arrCommandStr) > 1{
			return CurrentDirectory, errors.New("pwd does not take any arguments")
		}

		// Print the current path from CurrentDirectory.Path
		fmt.Println("PATH: " + CurrentDirectory.Path)
		return CurrentDirectory, nil
	
	default:
		return CurrentDirectory, errors.New("ERR: Cannot Recognize Input")
		
	}

	return CurrentDirectory, nil

}

