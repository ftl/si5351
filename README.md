# Si5351 with Go on the Raspberry Pi

This is a library to use the Si5351 on the Raspberry Pi. It comes with a command line tool to control the Si5351 from the command line.

## Disclaimer

I develop this software for myself and just for fun in my free time. If you find it useful, I'm happy to hear about that. If you have trouble using it, you have all the source code to fix the problem yourself (although pull requests are welcome). 

## Build

To build for the Raspberry Pi:

```
GOARCH=arm GOARM=5 GOOS=linux go build
```

## License

This software is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)