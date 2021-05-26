# Signpost

Signpost is a dead simple IRC bot you can put in one or more channels
to tell new folks to the channels a simple message.  This can be
useful if the network in question doesn't support an EntryMsg by
chanserv.

To configure, set the following variables:

  * `IRC_SERVER` - The connection string to point at the server.  You
    should provide this in host:port form.

  * `IRC_SSL` - Use SSL when connecting to the server, true if set.

  * `IRC_CHANNELS` - A comma seperated list of channels to join

  * `IRC_NICK` - A suitable nickname to join with.  You may optionally
    set a username and password to use with SASL.

  * `IRC_SASL` - Perform SASL negotiations during connect, true if set.

  * `IRC_USER` - SASL username.

  * `IRC_PASS` - SASL password.

  * `IRC_IGNOREHOSTS_REGEXP` - This regexp will be checked against the
    host component of a user's identity, and if it maches all messages
    that signpost would have delivered will be suppressed.

  * `IRC_MSGS` - A set of messages to be returned for each channel.  The format is:

  ```
  channel1:message for channel 1;channel2:message for channel 2
  ```

  Note that there is a maximum limit on environment variables, so keep
  your messages succinct.  Use a paste service if you need to.
