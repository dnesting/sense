# sense/realtime

This is kind of implemented as its own API and is even less stable than `[sense](..)`.
The [asyncapi.yaml] is an untested but nearly complete description of the
Sense real-time API.  I would have preferred to use a code generator, similar to how I did
in the [internal/client](../internal/client), but the code generators for AsyncAPI are not
mature and I couldn't get them to produce anything useful.  If you're reading this from
a timeline in which they are useful, consider splitting this package into an `internal/realtime`
(or integrate it with `internal/client`) and possibly merge this package or a higher-level
version of it with `sense`.