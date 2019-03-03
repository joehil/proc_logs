# proc_logs - a Linux log message processor
This is a Go programm to read UNIX logs and to react to log messages. The software can merge any number of logs and thus react to any event that is written to any of those logs. Examples of logs are syslog, mail.log, fail2ban.log, messages or apache error.log.

It is compiled by these statements:
go build --ldflags="-s -w"

The compiled module is proc_logs.

To get a compressed executable module you can use the command:
upx -9 proc_logs

The software consists of three parts:

1. proc_logs.go
This is the part that reads in the specified logs and is the driver of the processing rules.

2. proc_util.go
There are utilities in this file that can be used to process the rules.

3. proc_rules.go
This is the file where the processing rules are defined. This is the file that is to be customized do make the software do what you want it to do.

To be continued.
