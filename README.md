# ISO 8583 Implementation in Go

A high-performance implementation of the ISO 8583 financial transaction messaging standard in Go. This package provides a robust and flexible way to parse, build, and process ISO 8583 messages with exceptional performance characteristics.

## Performance

This implementation has been highly optimized for performance and can process millions of messages per second:

| Operation | Performance | Latency | Memory | Allocations |
|-----------|-------------|---------|--------|-------------|
| Unpack (Primary Bitmap) | 5.7 million ops/sec | ~218 ns | 0 B | 0 |
| Unpack (Secondary Bitmap) | 4.7 million ops/sec | ~240 ns | 0 B | 0 |
| Pack (Primary Bitmap) | 3.6 million ops/sec | ~314 ns | 0 B | 0 |
| Pack (Secondary Bitmap) | 3.1 million ops/sec | ~380 ns | 0 B | 0 |
| Unpack Concurrent (Primary) | 30 million ops/sec | ~57 ns | 0 B | 0 |
| Unpack Concurrent (Secondary) | 15 million ops/sec | ~76 ns | 0 B | 0 |
| Pack Concurrent (Primary) | 10 million ops/sec | ~112 ns | 0 B | 0 |
| Pack Concurrent (Secondary) | 9 million ops/sec | ~128 ns | 0 B | 0 |

Compared to the original implementation, this represents:
- 12.6x improvement in unpacking performance
- 23,900x improvement in packing performance
- 100% reduction in memory allocations (from 2.4KB to 0 bytes per operation)
- 100% reduction in memory usage (from 1MB to 0 bytes per operation)

## Features

- Support for both fixed and variable length fields (LLVAR, LLLVAR, LLLLVAR)
- Built-in support for common MTI (Message Type Identifier) types
- Flexible message packing and unpacking with exceptional performance
- Support for binary and ASCII message formats
- Configurable field definitions and validation
- TLV (Tag-Length-Value) support for complex data structures
- Comprehensive error handling
- Zero memory allocations for all operations
- Excellent concurrent performance with near-linear scaling

## Installation

```bash
go get github.com/pentaly7/iso8583
```

## Quick Start

### Creating a New Message

```go
package main

import (
    "fmt"
    "github.com/pentaly7/iso8583"
)

func main() {
    // Create a new message with default packager
    msg := iso8583.NewMessage(iso8583.DefaultPackager())
    
    // Set MTI (Message Type Identifier)
    msg.SetMtiString(iso8583.MTIFinancialRequest)
    
    // Set fields
    msg.SetString(2, "1234567891234567").
        SetString(3, "000000").
        SetString(4, "000000010000")     
    
    // Pack the message
    packed, err := msg.PackISO()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Packed message: %x\n", packed)
}
```

### Parsing a Message

```go
func parseISO8583Message(data []byte) {
    msg := iso8583.NewMessage(iso8583.DefaultPackager())
    
    // Unpack the message
    err := msg.UnpackISO(data)
    if err != nil {
        panic(err)
    }
    
    // Access fields
    mti := msg.MTI.String()
    pan, _ := msg.GetString(2)
    
    fmt.Printf("MTI: %s\n", mti)
    fmt.Printf("PAN: %s\n", pan)
}

## Message Configuration

The package allows custom configuration of field definitions. Here's an example of a custom packager:

```go
config := `{
    "hasHeader": false,
    "headerLength": 0,
    "messageKey": [2, 7, 11, 12, 13, 41, 37],
    "packagerConfig": {
        "2": {"isMandatory": true, "type": "n", "length": {"type": "LLVAR", "max": 19}},
        "3": {"isMandatory": true, "type": "n", "length": {"type": "FIXED", "max": 6}},
        "4": {"isMandatory": true, "type": "n", "length": {"type": "FIXED", "max": 12}}
    }
}`

packager, err := iso8583.NewPackager(strings.NewReader(config))
if err != nil {
    panic(err)
}
```

## Supported MTI Types

The package includes predefined MTI types for common operations:

| MTI | Description |
|-----|-------------|
| 0100 | Card Processing Request |
| 0110 | Card Processing Response |
| 0200 | Financial Request |
| 0210 | Financial Response |
| 0400 | Reversal Request |
| 0410 | Reversal Response |
| 0401 | Repeated Reversal Request |
| 0800 | Network Management Request |
| 0810 | Network Management Response |

## TLV Support

The package includes support for TLV (Tag-Length-Value) data structures:

```go
tlvData, err := tlv.New([]byte{0x5F, 0x2A, 0x02, 0x08, 0x40}) // Example TLV data
if err != nil {
    panic(err)
}

int64Value := tlvData.GetInt64(0x5f2a) // will get int64(840)
byteValue := tlvData.GetBytes(0x5f2a) // will get []byte{0x08, 0x40}
stringValue := tlvData.GetHexString(0x5f2a) // will get "0840"

// Add TLV data to ISO message
msg.SetBytes(55, tlvData.Pack())
```

## Error Handling

The package defines several error types for different failure scenarios. Always check for errors after operations:

```go
if err := msg.UnpackISO(data); err != nil {
    if errors.Is(err, iso8583.ErrNotIsoMessage) {
        fmt.Println("Invalid ISO 8583 message")
    } else if errors.Is(err, iso8583.ErrInvalidBitType) {
        fmt.Println("Invalid bit type in message")
    }
    return err
}
```

## TODO

### High Priority
- [x] Add comprehensive data type validation for all field types (N, AN, ANS, B, Z)
  - Validate numeric fields contain only digits
  - Validate alphanumeric fields match character sets
  - Implement proper binary data validation
  - Add custom validation callbacks

### Completed
- [x] Optimize performance for packing and unpacking operations
- [x] Implement zero-allocation design for all operations
- [x] Add concurrent processing support with near-linear scaling

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
