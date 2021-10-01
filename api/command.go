package api

import (
	"errors"
	db "github.com/diantanjung/wecom/db/sqlc"
	"github.com/diantanjung/wecom/token"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

type commandResponse struct {
	Path    string `json:"path"`
	Command string `json:"command"`
	Message string `json:"message"`
}

func (server *Server) RunCommand(ctx *gin.Context) {
	if !server.isUserHasDir(ctx) {
		err := errors.New("Directory doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	dir := "/" + ctx.Param("dir") + "/"
	path := server.config.BinPath + dir
	runner := path + ctx.Param("cmd")

	var args []string
	query := ctx.Request.URL.Query()

	for key, val := range query {
		args = append(args, "-"+key)
		args = append(args, val[0])
	}

	cmd := exec.Command(runner, args[0:]...)
	cmd.Dir = path
	stdout, err := cmd.Output()

	command := runner + " " + strings.Join(args, " ")

	var message string
	if err != nil {
		message = "Error! " + err.Error()
	} else {
		message = "Success! " + string(stdout)
	}
	res := commandResponse{
		Path:    path,
		Command: command,
		Message: message,
	}

	ctx.JSON(http.StatusOK, res)
}

type getSourceCodeResponse struct {
	MainStr  string `json:"main_str"`
	GoModStr string `json:"gomod_str"`
}

func (server *Server) GetSourceCode(ctx *gin.Context) {
	if !server.isUserHasDir(ctx) {
		err := errors.New("Directory doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	dir := ctx.Param("dir")
	cmd := ctx.Param("cmd")

	// main.go file
	mainFile := server.config.BinPath + "/" + dir + "/" + cmd + ".src/main.go"
	dirPath := server.config.BinPath + "/" + dir
	mainStr, err := ioutil.ReadFile(mainFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//go.mod file
	goModStr, err := ioutil.ReadFile(dirPath + "/" + cmd + ".src/go.mod")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := getSourceCodeResponse{
		MainStr:  strings.Trim(string(mainStr), " "),
		GoModStr: strings.Trim(string(goModStr), " "),
	}
	ctx.JSON(http.StatusOK, res)
}

type getFileContentResponse struct {
	FileStr string `json:"file_str"`
}

func (server *Server) GetFileContent(ctx *gin.Context) {
	if !server.isUserHasDir(ctx) {
		err := errors.New("Directory doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	dir := ctx.Param("dir")
	cmd := ctx.Param("cmd")
	file := ctx.Param("file")

	// file
	filePath := server.config.BinPath + "/" + dir + "/" + cmd + ".src/" + file
	fileString, err := ioutil.ReadFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := getFileContentResponse{
		FileStr: string(fileString),
	}
	ctx.JSON(http.StatusOK, res)
}

type updateFileContentRequest struct {
	FileStr string `json:"file_str" binding:"required"`
}

func (server *Server) UpdateFileContent(ctx *gin.Context) {
	if !server.isUserHasDir(ctx) {
		err := errors.New("Directory doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	var req updateFileContentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dir := ctx.Param("dir")
	cmd := ctx.Param("cmd")
	file := ctx.Param("file")

	pathFile := server.config.BinPath + "/" + dir + "/" + cmd + ".src/" + file

	err := ioutil.WriteFile(pathFile, []byte(req.FileStr), 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := commandResponse{
		Message: "Success update file",
	}

	ctx.JSON(http.StatusOK, res)
}

type updateSourceCodeRequest struct {
	MainStr  string `json:"main_str" binding:"required"`
	GoModStr string `json:"gomod_str" binding:"required"`
}

func (server *Server) UpdateSourceCode(ctx *gin.Context) {
	if !server.isUserHasDir(ctx) {
		err := errors.New("Directory doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	var req updateSourceCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dir := ctx.Param("dir")
	cmd := ctx.Param("cmd")

	mainFile := server.config.BinPath + "/" + dir + "/" + cmd + ".src/main.go"
	dirPath := server.config.BinPath + "/" + dir

	err := ioutil.WriteFile(mainFile, []byte(req.MainStr), 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = ioutil.WriteFile(dirPath+"/"+cmd+".src/go.mod", []byte(req.GoModStr), 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	exeCmd := exec.Command("/usr/local/go/bin/go", "build", "-o", cmd, mainFile)
	exeCmd.Dir = dirPath
	stdout, err := exeCmd.Output()

	var message string
	if err != nil {
		message = "Error! " + err.Error()
	} else {
		message = "Success! " + string(stdout)
	}
	res := commandResponse{
		Message: message,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateCommand(ctx *gin.Context) {
	dir := ctx.Param("dir")
	cmd := ctx.Param("cmd")
	// mkdir newdir ; cp -r master/main.src newdir/newcommand.src ; cp master/main newdir/newcommand
	exeCmd := exec.Command("mkdir", dir)
	exeCmd.Dir = server.config.BinPath
	_, _ = exeCmd.Output()

	exeCmd = exec.Command("cp", "-r", "master/main.src", dir+"/"+cmd+".src")
	exeCmd.Dir = server.config.BinPath
	_, _ = exeCmd.Output()

	exeCmd = exec.Command("cp", "master/main", dir+"/"+cmd)
	exeCmd.Dir = server.config.BinPath
	_, _ = exeCmd.Output()

	if !server.isUserHasDir(ctx) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		req := db.CreateUserDirParams{
			UserID: authPayload.UserID,
			Name:   dir,
		}
		_, err := server.querier.CreateUserDir(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	res := commandResponse{
		Message: "Command successfuly created!",
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) DeleteCommand(ctx *gin.Context) {
	if !server.isUserHasDir(ctx) {
		err := errors.New("Directory doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	dir := ctx.Param("dir")
	cmd := ctx.Param("cmd")

	exeCmd := exec.Command("rm", "-r", dir+"/"+cmd+".src")
	exeCmd.Dir = server.config.BinPath
	_, _ = exeCmd.Output()

	exeCmd = exec.Command("rm", dir+"/"+cmd)
	exeCmd.Dir = server.config.BinPath
	_, _ = exeCmd.Output()

	res := commandResponse{
		Message: "Command successfuly deleted!",
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) isUserHasDir(ctx *gin.Context) (res bool) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	req := db.CheckUserDirParams{
		UserID: authPayload.UserID,
		Name:   ctx.Param("dir"),
	}
	_, err := server.querier.CheckUserDir(ctx, req)
	if err == nil {
		return true
	}
	return false
}

func (server *Server) isUserHasDirectory(ctx *gin.Context, userId int64, dirName string) (res bool) {
	req := db.CheckUserDirParams{
		UserID: userId,
		Name:   dirName,
	}
	_, err := server.querier.CheckUserDir(ctx, req)
	if err == nil {
		return true
	}
	return false
}

type getCommandsResponse struct {
	Dir string `json:"dir"`
	Cmd string `json:"cmd"`
}

func (server *Server) GetCommands(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userDirs, err := server.querier.GetUserDirs(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	userDirsStr := "|"

	for _, dir := range userDirs {
		userDirsStr += dir.Name + "|"
	}

	res := []getCommandsResponse{}
	dirs, err := ioutil.ReadDir(server.config.BinPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	//var strCommand string
	var cmd, dirStr string
	for _, dir := range dirs {
		if dir.IsDir() && strings.Index(userDirsStr, "|"+dir.Name()+"|") > -1 {
			files, _ := ioutil.ReadDir(server.config.BinPath + "/" + dir.Name())
			dirStr = dir.Name()
			for _, file := range files {
				if !file.IsDir() && strings.Index(file.Name(), ".") < 0 {
					cmd = file.Name()
					if cmd != "" && dirStr != "" {
						res = append(res, getCommandsResponse{Cmd: cmd, Dir: dirStr})
					}
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, res)
}

//func (server *Server) Terminal(ctx *gin.Context) {
//	if !server.isUserHasDir(ctx) {
//		err := errors.New("Directory doesn't belong to the authenticated user")
//		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
//		return
//	}
//
//	dir := ctx.Param("dir")
//	cmd := ctx.Param("cmd")
//	exe := ctx.Param("exe")
//	mainFile := server.config.BinPath + "/" + dir + "/" + cmd + ".src/main.go"
//	dirRunner := server.config.BinPath + "/" + dir
//	dirSrc := server.config.BinPath + "/" + dir + "/" + cmd + ".src"
//
//	var commandStr string
//	var commandDir string
//	var args []string
//
//	switch exe {
//	case "ls":
//		commandStr = "ls"
//		commandDir = dirSrc
//	case "code":
//		commandStr = "ls"
//		commandDir = dirSrc
//	case "build":
//		commandStr = "/usr/local/go/bin/go"
//		args = append(args,"build", "-o", ctx.Param("cmd"), mainFile)
//		commandDir = dirRunner
//	case "run":
//		commandStr = dirRunner + "/" + cmd
//		commandDir = dirRunner
//	default:
//		ctx.JSON(http.StatusOK, commandResponse{Message: "Command not found. Try running `help`."})
//	}
//
//	exeCmd := exec.Command(commandStr, args[0:]...)
//	exeCmd.Dir = commandDir
//	stdout, err := exeCmd.Output()
//
//	var message string
//	if err != nil {
//		message = err.Error()
//	} else {
//		if len(stdout) > 0 {
//			message = string(stdout)
//		}else{
//			message = "Succes to execute `" + exe + "`"
//		}
//
//	}
//	res := commandResponse{
//		Message: message,
//	}
//
//	ctx.JSON(http.StatusOK, res)
//}

type terminalRequest struct {
	Exe string `json:"exe" binding:"required"`
}

func (server *Server) Terminal(ctx *gin.Context) {
	var req terminalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	exeArr := strings.Split(req.Exe, " ")

	if len(exeArr) > 2 {
		err := errors.New("Format command not found. Try running `help`.")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	message := ""
	path := ""
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	switch exeArr[0] {
	case "ls":
		userDirs, err := server.querier.GetUserDirs(ctx, authPayload.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		userDirsStr := "|"

		for _, dir := range userDirs {
			userDirsStr += dir.Name + "|"
		}

		//var listCmd []string
		dirs, err := ioutil.ReadDir(server.config.BinPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		//var strCommand string
		var cmd, dirStr string
		for _, dir := range dirs {
			if dir.IsDir() && strings.Index(userDirsStr, "|"+dir.Name()+"|") > -1 {
				files, _ := ioutil.ReadDir(server.config.BinPath + "/" + dir.Name())
				dirStr = dir.Name()
				for _, file := range files {
					if !file.IsDir() && strings.Index(file.Name(), ".") < 0 {
						cmd = file.Name()
						if cmd != "" && dirStr != "" {
							//listCmd = append(listCmd, dirStr + "/" + cmd + "/main.go")
							//listCmd = append(listCmd, dirStr + "/" + cmd + "/go.mod")
							message += dirStr + "/" + cmd + "/main.go\n\r"
							message += dirStr + "/" + cmd + "/go.mod\n\r\n\r"
						}
					}
				}
			}
		}
	case "build":
		dirCmdArr := strings.Split(exeArr[1], "/")
		if len(dirCmdArr) > 2 {
			err := errors.New("Format command not found. Try running `help`.")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		dir := dirCmdArr[0]
		cmd := dirCmdArr[1]

		if !server.isUserHasDirectory(ctx, authPayload.UserID, dir) {
			err := errors.New("Directory doesn't belong to the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		var args []string
		mainFile := server.config.BinPath + "/" + dir + "/" + cmd + ".src/main.go"
		dirRunner := server.config.BinPath + "/" + dir
		//dirSrc := server.config.BinPath + "/" + dir + "/" + cmd + ".src"

		commandStr := "/usr/local/go/bin/go"
		args = append(args, "build", "-o", cmd, mainFile)
		exeCmd := exec.Command(commandStr, args[0:]...)
		exeCmd.Dir = dirRunner
		stdout, err := exeCmd.Output()
		if err != nil {
			message = err.Error()
		} else {
			if len(stdout) > 0 {
				message = string(stdout)
			} else {
				message = "Succes to build command"
			}
		}
	case "run":
		dirCmdArr := strings.Split(exeArr[1], "/")
		if len(dirCmdArr) > 2 {
			err := errors.New("Format command not found. Try running `help`.")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		dir := dirCmdArr[0]
		cmd := dirCmdArr[1]

		if !server.isUserHasDirectory(ctx, authPayload.UserID, dir) {
			err := errors.New("Directory doesn't belong to the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		var args []string
		dirRunner := server.config.BinPath + "/" + dir

		commandStr := dirRunner + "/" + cmd
		exeCmd := exec.Command(commandStr, args[0:]...)
		exeCmd.Dir = dirRunner
		stdout, err := exeCmd.Output()
		if err != nil {
			message = err.Error()
		} else {
			if len(stdout) > 0 {
				message = string(stdout)
			} else {
				message = "Succes to execute command"
			}
		}
	case "open":
		dirCmdArr := strings.Split(exeArr[1], "/")
		if len(dirCmdArr) > 3 {
			err := errors.New("Format command not found. Try running `help`.")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		dir := dirCmdArr[0]
		cmd := dirCmdArr[1]
		file := dirCmdArr[2]

		if !server.isUserHasDirectory(ctx, authPayload.UserID, dir) {
			err := errors.New("Directory doesn't belong to the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		//path = "<a href=\"editcode/" + dir + "/" + cmd + "\" >Edit command" + dir + "/" + cmd + " </a>"
		path = "/openfile/" + dir + "/" + cmd + "/" + file

	default:
		err := errors.New("Format command not found. Try running `help`.")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	res := commandResponse{
		Message: message,
		Path:    path,
	}

	ctx.JSON(http.StatusOK, res)
	return
}
