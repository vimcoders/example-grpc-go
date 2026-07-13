@echo off
%~dp0protoc-35.1-win64\bin\protoc ^
--plugin=protoc-gen-go=%~dp0..\bin\protoc-gen-go.exe ^
--plugin=protoc-gen-go-grpc=%~dp0..\bin\protoc-gen-go-grpc.exe ^
--proto_path=%~dp0 ^
--go_out=%~dp0kubeapi ^
--go-grpc_out=%~dp0kubeapi ^
%~dp0proto\*.proto