@echo off
chcp 65001 >nul
cls
%~dp0protoc-35.1-win64\bin\protoc ^
--plugin=protoc-gen-go=%GOPATH%\bin\protoc-gen-go.exe ^
--plugin=protoc-gen-go-grpc=%GOPATH%\bin\protoc-gen-go-grpc.exe ^
--proto_path=%~dp0proto ^
--go_out=%~dp0kubeapi ^
--go-grpc_out=%~dp0kubeapi ^
%~dp0proto\*.proto