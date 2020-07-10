# Synse SNMP Plugin Base

A common base plugin definition for SNMP-type Synse plugins.

> NOTE: This project is currently under development and is not considered stable or ready for general use.

## Getting Started

This plugin is not intended for direct distribution/usage. Instead, it should be used as a package
to enable building SNMP-based plugins. As an example, see the [snmp-ups-plugin](https://github.com/vapor-ware/synse-snmp-ups-plugin).

When writing a plugin, you get get the SNMP plugin base with

```
go get github.com/vapor-ware/synse-snmp-base
```

## SNMP Plugin Configuration

Plugin and device configuration are described in detail in the [Synse SDK Documentation][sdk-docs].

This base plugin defines the configuration options which all plugins that use it will inherit.
Plugins will need to define their own plugin configuration (e.g. `config.yaml`) with the `dynamicRegistration`
block defined. This is how the SNMP plugin knows which servers to communicate with.

The options expected in the dynamic registration block are defined below. Note that each plugin
which uses this base is effectively just codifying one or more MIBs that this base can load and use.
As such, it is expected that the full MIB be represented, as no MIB walk is performed on startup.

An example configuration:

```yaml
dynamicRegistration:
  config:
  - mib: UPS-MIB
    version: v3
    agent: 'udp://127.0.0.1:1024'
    security:
      level: authPriv
      context: public
      username: simulator
      authentication:
        protocol: SHA
        passphrase: foobar
      privacy:
        protocol: AES
        passphrase: foobar
```

### Dynamic Registration Options

Below are the fields that are expected in each of the dynamic registration items.
If no default is specified (`-`), the field is required.

| Field                              | Description | Default |
| ---------------------------------- | ----------- | ------- |
| mib                                | The name of the MIB to use for the configured agent. The MIB name(s) are defined by the plugin implementation using the SNMP base. | `-` |
| version                            | The SNMP protocol version. **Note**: The security parameters, below, are only valid for SNMP `v3`. (Valid values include: `v1`, `v2`, `v2c`, `v3`) | `-` |
| agent                              | The address of the SNMP server to connect to. If this does not contain a protocol prefix, `udp://` is used by default. Only `udp` and `tcp` are supported protocols. If no port is specified, `161` is used by default. | `-` |
| community                          | The SNMP community string. | `""` |
| timeout                            | The timeout to use for SNMP requests. | `3s` |
| retries                            | The number of times to retry a request within the timeout period. | `1` |
| security.level                     | (`v3` only) The security message flag. Valid values include (case insensitive): `noauthnopriv`, `authnopriv`, `authpriv`, `reportable`| `-` |
| security.context                   | (`v3` only) The SNMPv3 context name. | `""` |
| security.username                  | (`v3` only) The SNMPv3 user name. | `""` |
| security.authentication.protocol   | (`v3` only) The SNMPv3 authentication protocol. Supported values include (case insensitive): `md5`, `sha`, `none`. | `-` |
| security.authentication.passphrase | (`v3` only) The passphrase for authentication. | `""` |
| security.privacy.protocol          | (`v3` only) The SNMPv3 privacy protocol. Supported values include (case insensitive): `aes`, `des`, `none`.| `-` |
| security.privacy.passphrase        | (`v3` only) The passphrase for privacy. | `""` |

### Reading Outputs

Outputs are referenced by name. A single device may have more than one instance
of an output type. The base plugin does **not** defined any custom outputs. All
outputs should be configured via device configuration and should reference either
the [built-in outputs](https://synse.readthedocs.io/en/latest/sdk/concepts/reading_outputs/#built-ins)
or outputs defined by a plugin built on top of this base.

### Device Handlers

Device Handlers are referenced by name.

| Name       | Description                                    | Outputs              | Read  | Write | Bulk Read | Listen |
| ---------- | ---------------------------------------------- | -------------------- | :---: | :---: | :-------: | :----: |
| read-only  | A handler only supporting OID reads.           | any (device defined) | ✓     | ✗     | ✗         | ✗      |

### Write Values

The SNMP base plugin does not currently support writing.

## Compatibility

Below is a table describing the compatibility of the plugin base versions with Synse platform versions.

|             | Synse v2 | Synse v3 |
| ----------- | -------- | -------- |
| plugin v0.x | ✗        | ✓        |

## Contributing / Reporting

If you experience a bug, would like to ask a question, or request a feature, open a
[new issue](https://github.com/vapor-ware/synse-snmp-base/issues) and provide as much
context as possible. All contributions, questions, and feedback are welcomed and appreciated.

## License

The Synse SNMP plugin base is licensed under GPLv3. See [LICENSE](LICENSE) for more info.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-base.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-base?ref=badge_large)
