# ISO 8583 Implementation in Go

A comprehensive implementation of the ISO 8583 financial transaction messaging standard in Go. This package provides a robust and flexible way to parse, build, and process ISO 8583 messages.

## Features

- Support for both fixed and variable length fields (LLVAR, LLLVAR, LLLLVAR)
- Built-in support for common MTI (Message Type Identifier) types
- Flexible message packing and unpacking
- Support for binary and ASCII message formats
- Configurable field definitions and validation
- TLV (Tag-Length-Value) support for complex data structures
- Comprehensive error handling

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
    msg.MTI = iso8583.MTIFinancialRequestByte
    
    // Set fields
    msg.SetString(2, "4761739001010119")  // Primary Account Number
    msg.SetString(3, "000000")            // Processing Code
    msg.SetString(4, "000000010000")      // Transaction Amount
    
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
    mti := string(msg.MTI)
    pan, _ := msg.GetString(2)
    
    fmt.Printf("MTI: %s\n", mti)
    fmt.Printf("PAN: %s\n", pan)
}
```

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
- [ ] Add comprehensive data type validation for all field types (N, AN, ANS, B, Z)
  - Validate numeric fields contain only digits
  - Validate alphanumeric fields match character sets
  - Implement proper binary data validation
  - Add custom validation callbacks

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.