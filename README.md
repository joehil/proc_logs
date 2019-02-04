# proc_logs - a Linux log message processor
This is a Go programm to read UNIX logs and to react to log messages. The software can merge any number of logs and thus react to any event that is written to any of those logs. Examples of logs are syslog, mail.log, fail2ban.log, messages or apache error.log.

It is compiled by these statements:
go build

The compiled module is proc_logs.
