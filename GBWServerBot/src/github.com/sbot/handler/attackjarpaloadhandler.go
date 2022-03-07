package handler

import (
	"context"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/utils/fileutils"
	"os"

	"os/exec"
	"path/filepath"
	"time"
)

const (

	StorePath = "/var/tmp/"
	JarMainName = "JarMain"
	JarMainJavaName = "JarMain.java"
	JarMainJavaClassName = "JarMain.class"
	JarMetaName = "MANIFEST.MF"
	JarMetaDir = "META-INF"
)

type AttackJarPayloadHandle struct {

	jarTemplateStorePath string

	javaVersion string

}

type TemplateJarMainData struct {

	Cmd string
}

type TemplateMetaFileData struct {

	Version string
	Main string
}


func NewAttackJarPayloadHandle(jarTemplateStorePath, javaVersion string) *AttackJarPayloadHandle {

	return &AttackJarPayloadHandle{
		jarTemplateStorePath: jarTemplateStorePath,
		javaVersion:          javaVersion,
	}
}

func (ajh *AttackJarPayloadHandle) compilerJarMain(cmd string) (string,error){

	javaMainTemPath := filepath.Join(ajh.jarTemplateStorePath,fmt.Sprintf("%s.tpl",JarMainJavaName))

	if !fileutils.FileIsExisted(javaMainTemPath) {
		errS := fmt.Sprintf("Jar Java Main source template file not found in path:%s",javaMainTemPath)
		log.Error(errS)
		return "",fmt.Errorf(errS)
	}

	tempData := &TemplateJarMainData{Cmd:cmd}

	javaFile := filepath.Join(StorePath,JarMainJavaName)

	if err := fileutils.GenerateFileFromTemplateFile(javaFile,javaMainTemPath,tempData);err!=nil||!fileutils.FileIsExisted(javaFile){

		errS := fmt.Sprintf("Generate java main source file from template file:%s,failed:%v",javaMainTemPath,err)
		log.Error(errS)
		return "",fmt.Errorf(errS)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	osCmd := exec.CommandContext(ctx, "javac", javaFile)
	osCmd.CombinedOutput()

	jclassFile := filepath.Join(StorePath,JarMainJavaClassName)

	if !fileutils.FileIsExisted(jclassFile){
		errS := fmt.Sprintf("Compile java main source file:%s is error",javaFile)
		log.Error(errS)
		return "",fmt.Errorf(errS)
	}

	return jclassFile,nil
}

func (ajh *AttackJarPayloadHandle) makeMetaFile()  (string,error) {

	metaFileTemPath := filepath.Join(ajh.jarTemplateStorePath,fmt.Sprintf("%s.tpl",JarMetaName))

	if !fileutils.FileIsExisted(metaFileTemPath) {
		errS := fmt.Sprintf("Jar Java MetaFile template file not found in path:%s",metaFileTemPath)
		log.Error(errS)
		return "",fmt.Errorf(errS)
	}

	tempData := &TemplateMetaFileData{
		Version: ajh.javaVersion,
		Main:   JarMainName ,
	}

	metaFile := filepath.Join(StorePath,JarMetaDir,JarMetaName)
	os.MkdirAll(filepath.Join(StorePath,JarMetaDir),0755)

	if err := fileutils.GenerateFileFromTemplateFile(metaFile,metaFileTemPath,tempData);err!=nil||!fileutils.FileIsExisted(metaFile){

		errS := fmt.Sprintf("Generate java meta file from template file:%s,failed:%v",metaFileTemPath,err)
		log.Error(errS)
		return "",fmt.Errorf(errS)
	}

	return metaFile,nil
}

func (ajh *AttackJarPayloadHandle) Handle(request *model.MakeJarAttackPayloadRequest) (string,error){

	jarFiles := make([]string,0)

	jarFile := filepath.Join(StorePath,JarMainName+".jar")


	classFile,err:= ajh.compilerJarMain(request.Cmd)

	if err!=nil {
		return "",err
	}

	jarFiles = append(jarFiles,classFile)

	metaFile,err := ajh.makeMetaFile()

	if err!=nil {
		return "",err
	}

	jarFiles = append(jarFiles,metaFile)


	if err =fileutils.ZipFiles(jarFile,StorePath,jarFiles,true);err!=nil {

		return "",err
	}

	if !fileutils.FileIsExisted(jarFile) {

		errS := fmt.Sprintf("Cannot make jar file:%s for cmd:%s",jarFile,request.Cmd)
		log.Errorf(errS)
		return "",fmt.Errorf(errS)
	}

	return jarFile,nil
}


