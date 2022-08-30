# serial-data-decryptor
Weeve module to receive serial data in JSON format by Ingress module and decrypt the value, then send the decrypted data to next moudle.

Expected data schema
{
	IV         string `json:"iv"` # Initialization Vector
	CypherText string `json:"cyphertext"` # Cypher Text
}

Sample encrypted data
{
    "iv": "jc/GHiyMZmDkj1FK",
    "cyphertext": "TlMGbx5LRKPAFHRcOVKvy9veapflEdXJo48PQ27u95HdcchyqgeQSzLFetmcT2EjswXITGIAjcVUVntIHPNGL8ZsIzGbdik3kdilZtq8ADyZsQ=="
}

Sample decryped data
cyccnt 0059d6ac (+0007eb06), ocnt 0002ffcd, ent +038d, oh 378782µ