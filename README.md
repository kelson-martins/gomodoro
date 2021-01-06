# Gomodoro

`[Go]modoro` is a CLI application for those interested in the Pomodoro technique. The application can be used to run Pomodoros from customized lenghts (defaults to 25m), allowing also to assign a `category` and `subcategory` to each pomodoro session.

Characteristics:
* Uses `SQLite` for storing local pomodoro sessions.
* Support backup of local database via Google Drive integration.
* Support synchronization via `Cassandra` database (powered by DataStax).
* Configuration file is stored at `~/gomodoro/config.yaml`
* Current version supports Linux and Darwin

## Todo
* Windows support.
* Support for Gomodoros from remote Cassandra database.

## Help

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
  sync        A brief description of your command
  totals      Query information about [Go]modoro usage totals
  version     A brief description of your command
```

## Installation

```
sudo make -E install
```

## Usage

#### Running a Gomodoro

```
gomodoro run --catergory coding
Gomodoro Started: 06-01-2021 21:13:45
▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆                                   %33.3
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
    username: ""
    keyspace: ""
    password: ""
```

3.  You are ready to sync with `gomodoro sync`
```
gomodoro sync
2021/01/06 20:57:03 [INFO] cassandra token retrieved successfuly
2021/01/06 20:57:03 [INFO] cassandra table synced successfully. DSE response: {"success":true}
2021/01/06 20:57:16 [INFO] local [Go]modoros pushed to remote successfully
2021/01/06 20:57:16 [INFO] [Go]modoro already synchronized, no remote changes were pulled.

```