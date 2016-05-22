TOTP SSH port fluxing
===

Some people change their SSH port on their servers so that it is slightly harder to find for bots or other nasties, and while that is generally viewed as an action of [security through obscurity](https://en.wikipedia.org/wiki/Security_through_obscurity) it does work very well at killing a lot of the automated logins you always see in `/var/log/auth.log`

However what if we could go take this to a ridiculous level? What if we could use [TOTP](https://en.wikipedia.org/wiki/Time-based_One-time_Password_Algorithm) codes that are normally used as 2nd factor codes to login to websites to actually know what port the sshd server is listening on?

For this, I present [totp-ssh-flux](https://github.com/benjojo/totp-ssh-fluxer), a way to make sure your sshd port changes every 30 seconds, and possibly causing your adversaries a small period of frustration.

Demo:

![gif](https://blog.benjojo.co.uk/asset/O7HwIbd7i0)

What you can see here is my phone (using a generic TOTP client) generating codes, that I can then use as the port to SSH into on a server.

The software behind it is fairly simple, It runs in a loop that does the following

* Generates a TOTP token
* Takes the last digit, if the result is above 65536, do that again
* Adds a iptables PREROUTING rule to redirect that number generated above
* Waits 30 seconds, removes that rule, repeat.

The neat thing is, because this is done in `PREROUTING`, even if the code expires, established connections stay connected.

## Installation

### You will most likely find more up to date instructions on the [totp-ssh-flux](https://github.com/benjojo/totp-ssh-fluxer) project readme

### Beware, currently I would not really recommend running this software, it was only written as a joke.

At the time of writing the project is just a single file, You will need to install [golang](https://golang.org/) and then `go get` and `go build`

Run the program as root ( it needs to, sorry, it's editing iptables )

Upon first run, the program will generate a token for the host in `/etc/ssh-flux-key` ( you can use the `-keypath` option to change that ) and you can input that into your phone or other clients.


You can confirm it works by running `watch iptables -vL -t nat` and waiting for the iptables rules to be inserted and removed.

---

Want to see more insanity like this? Follow me on twitter [@benjojo12](https://twitter.com/Benjojo12)