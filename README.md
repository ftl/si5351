# Si5351 with Go on the Raspberry Pi

This is a library to use the Si5351 on the Raspberry Pi. It comes with a command line tool to control the Si5351 from the command line.

## Disclaimer

I develop this software for myself and just for fun in my free time. If you find it useful, I'm happy to hear about that. If you have trouble using it, you have all the source code to fix the problem yourself (although pull requests are welcome). 

**Hint:** This software is currently work in progress. It is good enough for my experiments with the Si5351, but does not support all features of the device yet. Please be patient - pull requests are welcome.

## Usage

```
// define the crystal used with your device
crystal := si5351.Crystal{BaseFrequency: toCrystalFrequency(rootFlags.crystalFreq), Load: toCrystalLoad(rootFlags.crystalLoad), CorrectionPPM: rootFlags.ppm}

// open the I2C connection
bus, err := i2c.Open(rootFlags.address, rootFlags.bus)
if err != nil {
    log.Fatal(err)
}
defer bus.Close()
i2c.Debug = rootFlags.debugI2C

// create the device
device := si5351.New(crystal, bus)

// run the startup procedure
device.StartSetup()

// setup the PLL
device.SetupPLL(si5351.PLLA, 900*si5351.MHz)

// setup the output
device.PrepareOutputs(si5351.PLLA, false, si5351.ClockInputMultisynth, si5351.OutputDrive2mA, si5351.Clk1)
device.SetOutputFrequency(si5351.Clk1, frequency)

// finish the startup procedure
device.FinishSetup()

```

## Build

To build for the Raspberry Pi:

```
GOARCH=arm GOARM=7 GOOS=linux go build
```

## License

This software is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)