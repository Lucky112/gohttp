# gohttp
Plain and simple tool to make http requests in parallel

# Usage
```
gohttp -parallel <uint> [url1 ... urlN]
```

gohttp makes GET requests to urls provided in args and prints md5 hash of recieved response next to each URL. In case of an error, it will be printed instead of a hash. The order of printed urls may differ from the order in args.
Flag -parallel limits a number of parallel request. By default 10 requests can be done in parallel. 

# Example
```
gohttp -parallel 2 http://google.com facebook.com empty.com
Rate of parallelism: 2
Number of URLs:  3

http://google.com cda79ac80a6285dc31befe3e81fa96cc
http://facebook.com f3a5ce37c64a3fbc4b6ff8a8b9e313f5
http://empty.com : unexpected status: 403 Forbidden
```

```
gohttp http://google.com facebook.com empty.com reddit.com tradingview.com vk.com yandex.ru moex.ru alibaba.com instagram.com rambler.com unknown.ru 
Rate of parallelism: 10
Number of URLs:  12

http://empty.com : unexpected status: 403 Forbidden
http://moex.ru 62ba604c1f0983f8f078d4883bd8067a
http://unknown.ru 0a85c83127b4e13f647f37bb21d85c82
http://yandex.ru 48987956ff2e3d6300b69b3dcc231c43
http://vk.com 008d396a9e4205a18f0dbde6e069e6dd
http://google.com c8876be8bf45f2fbf24da39d0f20f558
http://facebook.com d34b62555b8abe82d031dd624d69c834
http://instagram.com 13bd7247ea60ca309ff4187e3a916b08
http://reddit.com d2d6123b604e820f49e2df682adb51f6
http://alibaba.com 9b0e0a8b67d5ca2677f916ca0a9ea065
http://tradingview.com 62328d1698c456122ff6858c063d8fbe
http://rambler.com 50fc609da691d6dce8ec3e01fb42975c
```

# License
This code is licensed under MIT License. Feel free to contibute and/or use it for your own benefit.
