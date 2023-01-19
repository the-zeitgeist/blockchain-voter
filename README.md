# BLOCKCHAIN VOTER

This API intends to initialize a voting system through distributed nodes and to allow the users to vote.

## BUILDING

In order to build a binary run:

```bash
go build .
```

## RUNNING

After your run the project, wheter using the binary or locally, it will check for the `candidates.json` file, which contains the voting information and which will be used to distribute the voting system through your nodes, this file should not be modified manually (a future version will verify the integrity of the file).

In case you don't have this file yet, you'll have two options:

1. Create a new system
2. Joining as a voting node to a previous system, to do so, you'll need to provide an url or a public ip to join the network.

Once the system is initialized, you'll be able to interact with the voting system.

### Sending a transaction

You can add a vote by sending a transaction to the network.

Method: `POST`
Endpoint: `/transaction`
Body:

```json
{
    "voter": "aasda1221a3121a4wa2",
    "candidate": "01"
}
```

**candidate must be a valid candidate id within the candidates.json file**
**voter musth be encrypted**

### Get all the transactions

You can see all the transactions within the network

Method: `GET`
Endpoint: `/transactions`

### Get the results

You can see the sumarized results

Method: `GET`
Endpoint: `/chains/result`

### Validating the transactions

Anyone can check the integrity of the data.

Method: `GET`
Endpoint: `/chains/valid`

Author: Nicolas Macias
