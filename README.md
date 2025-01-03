# What is this?

dkb2ynab fetches transaction data from DKB (using [dkbrobot](https://github.com/pczora/dkbrobot)) and converts it to a 
[YNAB](https://youneedabudget.com)-importable format.

# Who is this for?

Currently: **ME**. It might not work for you, but if you find a bug and create an issue I will try to fix things. But 
this is first and foremost something I use, and I am happy as long as it works for me.

# How can I use this?

First, start by cloning this repo. Also make sure you have at least Go 1.18 installed. There are no downloadable 
releases yet so you have to build the project on your machine.

## Create a config file

It should be called `config.yaml` and be placed in the project root. 
Replace `$YOUR_DKB_USERNAME` with... your DKB
username. Currently, (and for the foreseeable future) only DKB is supported (hence the name) so, `bank` has to be 
`"DKB"`.

### Interactive password input

```yaml
banks:
  - name: "DKB"
    bank: "DKB"
    credentials:
      username: "$YOUR_DKB_USERNAME"
```

### Password in keyring

If you don't want to input your DKB password each time you run the app, you can add it to your system's keyring.
In this case, use the following config. To add the password to the keyring on macOS, run 
`security add-generic-password -a $YOUR_DKB_USERNAME -s dkb2ynab_dkb -w` in a terminal.

```yaml
banks:
  - name: "DKB"
    bank: "DKB"
    credentials:
      username: "$YOUR_DKB_USERNAME"
      password:
        fromKeyring:
          key: "DKB"
```

## Run the app

```shell
go run cmd/main.go
```
