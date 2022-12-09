<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [fx/base/v1/query.proto](#fx/base/v1/query.proto)
    - [GetGasPriceRequest](#fx.base.v1.GetGasPriceRequest)
    - [GetGasPriceResponse](#fx.base.v1.GetGasPriceResponse)
  
    - [Query](#fx.base.v1.Query)
  
- [fx/ibc/applications/transfer/v1/query.proto](#fx/ibc/applications/transfer/v1/query.proto)
    - [Query](#fx.ibc.applications.transfer.v1.Query)
  
- [fx/ibc/applications/transfer/v1/transfer.proto](#fx/ibc/applications/transfer/v1/transfer.proto)
    - [FungibleTokenPacketData](#fx.ibc.applications.transfer.v1.FungibleTokenPacketData)
  
- [fx/ibc/applications/transfer/v1/tx.proto](#fx/ibc/applications/transfer/v1/tx.proto)
    - [MsgTransfer](#fx.ibc.applications.transfer.v1.MsgTransfer)
  
    - [Msg](#fx.ibc.applications.transfer.v1.Msg)
  
- [fx/legacy/other/query.proto](#fx/legacy/other/query.proto)
    - [GasPriceRequest](#fx.other.GasPriceRequest)
    - [GasPriceResponse](#fx.other.GasPriceResponse)
  
    - [Query](#fx.other.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="fx/base/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/base/v1/query.proto



<a name="fx.base.v1.GetGasPriceRequest"></a>

### GetGasPriceRequest







<a name="fx.base.v1.GetGasPriceResponse"></a>

### GetGasPriceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_prices` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.base.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GetGasPrice` | [GetGasPriceRequest](#fx.base.v1.GetGasPriceRequest) | [GetGasPriceResponse](#fx.base.v1.GetGasPriceResponse) |  | GET|/fx/base/v1/gas_price|

 <!-- end services -->



<a name="fx/ibc/applications/transfer/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/ibc/applications/transfer/v1/query.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ibc.applications.transfer.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.
Deprecated: This service is deprecated. It may be removed in the next
version. Replace ibc.applications.transfer.v1.Query

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DenomTrace` | [.ibc.applications.transfer.v1.QueryDenomTraceRequest](#ibc.applications.transfer.v1.QueryDenomTraceRequest) | [.ibc.applications.transfer.v1.QueryDenomTraceResponse](#ibc.applications.transfer.v1.QueryDenomTraceResponse) | DenomTrace queries a denomination trace information. | |
| `DenomTraces` | [.ibc.applications.transfer.v1.QueryDenomTracesRequest](#ibc.applications.transfer.v1.QueryDenomTracesRequest) | [.ibc.applications.transfer.v1.QueryDenomTracesResponse](#ibc.applications.transfer.v1.QueryDenomTracesResponse) | DenomTraces queries all denomination traces. | |
| `Params` | [.ibc.applications.transfer.v1.QueryParamsRequest](#ibc.applications.transfer.v1.QueryParamsRequest) | [.ibc.applications.transfer.v1.QueryParamsResponse](#ibc.applications.transfer.v1.QueryParamsResponse) | Params queries all parameters of the ibc-transfer module. | |
| `DenomHash` | [.ibc.applications.transfer.v1.QueryDenomHashRequest](#ibc.applications.transfer.v1.QueryDenomHashRequest) | [.ibc.applications.transfer.v1.QueryDenomHashResponse](#ibc.applications.transfer.v1.QueryDenomHashResponse) | DenomHash queries a denomination hash information. | |
| `EscrowAddress` | [.ibc.applications.transfer.v1.QueryEscrowAddressRequest](#ibc.applications.transfer.v1.QueryEscrowAddressRequest) | [.ibc.applications.transfer.v1.QueryEscrowAddressResponse](#ibc.applications.transfer.v1.QueryEscrowAddressResponse) | EscrowAddress returns the escrow address for a particular port and channel id. | |

 <!-- end services -->



<a name="fx/ibc/applications/transfer/v1/transfer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/ibc/applications/transfer/v1/transfer.proto



<a name="fx.ibc.applications.transfer.v1.FungibleTokenPacketData"></a>

### FungibleTokenPacketData
FungibleTokenPacketData defines a struct for the packet payload
See FungibleTokenPacketData spec:
https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | the token denomination to be transferred |
| `amount` | [string](#string) |  | the token amount to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |
| `router` | [string](#string) |  | the router is hook destination chain |
| `fee` | [string](#string) |  | the fee is destination fee |
| `memo` | [string](#string) |  | optional memo |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/ibc/applications/transfer/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/ibc/applications/transfer/v1/tx.proto



<a name="fx.ibc.applications.transfer.v1.MsgTransfer"></a>

### MsgTransfer
MsgTransfer defines a msg to transfer fungible tokens (i.e Coins) between
ICS20 enabled chains. See ICS Spec here:
https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source_port` | [string](#string) |  | the port on which the packet will be sent |
| `source_channel` | [string](#string) |  | the channel by which the packet will be sent |
| `token` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | the tokens to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |
| `timeout_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Timeout height relative to the current block height. The timeout is disabled when set to 0. |
| `timeout_timestamp` | [uint64](#uint64) |  | Timeout timestamp (in nanoseconds) relative to the current block timestamp. The timeout is disabled when set to 0. |
| `router` | [string](#string) |  | the router is hook destination chain |
| `fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | the tokens to be destination fee |
| `memo` | [string](#string) |  | optional memo |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ibc.applications.transfer.v1.Msg"></a>

### Msg
Msg defines the ibc/transfer Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Transfer` | [MsgTransfer](#fx.ibc.applications.transfer.v1.MsgTransfer) | [.ibc.applications.transfer.v1.MsgTransferResponse](#ibc.applications.transfer.v1.MsgTransferResponse) | Transfer defines a rpc handler method for MsgTransfer. | |

 <!-- end services -->



<a name="fx/legacy/other/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/legacy/other/query.proto



<a name="fx.other.GasPriceRequest"></a>

### GasPriceRequest







<a name="fx.other.GasPriceResponse"></a>

### GasPriceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_prices` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.other.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GasPrice` | [GasPriceRequest](#fx.other.GasPriceRequest) | [GasPriceResponse](#fx.other.GasPriceResponse) |  | GET|/other/v1/gas_price|

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
