@echo off
cd Common
echo Current directory: %cd%
::生成消息结构、表结构
go run GenModule\genStructs\main.go
::表转JSON 一定要在表结构生成完之后才能执行
go run GenModule\genOthers\main.go
pause