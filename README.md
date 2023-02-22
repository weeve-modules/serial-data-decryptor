# AES-GCM Decryptor

|           |                                                                                                   |
| --------- | ------------------------------------------------------------------------------------------------- |
| name      | AES-GCM Decryptor                                                                             |
| version   | v2.0.0                                                                                            |
| GitHub    | [serial-data-decryptor](https://github.com/weeve-modules/serial-data-decryptor)                   |
| DockerHub | [weevenetwork/serial-data-decryptor](https://hub.docker.com/r/weevenetwork/serial-data-decryptor) |
| authors   | Shrikant Bhusalwad, Paul Gaiduk, Vadim Vygonets                                                   |

***
## Table of Content

- [AES-GCM Decryptor](#aes-gcm-decryptor)
  - [Table of Content](#table-of-content)
  - [Description](#description)
  - [Module Variables](#module-variables)
  - [Module Testing](#module-testing)
  - [Dependencies](#dependencies)
  - [Input](#input)
  - [Output](#output)
***

## Description

Weeve module that decrypts the cyphertext received from the previous module using AES-GCM algorithm and forwards the plaintext to the next module. For incoming and outgoing message format see [Input](#input) and [Output](#output) sections.

## Module Variables

Weave module to receive serial data in JSON format from Ingress module and decrypt the cyphertext, then send the decrypted data to the next module.

| Environment Variables | type   | Description                                                                                          |
| --------------------- | ------ | ---------------------------------------------------------------------------------------------------- |
| MODULE_NAME           | string | Name of the module                                                                                   |
| MODULE_TYPE           | string | Type of the module (Input, Processing, Output)                                                       |
| LOG_LEVEL             | string | Allowed log levels: DEBUG, INFO, WARNING, ERROR, CRITICAL. Refer to `logging` package documentation. |
| INGRESS_HOST          | string | Host to which data will be received                                                                  |
| INGRESS_PORT          | string | Port to which data will be received                                                                  |
| EGRESS_URLS           | string | HTTP ReST endpoint for the next module                                                               |
| AES_KEY               | string | AES key to decrypt cyphertext in base64 format                                                       |

## Module Testing

TBD

## Dependencies

The following are module dependencies:

* github.com/sirupsen/logrus
* github.com/go-playground/validator/v10

## Input

Input to this module is JSON body single object with fields `iv` and `cyphertext`. The `iv` field is the initialization vector in base64 format and the `cyphertext` field is the cyphertext in base64 format. For example:

```json
{
    "iv": "jc/GHiyMZmDkj1FK",
    "cyphertext": "TlMGbx5LRKPAFHRcOVKvy9veapflEdXJo48PQ27u95HdcchyqgeQSzLFetmcT2EjswXITGIAjcVUVntIHPNGL8ZsIzGbdik3kdilZtq8ADyZsQ=="
}
```

## Output

Output of this module is a JSON object with the original plaintext. In this case this is the decrypted plaintext of the cyphertext from above:
```json
{
    "plaintext": "cyccnt 0059d6ac (+0007eb06), ocnt 0002ffcd, ent +038d, oh 378782µ"
}
```
