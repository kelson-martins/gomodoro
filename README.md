# Gomodoro

`[Go]modoro` is a CLI application for those interested in the Pomodoro technique. The application can be used to run Pomodoros from customized lengths (defaults to 25m), allowing also to assign a `category` and `subcategory` to each pomodoro session.

Characteristics:
* Uses `SQLite` for storing the local Pomodoro sessions.
* Support backup of local database via Google Drive integration.
* Support synchronization via `Cassandra` database (powered by DataStax).
* Configuration file `config.yaml` and SQLite database `gomodoro.db` are stored at `~/gomodoro`.
* Current version supports Linux and Darwin.

## Todo
* Windows support.
* Support for deletion of sessions from remote Cassandra database.
* Change the visual alert for the Pomodoro conclusion to another graphical library.

## Installation

1. If you have a `go` development environment setup, run:

```
sudo make install
```

The `gomodoro` binary will be added automatically to the `/usr/local/bin/` location, so you can run `gomodoro` command straightaway.

2. If you do not have a development environment:

a. Download the binary related to your architecture from the [releases page](https://github.com/kelson-martins/gomodoro/releases/tag) and move it into any desired location.

b. From the root of the repository, run `sudo ./install.sh` to initialize the app folder structure.

c. You are ready to run the `gomodoro` command.

## Usage

```
gomodoro -h
Gomodoro - The ultimate productivity tool. Keep [Go]ing

Usage:
  gomodoro [command]

Available Commands:
  backup      Backup the [Go]modoro database into Google Drive
  delete      Delete a [Go]modoro from the database
  help        Help about any command
  list        List [Go]modoro sessions
  run         Run a [Go]modoro
  sync        Synchronize [Go]modoros from local/remote
  totals      Query information about [Go]modoro usage totals
  version     Display [Go]modoro version
```

#### Running a Gomodoro

```
gomodoro run --category coding
Gomodoro Started: 06-01-2021 21:13:45
▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆                                   %33.3
```

#### Totals

```
gomodoro totals
[Go]modoros today:  2
[Go]modoros yesterday:  3
[Go]modoros this month:  18
[Go]modoros last month:  99
[Go]modoros all-time:  207
```

#### Listing Gomodoros

```
gomodoro list
ID: 207  2021-01-06 14:53:34     Category: coding        SubCategory: gomodoro feature x
ID: 206  2021-01-06 10:56:26     Category: coding
ID: 203  2021-01-06 08:33:04     Category: coding
ID: 202  2021-01-05 15:27:15     Category: coding        SubCategory: gomodoro feature y
ID: 201  2021-01-05 14:41:28     Category: coding
ID: 200  2021-01-05 11:37:45     Category: coding
```

## Remote Synchronization (DataStax Cassandra)

The Gomodoro application supports remote synchronization with the `gomodoro sync` command. The sync functionality allows you to keep your Gomodoro progress across different workstations. The functionality will:

1. Push all local finished **gomodoros** into your remote Cassandra database
2. Pull all remote gomodoros into local

#### Synchronization Requirements

1. Create your Cassandra database (powered by DataStax) [here](https://astra.datastax.com/register). The free tier supports up to 5GB of data, which is enough as the Gomodoro sync functionality consumes small bandwidth footprint.

2. Update the Gomodoro app `config.yaml` located at `~/gomodoro/config.yaml` with the information from your Cassandra database, more specifically the following fields:

```
cassandra:
    cluster_id: ""
    cluster_region: ""
    keyspace: ""
    token: ""
```

3.  You are ready to sync with `gomodoro sync`
```
gomodoro sync
2021/01/06 20:57:03 [INFO] cassandra token retrieved successfuly
2021/01/06 20:57:03 [INFO] cassandra table synced successfully. DSE response: {"success":true}
2021/01/06 20:57:16 [INFO] local [Go]modoros pushed to remote successfully
2021/01/06 20:57:16 [INFO] [Go]modoro already synchronized, no remote changes were pulled.

```

## Backup

Backup of the local SQLite database is supported via Google Drive integration.

At the first backup execution, a link for integration will be displayed on the console.
```
gomodoro backup
[INFO] initiating Gomodoro database backup
[INFO] backup folder gomodoro_backup found, re-using it.
[INFO] [Go]modoro database 'gomodoro.db' successfully saved at 'gomodoro_backup' directory
```

The backup will be stored at the root of your Google Drive filesystem at the `/gomodoro_backup` location.
