# alsaRemoteControl
Control your alsa mixer based audio system from a simple web(api?)

```
GOOS=linux GOARCH=arm go build .
scp alsaRemoteControl <rasp-pi_ip-or-domain>:
<...>
./alsaRemoteControl
```

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