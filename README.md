# EIR
This repository offers a EIR NF according to the ETSI TS 129 511 V17.3.0's 3GPP spec

This implementation is intended for review by the Free5GC team.

As Free5GC does not natively support the EIR NF, this repository should be tested in conjunction with the following pull requests:
- [!4](https://github.com/adjivas/amf/pull/4) adds to the AMF the N17 EIR verification step
- [!1](https://github.com/adjivas/nrf/pull/3) adapts nRF to provide a EIR registration/deregistration notifications

To populate the database with EIR data, we can use these commands:
```shell
% mongosh
> use free5gc
> db.policyData.ues.eirData.insertOne( { "pei": "imeisv-4370816125816151", "equipement_status": "WHITELISTED" })
```

To run and test this NF, use the following commands:
```shell
% go run cmd/main.go --config config/eircfg.yaml
% http://127.0.0.8:8000/n5g-eir-eic/v1/equipement-status?pei=imeisv-4370816125816151
```

To run linters, use:
```shell
golangci-lint run
```

We can consider to add a EquipmentStatus enum on the front of the Free5GC [Webconsole](https://github.com/free5gc/webconsole)

Like every Free5gc NF service, this EIR NF is executable with the `go run cmd/main.go` command
It's will use a default configuration path as `config/eircfg.yaml`, the [eircfg.yaml](https://github.com/adjivas/eir/blob/main/config/eircfg.yaml) configuration file was added as an example.

The EIR configuration file supports a optional `configuration.defaultStatus` to set the default EquipmentStatus when it's wasn't provided on the database.

This work is sponsored by [Free Mobile](https://mobile.free.fr)!
