# DnsMadeEasy DNS Update

[![Build Status](https://secure.travis-ci.org/gblmarquez/dnsmadeeasy.png)](https://travis-ci.org/gblmarquez/dnsmadeeasy)

It works as a service on Windows, Mac, and Linux!

This is a simple tool that periodically updates your IP when it changes for [DnsMadeEasy](http://www.dnsmadeeasy.com/)'s dynamic DNS.

## Download

Download the pre-built binary from Gobuild.io:

[![Gobuild Download](http://gobuild.io/badge/github.com/gblmarquez/dnsmadeeasy/download.png)](http://gobuild.io/github.com/gblmarquez/dnsmadeeasy)

## Commands

### Install

To install the service you run the follow command, that will install and create `dnsmadeeasy.cfg` config file

`dnsmadeeasy install [USER] [PASSWORD] [RECORD_ID]`

Arguments description:
- `USER` from your DnsMadeEasy credentials 
- `PASSWORD` from DnsMadeEasy credential or configured per record, so you don't need to use your account password. 
- `RECORD_ID` identifies the record to update.

### Uninstall

To un-install the service you run the follow command, that will remove the service from registry

`dnsmadeeasy remove`

### Self Run

To run the application without installing you need to use the follow command

`dnsmadeeasy run [USER] [PASSWORD] [RECORD_ID]`

Arguments details:
- `USER` from your DnsMadeEasy credentials 
- `PASSWORD` from DnsMadeEasy credential or configured per record, so you don't need to use your account password. 
- `RECORD_ID` identifies the record to update.

### Control 

To start the service use the follow command

`dnsmadeeasy start`

To stop the service use the follow command

`dnsmadeeasy stop`

## Troubleshooting

### Linux and Mac
The service uses the syslogd and the log entries are created on file `/var/log/system.log`.

### Windows
The service uses the system events and the log entries can be viewed using event viewer.

## License

The tool is released as OSS under the [New BSD license](https://github.com/gblmarquez/dnsmadeeasy/blob/master/LICENSE).
