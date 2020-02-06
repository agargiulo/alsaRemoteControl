# alsaRemoteControl
Control your alsa mixer based audio system from a simple web(api?)

## Build/deploy
```sh
GOOS=linux GOARCH=arm go build .
scp alsaRemoteControl <rasp-pi_ip-or-domain>:
```

## Running
#### With TLS
```sh
export ALSA_REMOTE_SSL_CERT=cert.pem ALSA_REMOTE_SSL_KEY=key.pem
./alsaRemoteControl [8443]
```
#### No TLS
```sh
./alsaRemoteControl [8080]
```

##### Notes
if no port is specified, it will run on port 12345.


## API
None of this is a real REST API, just a heads up

|Request|Description|
|:-------|:-----------:|
|`GET /` | Loads the index page, with buttons \[`TODO`\] |
|`GET /volume` | Prints the current volume |
|`GET /toggle` | Mute or Unmute the audio |
|`GET /mute`   | Mute the audio |
|`GET /unmute` | Unmute the audio |
|`GET /up`     | Increase volume by 5% |
|`GET /down`   | Decrease volume by 5% |
|`POST /volume/(0-100)` | Sets the volume to the speficied value |