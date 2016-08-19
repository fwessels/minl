
WebRTC Experiments
==================

Install
-------

```
go get github.com/keroserene/go-webrtc
```

Run
---

In order to establish a connection, for now the `offer` and `answer` need to be manually copied between the invoker and the server.

Run as follows for the invoker:

```
go run invoke/invoke.go --payload '{"key1":"value1", "key2":"value2", "key3":"value3"}'
```

and for the server

```
go run serve/serve.go
```

Expected output
---------------

Expected output should be similar to:

```
$ go run invoke/invoke.go --payload '{"key1":"value1", "key2":"value2", "key3":"value3"}'
Starting up PeerConnection...
Initializing datachannel....
Generating offer...
Finished gathering ICE candidates.

 ---- Please copy below to peer ---- 

{"type":"offer","sdp":"v=0\r\no=- 6303283194971796048 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=msid-semantic: WMS\r\nm=application 63762 DTLS/SCTP 5000\r\nc=IN IP4 77.160.254.108\r\na=candidate:3498311136 1 udp 2122260223 192.168.2.132 63762 typ host generation 0 network-id 1\r\na=candidate:1372322644 1 udp 1686052607 77.160.254.108 63762 typ srflx raddr 192.168.2.132 rport 63762 generation 0 network-id 1\r\na=candidate:2650800400 1 tcp 1518280447 192.168.2.132 55410 typ host tcptype passive generation 0 network-id 1\r\na=ice-ufrag:gtjlAytCagX5VlRR\r\na=ice-pwd:TQnotiP1VnLFrJeQy4ABz7gV\r\na=fingerprint:sha-256 11:0F:70:57:37:4A:BF:E7:44:78:99:A5:90:1A:90:5C:F1:8E:88:F6:F1:83:32:D3:88:03:FA:25:4E:BC:75:76\r\na=setup:actpass\r\na=mid:data\r\na=sctpmap:5000 webrtc-datachannel 1024\r\n"}

{"type":"answer","sdp":"v=0\r\no=- 5236220174479878421 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=msid-semantic: WMS\r\nm=application 61023 DTLS/SCTP 5000\r\nc=IN IP4 77.160.254.108\r\nb=AS:30\r\na=candidate:3498311136 1 udp 2122260223 192.168.2.132 61023 typ host generation 0 network-id 1\r\na=candidate:1372322644 1 udp 1686052607 77.160.254.108 61023 typ srflx raddr 192.168.2.132 rport 61023 generation 0 network-id 1\r\na=ice-ufrag:+PC29ig1c8mro9WA\r\na=ice-pwd:Wh4g/z8A4WQlA4HXl7Cv2hNr\r\na=fingerprint:sha-256 35:CB:AC:B4:79:9E:45:E8:4B:79:91:FA:82:52:5F:7E:B2:C2:56:53:AD:F1:F3:45:6F:0B:8E:BA:7C:D9:4E:4B\r\na=setup:active\r\na=mid:data\r\na=sctpmap:5000 webrtc-datachannel 1024\r\n"}
SDP answer successfully received.
Data Channel established
Invoking server with payload: {"key1":"value1", "key2":"value2", "key3":"value3"}
Server result: {"result":"response"}
```

