# Comparison Filter

|              |                                                                  |
| ------------ | ---------------------------------------------------------------- |
| name         | Serial Data Decryptor                                            |
| version      | v1.0.0                                                           |
| GitHub       | [serial-data-decryptor](https://github.com/weeve-modules/serial-data-decryptor) |
| DockerHub    | [weevenetwork/serial-data-decryptor](https://hub.docker.com/r/weevenetwork/serial-data-decryptor)     |
| authors      | Shrikant Bhusalwad, Paul Gaiduk, Vadim Vygonets                  |

***
## Table of Content

- [Comparison Filter](#comparison-filter)
  - [Table of Content](#table-of-content)
  - [Description](#description)
  - [Module Variables](#module-variables)
  - [Module Testing](#module-testing)
  - [Dependencies](#dependencies)
  - [Input](#input)
  - [Output](#output)
***

## Description

Weeve module to receive serial data in JSON format by Ingress module and decrypt the value, then send the decrypted data to next moudle.

## Module Variables

The following module configurations can be provided in a data service designer section on weeve platform:

| Environment Variables | type   | Description                                       |
| --------------------- | ------ | ------------------------------------------------- |
| MODULE_NAME           | string | Name of the module                                |
| MODULE_TYPE           | string | Type of the module (Input, Processing, Output)    |
| LOG_LEVEL             | string | Allowed log levels: DEBUG, INFO, WARNING, ERROR, CRITICAL. Refer to `logging` package documentation. |
| INGRESS_HOST          | string | Host to which data will be received               |
| INGRESS_PORT          | string | Port to which data will be received               |
| EGRESS_URLS           | string | HTTP ReST endpoint for the next module            |
| AES_KEY               | string | AES key to decrypt cyphertext in base64 format    |

## Module Testing

TBD

## Dependencies

The following are module dependencies:

* github.com/sirupsen/logrus
* github.com/go-playground/validator/v10

## Input

Input to this module is JSON body single object:

```json
{
    "iv": "jc/GHiyMZmDkj1FK",
    "cyphertext": "TlMGbx5LRKPAFHRcOVKvy9veapflEdXJo48PQ27u95HdcchyqgeQSzLFetmcT2EjswXITGIAjcVUVntIHPNGL8ZsIzGbdik3kdilZtq8ADyZsQ=="
}
```

## Output

Output of this module is the original plaintext. In this case this is the decrypted plaintext of the cyphertext from above:
```text
cyccnt 0059d6ac (+0007eb06), ocnt 0002ffcd, ent +038d, oh 378782µ
```
